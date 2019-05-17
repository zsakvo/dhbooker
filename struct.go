package main

import (
	"sync"

	"github.com/Unknwon/goconfig"
	"github.com/schollz/progressbar"
	"github.com/tidwall/gjson"
)

const appVersion = "2.1.032"
const initEncryptKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"

var downloadType string
var bookCoverURL string
var tocNcx string
var contentOpf string
var bookID string
var bookName string
var bookAuthor string
var bookChapterNum int
var validChapterNum int
var invalidChapterID []string
var chapterTitles []string
var chapterIDs []gjson.Result
var validChapterIDs []string
var chapters sync.Map
var invalidChapters sync.Map
var bookTmpPath string

// var chapters map[string]string
// var counter = struct{
//     sync.RWMutex
//     chapters map[string]string
// }
var bar progressbar.ProgressBar

var token ctoken
var path pathSettings
var config *goconfig.ConfigFile

type ctoken struct {
	readerID   string
	loginToken string
	account    string
}

type accountSettings struct {
	username string
	password string
}

type pathSettings struct {
	tmp string
	out string
}
