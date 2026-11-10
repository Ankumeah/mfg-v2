package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

func parseArgs() {
	if len(os.Args) < 1 {
		fmt.Fprintln(os.Stderr, "argv is of len 0. Your runnning this on a potato arent you?")
	}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [OPTIONS]... [MANGAS]...\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "  --help -h")
		fmt.Fprintln(os.Stderr, "\tPrint this text")
	}

	flag.UintVar(&maxDownloads, "j", 10, "Max parallel downloads, Default 10")
	t := flag.String("t", "10s", "Max download timeout. Default 10s. Eg - 1h1m1s1ms1us1ns")
	flag.Parse()
	mangas = flag.Args()

	if len(mangas) < 1 {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Must provide atlest one manga")
		os.Exit(1)
	}

	var err error
	timeout, err = time.ParseDuration(*t)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid timeout: %v\n", err.Error())
		os.Exit(1)
	}

	client = &http.Client{Timeout: timeout}
}
