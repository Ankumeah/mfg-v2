package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"
)

func downloadFile(url string, savePath string) error {
	if _, err := os.Stat(savePath); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) && err != nil {
		return err
	}

	res, err := client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return errors.New("http error: " + res.Status)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, res.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(savePath, buf.Bytes(), 0o644)
	if err != nil {
		return err
	}

	return nil
}
