package main

func main() {
	paraseFlags()
	initConfig()
	destoryTemp(true)
	login()
	getBookInfo()
	downloadChapters(chapterIDs)
	destoryTemp(false)
}
