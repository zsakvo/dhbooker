package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	pb "gopkg.in/cheggaaa/pb.v1"
)

var downloadIndex int
var downIndexChan (chan int)
var downloadFailedChan (chan string)
var channelClosed = false

//获取书籍信息
func getBookInfo() {
	if len(book.id) == 2 {
		println("请至少输入书籍 ID")
		os.Exit(1)
	}
	bookInfoURL := "https://app.hbooker.com/book/get_info_by_id"
	paramsMap := map[string]string{"app_version": appVersion, "login_token": token.loginToken, "book_id": book.id, "account": token.account}
	res, err := httpGet(bookInfoURL, paramsMap)
	if err != nil {
		panic(err)
	}
	body, err := getBody(res)
	check(err)
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
	chapterListURL := "https://www.ciweimao.com/chapter-list/" + book.id + "/book_detail"
	chapterListRes, err := httpGet(chapterListURL, nil)
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(chapterListRes.Body)
	if err != nil {
		panic(err)
	}
	doms := doc.Find("ul[class=book-chapter-list]>li>a")
	book.chapterNum = doms.Length()
	doms.Each(func(i int, selection *goquery.Selection) {
		attr, _ := selection.Attr("href")
		id := strings.Replace(attr, "https://www.ciweimao.com/chapter/", "", -1)
		book.chapterIDs = append(book.chapterIDs, id)
	})
	fmt.Println("《" + book.name + "》，" + "共" + strconv.Itoa(book.chapterNum) + "章")
}

//下载章节
func downloadChapters() {
	downloadIndex = 0
	downIndexChan = make(chan int)
	downloadFailedChan = make(chan string)
	bar = pb.StartNew(book.chapterNum)
	bar.ShowTimeLeft = false
	// bar.SetWidth(80)
	for _, chapterID := range book.chapterIDs {
		go writeChapterTemp(chapterID, book.chapterNum)
	}

	for {
		if channelClosed {
			return
		}
		select {
		case ix, _ := <-downIndexChan:
			if ix < 2 {
				bar.Increment()
			}
			downloadIndex++

		case df, _ := <-downloadFailedChan:
			download.failedChapters = append(download.failedChapters, df)
		}
		watchChan(book.chapterNum)
	}
}

func redownlodChapters() {
	downloadIndex = 0
	chapterIDs := download.failedChapters
	download.failedChapters = download.failedChapters[:0]
	for _, chapterID := range chapterIDs {
		go writeChapterTemp(chapterID, len(chapterIDs))
	}

	for {
		if channelClosed {
			return
		}
		select {
		case ix, _ := <-downIndexChan:
			if ix < 2 {
				bar.Increment()
			}
			downloadIndex++

		case df, _ := <-downloadFailedChan:
			download.failedChapters = append(download.failedChapters, df)

		}
		watchChan(len(chapterIDs))
	}
}

func watchChan(num int) {
	if channelClosed {
		println("关了")
		return
	}
	if downloadIndex == num {
		if len(download.failedChapters) == 0 {
			close(downIndexChan)
			close(downloadFailedChan)
			channelClosed = true
			genBook()
		}
		redownlodChapters()
	}
}

//获取章节内容
func getChapterContent(chapterID string) (string, int) {
	contentKeyURL := "https://app.hbooker.com/chapter/get_chapter_cmd?app_version=" + appVersion + "&chapter_id=" + chapterID + "&login_token=" + token.loginToken + "&account=" + token.account
	paramsMap := map[string]string{"app_version": appVersion, "chapter_id": chapterID, "login_token": token.loginToken, "account": token.account}
	keyRes, err := httpGet(contentKeyURL, paramsMap)
	if err != nil {
		return "", 2
	}
	keyBody, err := getBody(keyRes)
	if err != nil {
		return "", 2
	}
	keyBody = decode(keyBody, initEncryptKey)
	contentKey := gjson.Get(keyBody, "data.command").String()
	contentURL := "https://app.hbooker.com/chapter/get_cpt_ifm"
	paramsMap1 := map[string]string{"chapter_command": contentKey, "app_version": appVersion, "login_token": token.loginToken, "chapter_id": chapterID, "account": token.account}
	contentRes, err := httpGet(contentURL, paramsMap1)
	if err != nil {
		return "", 2
	}
	contentBody, err := getBody(contentRes)
	if err != nil {
		return "", 2
	}
	contentBody = decode(contentBody, initEncryptKey)
	chapterTitle := gjson.Get(contentBody, "data.chapter_info.chapter_title").String()
	content := gjson.Get(contentBody, "data.chapter_info.txt_content").String()
	auth := gjson.Get(contentBody, "data.chapter_info.auth_access").String()
	if auth == "0" {
		book.invalidChapters.Store(chapterID, nil)
		return "", 1
	}
	content = decode(content, contentKey)
	if book.format == "epub" {
		titleElement := "<h2 id=\"title\" class=\"titlel2std\">" + chapterTitle + "</h2>"
		content = strings.Replace(content, "　　", "<p class=\"a\">　　", -1)
		content = strings.Replace(content, "\n", "</p>", -1)
		// chapterTitles = append(chapterTitles, chapterTitle)
		book.chapters.Store(chapterID, chapterTitle)
		return contentHeader + "\n" + titleElement + "\n" + content + "\n" + contentFooter, 0
	}
	return chapterTitle + "\n\n" + content + "\n\n\n\n", 0
}

//写出章节缓存
func writeChapterTemp(chapterID string, num int) {
	var tmpPath string
	var fileName string
	if book.format == "epub" {
		tmpPath = book.tmpPath + "OEBPS" + pathSeparator
		fileName = "chapter" + chapterID + ".html"
	} else {
		tmpPath = book.tmpPath
		fileName = chapterID + ".txt"
	}
	content, code := getChapterContent(chapterID)
	switch {
	case code == 0:
		writeOut(content, tmpPath, fileName)
	case code == 1:
		break
	case code == 2:
		downloadFailedChan <- chapterID
	}
	downIndexChan <- code
}

func genBook() {
	if book.format == "epub" {
		coverElement := coverHeader + "\n" + "<img src=\"cover.jpg\" alt=\"" + book.name + "\" />" + coverFooter
		res, err := httpGet(book.coverURL, nil)
		if err != nil {
			panic(err)
		}
		coverBody, err := getBody(res)
		check(err)
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
