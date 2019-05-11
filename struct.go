package main

import (
	"github.com/Unknwon/goconfig"
	"github.com/schollz/progressbar"
	"github.com/tidwall/gjson"
)

const appVersion = "2.1.032"
const initEncryptKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"

var bookID string
var bookName string
var bookChapterNum int
var validChapterNum int
var invalidChapterID []string
var chapterIDs []gjson.Result
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
