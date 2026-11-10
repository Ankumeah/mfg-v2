default: run

run *args:
  go build -ldflags="-s -w" -o ./bin/mfg .
  ./bin/mfg {{ args }}
