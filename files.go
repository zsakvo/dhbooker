package main

import (
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
func mergeTemp() bool {
	var content []byte
	if !isFileExist(path.out) {
		os.MkdirAll(path.out, os.ModePerm)
	}
	outPath := path.out + "/" + bookName + ".txt"
	tmpPath := path.tmp + "/" + bookName + "/"
	for _, cid := range chapterIDs {
		d, err := ioutil.ReadFile(tmpPath + cid.String() + ".txt")
		if err == nil {
			content = append(content, d...)
		}
	}
	err1 := ioutil.WriteFile(outPath, content, 0644)
	if err1 != nil {
		return false
	}
	return true
}

//移除缓存
func destoryTemp() error {
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
	return nil
}
