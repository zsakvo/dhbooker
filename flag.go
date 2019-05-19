package main

import "flag"

func paraseFlags() {
	flag.StringVar(&book.format, "t", "txt", "下载类型，可选 txt 或 epub")
	flag.Int64Var(&ping, "p", 5000, "超时时间，单位为毫秒")
	flag.StringVar(&book.id, "b", "", "bookID，请在对应网页 url 中获取")
	flag.Parse()
}
