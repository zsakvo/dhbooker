package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/schollz/progressbar"
	"github.com/tidwall/gjson"
)

var quit = make(chan int)

func httpGet(url string) string {
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	check(err)
	return string(s)
}

func httpPost(url string, content string) string {
	resp, err := http.Post(url,
		"application/x-www-form-urlencoded",
		strings.NewReader(content))
	check(err)
	body, err := ioutil.ReadAll(resp.Body)
	check(err)
	return string(body)
}

//账号登入
func loginByPass(account accountSettings) {
	name := account.username
	pass := account.password
	url := "https://app.hbooker.com/signup/login?login_name=" + name + "&app_version=" + appVersion + "&passwd=" + pass
	body := httpGet(url)
	body = decode(body, initEncryptKey)
	code := gjson.Get(body, "code").String()
	switch {
	case code == "210002":
		fmt.Println("用户不存在，请重写配置文件")
	case code == "210003":
		fmt.Println("密码错误，请重写配置文件")
	case code == "100000":
		localToken := gjson.Get(body, "data").String()
		setConfig("token", "token", localToken)
		writeConfig()
		token.readerID = gjson.Get(localToken, "reader_info.reader_id").String()
		token.loginToken = gjson.Get(localToken, "login_token").String()
		token.account = gjson.Get(localToken, "reader_info.account").String()
		clear()
	}
}

//凭证登入
func loginByToken(localToken string) {
	token.readerID = gjson.Get(localToken, "reader_info.reader_id").String()
	token.loginToken = gjson.Get(localToken, "login_token").String()
	token.account = gjson.Get(localToken, "reader_info.account").String()
	url := "https://app.hbooker.com/reader/get_my_info?reader_id=" + token.readerID + "&app_version=" + appVersion + "&login_token=" + token.loginToken + "&account=" + token.account
	body := httpGet(url)
	body = decode(body, initEncryptKey)
	code := gjson.Get(body, "code").String()
	if code != "100000" {
		fmt.Println(gjson.Get(body, "tip"))
		fmt.Println("凭证登入失败，尝试使用用户密码登入")
		account := getAccountSettings()
		loginByPass(account)
	} else {
		clear()
	}
}

//获取书籍名称
func getBookName() string {
	bookInfoURL := "https://app.hbooker.com/book/get_info_by_id?module_id=&tab_type=&app_version=" + appVersion + "&recommend=&carousel_position=&book_id=" + bookID + "&login_token=" + token.loginToken + "&book_id=" + bookID + "&account=" + token.account
	body := httpGet(bookInfoURL)
	body = decode(body, initEncryptKey)
	code := gjson.Get(body, "code").String()
	if code != "100000" {
		tip := gjson.Get(body, "tip").String()
		fmt.Println(tip)
		os.Exit(1)
	}
	return gjson.Get(body, "data.book_info.book_name").String()
}

//获取分卷信息
func getBookRolls() ([]gjson.Result, int) {
	url := "https://app.hbooker.com/book/get_division_list?app_version=" + appVersion + "&login_token=" + token.loginToken + "&book_id=" + bookID + "&account=" + token.account
	body := httpGet(url)
	body = decode(body, initEncryptKey)
	rolls := gjson.Get(body, "data.division_list.#.division_id").Array()
	num := len(rolls)
	return rolls, num
}

//获取章节信息
func getChapters(rolls []gjson.Result) ([]gjson.Result, int) {
	for _, roll := range rolls {
		content := "last_update_time=0&app_version=" + appVersion + "&division_id=" + roll.String() + "&login_token=" + token.loginToken + "&account=" + token.account
		body := httpPost("https://app.hbooker.com/chapter/get_updated_chapter_by_division_id", content)
		body = decode(body, initEncryptKey)
		chapterIDs = append(chapterIDs, gjson.Get(body, "data.chapter_list.#.chapter_id").Array()...)
	}
	num := len(chapterIDs)
	return chapterIDs, num
}

//获取章节内容
func getChapterContent(chapterID string) (string, string, bool) {
	contentKeyURL := "https://app.hbooker.com/chapter/get_chapter_cmd?app_version=" + appVersion + "&chapter_id=" + chapterID + "&login_token=" + token.loginToken + "&account=" + token.account
	keyBody := httpGet(contentKeyURL)
	keyBody = decode(keyBody, initEncryptKey)
	contentKey := gjson.Get(keyBody, "data.command").String()
	contentURL := "https://app.hbooker.com/chapter/get_cpt_ifm?chapter_command=" + contentKey + "&app_version=" + appVersion + "&chapter_id=" + chapterID + "&login_token=" + token.loginToken + "&account=" + token.account
	contentBody := httpGet(contentURL)
	contentBody = decode(contentBody, initEncryptKey)
	chapterTitle := gjson.Get(contentBody, "data.chapter_info.chapter_title").String()
	content := gjson.Get(contentBody, "data.chapter_info.txt_content").String()
	auth := gjson.Get(contentBody, "data.chapter_info.auth_access").String()
	if len(content) == 0 {
		return "", "", false
	}
	if auth == "0" {
		return "", "", false
	}
	validChapterNum++
	content = decode(content, contentKey)
	return chapterTitle, chapterTitle + "\n\n" + content + "\n\n\n\n", true
}

//下载索引
var di int

//写出章节缓存
func writeChapterTemp(chapterID string) {
	tmpPath := path.tmp + "/" + bookName
	os.MkdirAll(tmpPath, os.ModePerm)
	_, content, result := getChapterContent(chapterID)
	if result {
		d := []byte(content)
		filePath := tmpPath + "/" + chapterID + ".txt"
		err := ioutil.WriteFile(filePath, d, 0644)
		check(err)
	}
	bar.Add(1)
	if di == bookChapterNum {
		quit <- 1
	}
	di++
}

//下载章节
func downloadChapters(chapterIDs []gjson.Result) {
	di = 1
	bar = *progressbar.New(len(chapterIDs))
	for _, chapterID := range chapterIDs {
		go writeChapterTemp(chapterID.String())
	}
	<-quit
	mergeTemp()
	bar.Finish()
}
