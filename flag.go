package main

import "flag"

func paraseFlags() {
	flag.StringVar(&downloadType, "t", "txt", "下载类型，可选 txt 或 epub")
	flag.StringVar(&bookID, "b", "", "bookID，请在对应网页 url 中获取")
	flag.Parse()
}
