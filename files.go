package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func isFileExist(filePath string) bool {
	_, err := os.Stat(filePath) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

//合并缓存
func mergeTemp() {
	var content []byte
	if !isFileExist(path.out) {
		os.MkdirAll(path.out, os.ModePerm)
	}
	outPath := path.out + pathSeparator + book.name + ".txt"
	tmpPath := book.tmpPath
	for _, cid := range book.chapterIDs {
		d, err := ioutil.ReadFile(tmpPath + cid + ".txt")
		if err == nil {
			content = append(content, d...)
		}
	}
	err1 := ioutil.WriteFile(outPath, content, 0644)
	check(err1)
}

//移除缓存
func destoryTemp(b bool) error {
	if b {
		fmt.Println("正在清理临时目录")
	}
	d, err := os.Open(path.tmp)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path.tmp, name))
		if err != nil {
			return err
		}
	}
	if !b {
		fmt.Print("\n下载完毕")
	}
	return nil
}

//写出内容
func writeOut(content, dirPath, fileName string) {
	if !isFileExist(dirPath) {
		err := os.MkdirAll(dirPath, os.ModePerm)
		check(err)
	}
	outPath := dirPath + fileName
	d := []byte(content)
	err := ioutil.WriteFile(outPath, d, 0644)
	check(err)
}
