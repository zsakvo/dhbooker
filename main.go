package main

import (
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
	fmt.Print("请输入要下载的书籍的书号：")
	fmt.Scanln(&bookID)
	bookName = getBookName()
	fmt.Println("《" + bookName + "》")
	rolls, rollNum := getBookRolls()
	chapterIDs, chapterNum := getChapters(rolls)
	bookChapterNum = chapterNum
	fmt.Println("共" + strconv.Itoa(rollNum) + "卷，" + strconv.Itoa(chapterNum) + "章")
	downloadChapters(chapterIDs)
	destoryTemp()
	fmt.Print("\n下载完毕")
}
