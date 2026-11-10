package main

import (
	"io"
	"strings"
)

func getChapters(manga string) ([]string, error) {
	res, err := client.Get(baseMangaUrl + manga)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	html, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	match := regex1.FindAllString(string(html), -1)
	match = regex2.FindAllString(strings.Join(match, "\n"), -1)

	seen := make(map[string]struct{})
	var result []string
	for _, ch := range match {
		if _, ok := seen[ch]; ok {
			continue
		}

		seen[ch] = struct{}{}
		result = append(result, ch)
	}

	return result, nil
}
