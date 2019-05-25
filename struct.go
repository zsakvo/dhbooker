package main

import (
	"runtime"
	"sync"

	"github.com/Unknwon/goconfig"
	"github.com/tidwall/gjson"
	"gopkg.in/cheggaaa/pb.v1"
)

const appVersion = "2.1.032"
const initEncryptKey = "zG2nSeEfSHfvTCHy5LCcqtBbQehKNLXn"

var ostype = runtime.GOOS
var pathSeparator string

var bar *pb.ProgressBar
var ping int64
var book cbook
var token ctoken
var path pathSettings
var config *goconfig.ConfigFile
var download cdownload

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

type cbook struct {
	format          string
	coverURL        string
	id              string
	name            string
	author          string
	tmpPath         string
	chapterNum      int
	rollNum         int
	chapters        sync.Map
	invalidChapters sync.Map
	chapterIDs      []string
	rolls           []gjson.Result
}

type cdownload struct {
	failedChapters []string
}
