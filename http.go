package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/progressbar"
	"github.com/tidwall/gjson"
)

var quit = make(chan int)

func httpGet(url string, paramsMap map[string]string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Duration(ping * int64(time.Millisecond)),
	}
	request, err := http.NewRequest("GET", url, nil)
	request.Header.Set("User-Agent", "dhbooker")
	params := request.URL.Query()
	if paramsMap != nil {
		for m, n := range paramsMap {
			params.Add(m, n)
		}
		request.URL.RawQuery = params.Encode()
	}
	response, err := client.Do(request)
	check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	check(err)
	return string(body)
}

func httpPost(url string, content string) string {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := http.Client{
		Transport: tr,
		Timeout:   time.Duration(ping * int64(time.Millisecond)),
	}
	request, err := http.NewRequest("POST", url, strings.NewReader(content))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("User-Agent", "dhbooker")
	response, err := client.Do(request)
	check(err)
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	check(err)
	return string(body)
}

//账号登入
func loginByPass(account accountSettings) {
	name := account.username
	pass := account.password
	url := "https://app.hbooker.com/signup/login"
	paramsMap := map[string]string{"login_name": name, "app_version": appVersion, "passwd": pass}
	body := httpGet(url, paramsMap)
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
	url := "https://app.hbooker.com/reader/get_my_info"
	paramsMap := map[string]string{"reader_id": token.readerID, "app_version": appVersion, "login_token": token.loginToken, "account": token.account}
	body := httpGet(url, paramsMap)
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
func getBookInfo() {
	if len(book.id) == 2 {
		println("请至少输入书籍 ID")
		os.Exit(1)
	}
	bookInfoURL := "https://app.hbooker.com/book/get_info_by_id"
	paramsMap := map[string]string{"app_version": appVersion, "login_token": token.loginToken, "book_id": book.id, "account": token.account}
	body := httpGet(bookInfoURL, paramsMap)
	body = decode(body, initEncryptKey)
	code := gjson.Get(body, "code").String()
	if code != "100000" {
		tip := gjson.Get(body, "tip").String()
		fmt.Println(tip)
		os.Exit(1)
	}
	book.name = gjson.Get(body, "data.book_info.book_name").String()
	book.author = gjson.Get(body, "data.book_info.author_name").String()
	book.coverURL = gjson.Get(body, "data.book_info.cover").String()
	book.tmpPath = path.tmp + pathSeparator + book.name + pathSeparator
	fmt.Println("《" + book.name + "》")
	getBookRolls()
	getChapters()
	fmt.Println("共" + strconv.Itoa(book.rollNum) + "卷，" + strconv.Itoa(book.chapterNum) + "章")
}

//获取分卷信息
func getBookRolls() {
	url := "https://app.hbooker.com/book/get_division_list"
	paramsMap := map[string]string{"app_version": appVersion, "login_token": token.loginToken, "book_id": book.id, "account": token.account}
	body := httpGet(url, paramsMap)
	body = decode(body, initEncryptKey)
	book.rolls = gjson.Get(body, "data.division_list.#.division_id").Array()
	book.rollNum = len(book.rolls)
}

//获取章节信息
func getChapters() {
	for _, roll := range book.rolls {
		content := "last_update_time=0&app_version=" + appVersion + "&division_id=" + roll.String() + "&login_token=" + token.loginToken + "&account=" + token.account
		body := httpPost("https://app.hbooker.com/chapter/get_updated_chapter_by_division_id", content)
		body = decode(body, initEncryptKey)
		book.chapterIDs = append(book.chapterIDs, gjson.Get(body, "data.chapter_list.#.chapter_id").Array()...)
	}
	book.chapterNum = len(book.chapterIDs)
}

//获取章节内容
func getChapterContent(chapterID string) (string, string, bool) {
	contentKeyURL := "https://app.hbooker.com/chapter/get_chapter_cmd?app_version=" + appVersion + "&chapter_id=" + chapterID + "&login_token=" + token.loginToken + "&account=" + token.account
	paramsMap := map[string]string{"app_version": appVersion, "chapter_id": chapterID, "login_token": token.loginToken, "account": token.account}
	keyBody := httpGet(contentKeyURL, paramsMap)
	keyBody = decode(keyBody, initEncryptKey)
	contentKey := gjson.Get(keyBody, "data.command").String()
	contentURL := "https://app.hbooker.com/chapter/get_cpt_ifm"
	paramsMap1 := map[string]string{"chapter_command": contentKey, "app_version": appVersion, "login_token": token.loginToken, "chapter_id": chapterID, "account": token.account}
	contentBody := httpGet(contentURL, paramsMap1)
	contentBody = decode(contentBody, initEncryptKey)
	chapterTitle := gjson.Get(contentBody, "data.chapter_info.chapter_title").String()
	content := gjson.Get(contentBody, "data.chapter_info.txt_content").String()
	auth := gjson.Get(contentBody, "data.chapter_info.auth_access").String()
	if len(content) == 0 {
		book.invalidChapters.Store(chapterID, "1")
		return "", "", false
	}
	if auth == "0" {
		book.invalidChapters.Store(chapterID, "1")
		return "", "", false
	}
	// validChapterNum++
	content = decode(content, contentKey)
	// validChapterIDs = append(validChapterIDs, chapterID)
	if book.format == "epub" {
		titleElement := "<h2 id=\"title\" class=\"titlel2std\">" + chapterTitle + "</h2>"
		content = strings.Replace(content, "　　", "<p class=\"a\">　　", -1)
		content = strings.Replace(content, "\n", "</p>", -1)
		// chapterTitles = append(chapterTitles, chapterTitle)
		book.chapters.Store(chapterID, chapterTitle)
		return chapterTitle, contentHeader + "\n" + titleElement + "\n" + content + "\n" + contentFooter, true
	}
	return chapterTitle, chapterTitle + "\n\n" + content + "\n\n\n\n", true
}

//下载索引
var di int

//写出章节缓存
func writeChapterTemp(chapterID string) {
	var tmpPath string
	var fileName string
	if book.format == "epub" {
		tmpPath = book.tmpPath + "OEBPS" + pathSeparator
		fileName = "chapter" + chapterID + ".html"
	} else {
		tmpPath = book.tmpPath
		fileName = chapterID + ".txt"
	}
	_, content, result := getChapterContent(chapterID)
	if result {
		writeOut(content, tmpPath, fileName)
	}
	bar.Add(1)
	if di == book.chapterNum {
		quit <- 1
	}
	di++
}

//下载章节
func downloadChapters() {
	di = 1
	bar = *progressbar.New(len(book.chapterIDs))
	for _, chapterID := range book.chapterIDs {
		go writeChapterTemp(chapterID.String())
	}
	<-quit
	if book.format == "epub" {
		println(pathSeparator)
		coverElement := coverHeader + "\n" + "<img src=\"cover.jpg\" alt=\"" + book.name + "\" />" + coverFooter
		coverBody := httpGet(book.coverURL, nil)
		writeOut(mimetype, book.tmpPath, "mimetype")
		writeOut(container, book.tmpPath+"META-INF"+pathSeparator, "container.xml")
		writeOut(coverElement, book.tmpPath+"OEBPS"+pathSeparator, "cover.html")
		writeOut(coverBody, book.tmpPath+"OEBPS"+pathSeparator, "cover.jpg")
		writeOut(css, book.tmpPath+"OEBPS"+pathSeparator, "style.css")
		writeOut(genBookToc(), book.tmpPath+"OEBPS"+pathSeparator, "book-toc.html")
		writeOut(genContentOpf(), book.tmpPath+"OEBPS"+pathSeparator, "content.opf")
		writeOut(genTocNcx(), book.tmpPath+"OEBPS"+pathSeparator, "toc.ncx")
		compressEpub(book.tmpPath, path.out+pathSeparator+book.name+".epub")
	} else {
		mergeTemp()
	}
	bar.Finish()
}

func login() {
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
}
