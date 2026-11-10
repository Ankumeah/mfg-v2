package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func downloadCh(sem chan struct{}, ch string, manga string, savePath string) error {
	m := strings.ToLower(manga)

	res, err := client.Get(baseMangaPageUrl + manga + "_" + ch)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	pages := imageRegex.FindAllString(string(html), -1)

	dir, err := os.MkdirTemp("", fmt.Sprintf("mfg_ch_%s_*", ch))
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	maxLen := 1
	for _, page := range pages {
		maxLen = max(len(page), maxLen)
	}

	saves := []string{}

	var wg sync.WaitGroup
	var e error
	for _, page := range pages {
		save := page
		for range maxLen - len(save) {
			save = "0" + save
		}
		save = filepath.Join(dir, save)

		if e == nil {
			sem <- struct{}{}
			wg.Go(func() {
				err = downloadFile(
					fmt.Sprintf(
						"%s%s/%s_%s/%s_%s_%s",
						baseImagesUrl, m, m, ch, m, ch, page,
					), save,
				)
				<-sem
				if err != nil {
					e = err
				}
			})

			saves = append(saves, save)
		}

	}
	wg.Wait()
	if e != nil {
		return err
	}

	return zipFiles(saves, savePath)
}

func zipFiles(files []string, savePath string) error {
	buf := bytes.NewBuffer(nil)

	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	for _, file := range files {
		fileToZip, err := os.Open(file)
		if err != nil {
			return err
		}
		defer fileToZip.Close()

		w, err := zipWriter.Create(filepath.Base(file))
		if err != nil {
			return err
		}

		if _, err := io.Copy(w, fileToZip); err != nil {
			return err
		}
	}

	return os.WriteFile(savePath, buf.Bytes(), 0o644)
}
