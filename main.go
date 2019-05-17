package main

import (
	"flag"
	"fmt"
	"strconv"
)

func main() {
	getConfig()
	path = getPathSettings()
	fmt.Println("正在清理临时目录")
	destoryTemp()
	fmt.Println("正在检查凭证")
	token := getToken()
	if len(token) == 0 {
		fmt.Println("无凭证，尝试使用用户密码登入")
		account := getAccountSettings()
		loginByPass(account)
	} else {
		fmt.Println("尝试使用上次凭证登入")
		localToken := getSection("token")["token"]
		loginByToken(localToken)
	}
	flag.StringVar(&downloadType, "t", "txt", "下载类型，可选 txt 或 epub")
	flag.StringVar(&bookID, "b", "", "bookID，请在对应网页 url 中获取")
	flag.Parse()
	getBookInfo()
	bookTmpPath = path.tmp + "/" + bookName + "/"
	fmt.Println("《" + bookName + "》")
	rolls, rollNum := getBookRolls()
	chapterIDs, chapterNum := getChapters(rolls)
	bookChapterNum = chapterNum
	fmt.Println("共" + strconv.Itoa(rollNum) + "卷，" + strconv.Itoa(chapterNum) + "章")
	downloadChapters(chapterIDs)
	destoryTemp()
	fmt.Print("\n下载完毕")
}
