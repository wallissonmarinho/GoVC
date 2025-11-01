package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	var workers int
	var saveLogs bool
	flag.IntVar(&workers, "p", runtime.NumCPU(), "number of parallel ffmpeg processes")
	flag.BoolVar(&saveLogs, "logs", true, "save per-file ffmpeg logs to mp4/*.log (default true). Use -logs=false to disable")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: go run main.go [-p workers] /path/to/folder")
		return
	}

	dir := args[0]

	// Verifica se o diretório existe
	info, err := os.Stat(dir)
	if err != nil || !info.IsDir() {
		log.Fatalf("Directory not found: %s", dir)
	}

	// Cria a pasta mp4
	outDir := filepath.Join(dir, "mp4")
	if err := os.MkdirAll(outDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Lista arquivos MKV
	mkvFiles, err := filepath.Glob(filepath.Join(dir, "*.mkv"))
	if err != nil {
		log.Fatalf("Failed to list MKV files: %v", err)
	}

	if workers < 1 {
		workers = 1
	}

	sem := make(chan struct{}, workers)
	var wg sync.WaitGroup

	if len(mkvFiles) == 0 {
		fmt.Println("No .mkv files found to convert.")
	}

	// progress tracking
	progress := make(map[string]float64)
	var progMu sync.Mutex
	var completed int32
	total := len(mkvFiles)

	// printer goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			progMu.Lock()
			if len(progress) == 0 && atomic.LoadInt32(&completed) >= int32(total) {
				progMu.Unlock()
				return
			}
			for name, p := range progress {
				fmt.Printf("%s: %.1f%%\n", name, p)
			}
			progMu.Unlock()
		}
	}()

	for _, f := range mkvFiles {
		f := f // capture
		wg.Add(1)
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			filename := filepath.Base(f)
			nameOnly := strings.TrimSuffix(filename, filepath.Ext(filename))
			srtPath := filepath.Join(dir, nameOnly+".srt")

			var args []string
			outPath := filepath.Join(outDir, nameOnly+".mp4")
			logPath := filepath.Join(outDir, nameOnly+".log")

			// get duration via ffprobe
			durationSec, _ := getDuration(f)

			if _, err := os.Stat(srtPath); err == nil {
				fmt.Printf("Started: %s (with subtitles)\n", filename)
				args = []string{
					"-i", f,
					"-i", srtPath,
					"-map", "0:v",
					"-map", "0:a",
					// include internal subtitle streams (if any) and the external srt input's subtitles
					"-map", "0:s?",
					"-map", "1:s?",
					"-c:v", "copy",
					"-c:a", "aac", "-b:a", "192k",
					"-c:s", "mov_text",
					"-progress", "pipe:1",
					"-nostats",
					outPath,
				}
			} else {
				fmt.Printf("Started: %s (no subtitles)\n", filename)
				args = []string{
					"-i", f,
					"-map", "0:v",
					"-map", "0:a",
					// include internal subtitle streams if present
					"-map", "0:s?",
					"-c:v", "copy",
					"-c:a", "aac", "-b:a", "192k",
					"-c:s", "mov_text",
					"-progress", "pipe:1",
					"-nostats",
					outPath,
				}
			}

			cmd := exec.Command("ffmpeg", args...)

			// open logfile if enabled
			var lf *os.File
			if saveLogs {
				lf, err = os.Create(logPath)
				if err != nil {
					// fallback to discarding stderr if we can't create a log file
					cmd.Stderr = io.Discard
				} else {
					// don't defer Close() here because we may want to delete the file after cmd finishes;
					// close it explicitly after cmd.Wait()
					cmd.Stderr = lf
				}
			} else {
				// user opted out of logs: stream ffmpeg stderr to the terminal
				cmd.Stderr = os.Stderr
			}

			stdout, err := cmd.StdoutPipe()
			if err != nil {
				log.Printf("Failed to get stdout pipe for %s: %v", filename, err)
			}

			if err := cmd.Start(); err != nil {
				log.Printf("Failed to start ffmpeg for %s: %v", filename, err)
				if lf != nil {
					_, _ = lf.WriteString("failed to start ffmpeg: " + err.Error() + "\n")
				}
				<-sem
				atomic.AddInt32(&completed, 1)
				return
			}

			// parse progress from stdout
			scanner := bufio.NewScanner(stdout)
			var lastPercent float64
			for scanner.Scan() {
				line := scanner.Text()
				if lf != nil {
					_, _ = lf.WriteString(line + "\n")
				}
				if strings.Contains(line, "=") {
					parts := strings.SplitN(line, "=", 2)
					key := parts[0]
					val := parts[1]
					if key == "out_time_ms" {
						if v, err := strconv.ParseFloat(val, 64); err == nil {
							// out_time_ms is in microseconds
							processed := v / 1000000.0
							if durationSec > 0 {
								pct := (processed / durationSec) * 100.0
								if pct > 100 {
									pct = 100
								}
								progMu.Lock()
								progress[nameOnly] = pct
								progMu.Unlock()
								lastPercent = pct
							}
						}
					} else if key == "out_time" {
						if t, err := parseOutTime(val); err == nil {
							if durationSec > 0 {
								pct := (t / durationSec) * 100.0
								if pct > 100 {
									pct = 100
								}
								progMu.Lock()
								progress[nameOnly] = pct
								progMu.Unlock()
								lastPercent = pct
							}
						}
					} else if key == "progress" && val == "end" {
						progMu.Lock()
						progress[nameOnly] = 100
						progMu.Unlock()
						lastPercent = 100
					}
				}
			}

			// wait for ffmpeg to finish
			waitErr := cmd.Wait()
			if waitErr != nil {
				log.Printf("Failed to convert %s: %v (see %s)", filename, waitErr, logPath)
			}

			// ensure final percentage recorded
			progMu.Lock()
			progress[nameOnly] = lastPercent
			progMu.Unlock()

			// close logfile if we opened it
			if lf != nil {
				_ = lf.Close()
			}

			// if logs are enabled and conversion succeeded, remove the log file
			// Only delete if ffmpeg exited successfully and the output file exists and has size > 0
			if saveLogs && waitErr == nil {
				if st, err := os.Stat(outPath); err == nil {
					if st.Size() > 0 {
						if err := os.Remove(logPath); err != nil {
							// non-fatal: just log the failure to remove
							log.Printf("Warning: failed to remove log %s: %v", logPath, err)
						}
					} else {
						log.Printf("Not removing log for %s because output file is empty", filename)
					}
				} else {
					log.Printf("Not removing log for %s because output file not found: %v", filename, err)
				}
			}

			atomic.AddInt32(&completed, 1)

			fmt.Printf("Finished: %s -> %s\n", filename, outPath)
		}()
	}

	wg.Wait()
}

// getDuration uses ffprobe to retrieve the duration in seconds of a media file.
func getDuration(path string) (float64, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	s := strings.TrimSpace(string(out))
	if s == "" {
		return 0, fmt.Errorf("ffprobe returned empty duration")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// parseOutTime parses ffmpeg out_time formatted as HH:MM:SS.micro to seconds
func parseOutTime(s string) (float64, error) {
	// expect e.g. 00:01:23.456789
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid out_time format: %s", s)
	}
	h, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, err
	}
	m, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, err
	}
	sec, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0, err
	}
	return h*3600 + m*60 + sec, nil
}
