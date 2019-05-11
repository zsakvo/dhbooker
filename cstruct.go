package main

import (
	"github.com/tidwall/gjson"
)

//Token1 struct
type Token1 struct {
	ReaderID   string
	LoginToken string
	Account    string
	AppVersion string
}

//Book1 struct
type Book1 struct {
	BookName   string
	Division   string
	Volumes    []gjson.Result
	VolumeNum  int
	ChapterIDs []gjson.Result
	ChapterNum int
}

//Settings1 struct
type Settings1 struct {
	TmpPath    string
	OutputPath string
}
