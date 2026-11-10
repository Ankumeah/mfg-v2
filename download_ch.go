package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	for _, page := range pages {
		save := page
		for range maxLen - len(save) {
			save = "0" + save
		}
		save = filepath.Join(dir, save)

		sem <- struct{}{}
		if err := downloadFile(
			fmt.Sprintf(
				"%s%s/%s_%s/%s_%s_%s",
				baseImagesUrl, m, m, ch, m, ch, page,
			), save,
		); err != nil {
			return err
		}
		<-sem

		saves = append(saves, save)
	}

	buf := bytes.NewBuffer(nil)

	zipWriter := zip.NewWriter(buf)
	defer zipWriter.Close()

	for _, save := range saves {
		if err := addFile(zipWriter, save); err != nil {
			return err
		}
	}

	if err := os.WriteFile(savePath, buf.Bytes(), 0o644); err != nil {
		return err
	}

	return nil
}

func addFile(zipWriter *zip.Writer, file string) error {
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

	return nil
}
