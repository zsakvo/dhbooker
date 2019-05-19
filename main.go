package main

func main() {
	paraseFlags()
	initConfig()
	destoryTemp(true)
	login()
	getBookInfo()
	downloadChapters()
	destoryTemp(false)
}
