package main

import (
	"net/http"
	"regexp"
	"time"
)

const baseUrl = "https://images.mangafreak.net/downloads/"
const baseImagesUrl = "https://images.mangafreak.me/mangas/"
const baseMangaUrl = "https://ww3.mangafreak.me/Manga/"
const baseMangaPageUrl = "https://ww3.mangafreak.me/Read1_"

var regex1 *regexp.Regexp = regexp.MustCompile(`/Read1_[^"]*`)
var regex2 *regexp.Regexp = regexp.MustCompile(`(?m)[0-9]+[a-z]*$`)

var imageRegex *regexp.Regexp = regexp.MustCompile(`[0-9]+\.jpg`)

var maxDownloads uint
var timeout time.Duration
var mangas []string
var client *http.Client
