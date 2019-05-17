package main

import (
	"strconv"
	"strings"
)

func genContentOpf() string {
	var manifestStr string
	var spineStr string
	for i, chapterID := range chapterIDs {
		_, ok := invalidChapters.Load(chapterID.String())
		if !ok {
			manifestStr += "<item id=\"chapter" + strconv.Itoa(i) + "\" href=\"chapter" + chapterID.String() + ".html\" media-type=\"application/xhtml+xml\"/>"
			spineStr += "<itemref idref=\"chapter" + strconv.Itoa(i) + "\" linear=\"yes\"/>"
		}
	}
	contentOpfHeader = strings.Replace(contentOpfHeader, "bookTitle", bookName, 1)
	contentOpfHeader = strings.Replace(contentOpfHeader, "bookAuthor", bookAuthor, 1)
	return contentOpfHeader + contentOpfManifestHeader + manifestStr + contentOpfManifestFooter + contentOpfNcxtocHeader + spineStr + contentOpfNcxtocFooter + contentOpfFooter
}

func genTocNcx() string {
	docTitle := "<docTitle><text>" + bookName + "</text></docTitle>"
	docAuthor := "<docAuthor><text>" + bookAuthor + "</text></docAuthor>"
	navMap := "<navMap> <navPoint id=\"cover\" playOrder=\"1\"> <navLabel><text>封面</text></navLabel> <content src=\"cover.html\"/> </navPoint> <navPoint id=\"htmltoc\" playOrder=\"2\"> <navLabel><text>目录</text></navLabel> <content src=\"book-toc.html\"/> </navPoint>\""
	var str string
	for i, id := range chapterIDs {
		title, ok := chapters.Load(id.String())
		if ok {
			str += "<navPoint id=\"chapter" + strconv.Itoa(i) + "\" playOrder=\"" + strconv.Itoa(3+i) + "\"> <navLabel><text>" + title.(string) + "</text></navLabel> <content src=\"chapter" + chapterIDs[i].String() + ".html\"/> </navPoint>"
		}
	}
	return tocNcxHeader + docTitle + docAuthor + navMap + str + tocNcxFooter
}

func genBookToc() string {
	var str string
	for _, id := range chapterIDs {
		title, ok := chapters.Load(id.String())
		if ok {
			str += "<dt class=\"tocl2\"><a href=\"chapter" + id.String() + ".html\">" + title.(string) + "</a></dt>"
		}
	}
	return bookTocHeader + str + bookTocFooter
}
