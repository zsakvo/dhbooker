package main

import (
	"strconv"
	"strings"
)

func genEpub() {
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
}

func genContentOpf() string {
	var manifestStr string
	var spineStr string
	for i, chapterID := range book.chapterIDs {
		_, ok := book.invalidChapters.Load(chapterID)
		if !ok {
			manifestStr += "<item id=\"chapter" + strconv.Itoa(i) + "\" href=\"chapter" + chapterID + ".html\" media-type=\"application/xhtml+xml\"/>"
			spineStr += "<itemref idref=\"chapter" + strconv.Itoa(i) + "\" linear=\"yes\"/>"
		}
	}
	contentOpfHeader = strings.Replace(contentOpfHeader, "bookTitle", book.name, 1)
	contentOpfHeader = strings.Replace(contentOpfHeader, "bookAuthor", book.author, 1)
	return contentOpfHeader + contentOpfManifestHeader + manifestStr + contentOpfManifestFooter + contentOpfNcxtocHeader + spineStr + contentOpfNcxtocFooter + contentOpfFooter
}

func genTocNcx() string {
	docTitle := "<docTitle><text>" + book.name + "</text></docTitle>"
	docAuthor := "<docAuthor><text>" + book.author + "</text></docAuthor>"
	navMap := "<navMap> <navPoint id=\"cover\" playOrder=\"1\"> <navLabel><text>封面</text></navLabel> <content src=\"cover.html\"/> </navPoint> <navPoint id=\"htmltoc\" playOrder=\"2\"> <navLabel><text>目录</text></navLabel> <content src=\"book-toc.html\"/> </navPoint>\""
	var str string
	for i, id := range book.chapterIDs {
		title, ok := book.chapters.Load(id)
		if ok {
			str += "<navPoint id=\"chapter" + strconv.Itoa(i) + "\" playOrder=\"" + strconv.Itoa(3+i) + "\"> <navLabel><text>" + title.(string) + "</text></navLabel> <content src=\"chapter" + book.chapterIDs[i] + ".html\"/> </navPoint>"
		}
	}
	return tocNcxHeader + docTitle + docAuthor + navMap + str + tocNcxFooter
}

func genBookToc() string {
	var str string
	for _, id := range book.chapterIDs {
		title, ok := book.chapters.Load(id)
		if ok {
			str += "<dt class=\"tocl2\"><a href=\"chapter" + id + ".html\">" + title.(string) + "</a></dt>"
		}
	}
	return bookTocHeader + str + bookTocFooter
}
