# GoVC

Ferramenta simples para converter arquivos MKV para MP4 em lote, com paralelismo e suporte a legendas.

Principais comportamentos

- Paralelismo: use `-p N` para controlar quantos processos `ffmpeg` rodam simultaneamente. Padrão: número de CPUs.
- Logs: por padrão o programa salva logs temporários por arquivo em `mp4/<nome>.log`. Use `-logs=false` para desativar.
- Remoção de logs: quando `-logs=true`, os arquivos `.log` são removidos automaticamente ao final da conversão do respectivo vídeo apenas se a conversão teve sucesso e o arquivo MP4 de saída existir com tamanho > 0. Caso contrário, o `.log` é mantido para diagnóstico.

# GoVC

Simple batch MKV → MP4 converter with parallel processing and subtitle handling.

Key behaviors

- Parallelism: use `-p N` to control how many ffmpeg processes run in parallel. Default: number of CPUs.
- Logs: by default the program writes a per-file temporary log to `mp4/<name>.log`. Use `-logs=false` to disable log files.
- Log removal: when `-logs=true`, a `.log` file is removed automatically after the corresponding conversion finishes only if the conversion succeeded and the output MP4 exists and has size > 0. Otherwise the `.log` is kept for debugging.

Subtitles

- The tool attempts to include both embedded (internal) and external subtitles (a `.srt` file with the same base name):
  - Internal subtitle streams are mapped with `-map 0:s?` (optional mapping).
  - If a `name.srt` exists, it is added as an extra input and mapped with `-map 1:s?`.
  - Text-based subtitles are converted to `mov_text` (MP4-compatible) using `-c:s mov_text`.

Important limitations

- Image-based subtitles (PGS, VobSub) cannot be converted to `mov_text`. In those cases the program does not burn-in the subtitles automatically — such subtitles will be ignored as text tracks.
- If you want automatic burn-in of image-based subtitles, that requires re-encoding the video (much slower). I can add that option if you want.

Requirements

- `ffmpeg` and `ffprobe` must be installed and available in your PATH.

Usage examples

- Convert using 4 parallel processes and keep temporary logs (default):

```bash
go run main.go -p 4 /path/to/folder
```

- Convert using 2 parallel processes and do not save per-file logs (ffmpeg stderr will be shown in the terminal):

```bash
go run main.go -p 2 -logs=false /path/to/folder
```

Notes

- Temporary log files are useful for debugging failed conversions — they are preserved when an error occurs.
- Future improvements I can add on request:
  - Automatic detection and reporting of subtitle streams that can't be converted (before deleting logs),
  - `--burn` option to burn image-based subtitles into the video via re-encode,
  - `--log-dir` option to choose a custom directory for logs instead of `mp4/`.
