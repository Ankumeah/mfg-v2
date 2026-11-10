package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func main() {
  parseArgs()

	for _, manga := range mangas {
		fmt.Printf("Downloading %v\n", manga)

		chs, err := getChapters(manga)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error while getting %v: %v\n", manga, err.Error())
			continue
		}

		maxPad := 1
		for i, ch := range chs {
			if ch == "" {
				chs = append(chs[:i], chs[i+1:]...)
			}
			if len(ch) > maxPad {
				maxPad = len(ch)
			}
		}

		fmt.Printf("Found %v chapters\n", len(chs))

		var prog int
		go func() {
			for {
				if prog == len(chs) {
					return
				}
				fmt.Fprintf(os.Stderr, "\r%v / %v", len(chs), prog)
				time.Sleep(time.Second * 1)
			}
		}()

		os.MkdirAll(manga, 0o755)

		var wg sync.WaitGroup
		sem := make(chan struct{}, maxDownloads)
		for _, ch := range chs {
			url := baseUrl + manga + "_" + ch
      rawCh := ch

			for range maxPad - len(ch) {
				ch = "0" + ch
			}

			sem <- struct{}{}
			wg.Go(func() {
				if err := downloadFile(
					url, filepath.Join(manga, ch + ".cbz"),
				); err != nil {
					fmt.Fprintf(os.Stderr, "\nError while downloading chapter %v (falling back to page wise download): %v\n", ch, err.Error())
          if err := downloadCh(sem, rawCh, manga, filepath.Join(manga, ch + ".cbz")); err != nil {
					  fmt.Fprintf(os.Stderr, "\nPage wise download failed %v: %v\n", ch, err.Error())
          }
				}
				<-sem
				prog++
			})
		}
		wg.Wait()
		fmt.Fprintf(os.Stderr, "\r%v / %v", len(chs), prog)
		fmt.Fprintln(os.Stderr)
	}
}
