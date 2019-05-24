package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Unknwon/goconfig"
)

//读取配置文件
func getConfig() {
	if isFileExist("conf.ini") {
		cfg, err := goconfig.LoadConfigFile("conf.ini")
		if err != nil {
			os.Remove("conf.ini")
			initSettings()
		}
		config = cfg
	} else {
		initSettings()
		fmt.Println("已初始化 conf.ini 配置文件，请按提示填写后运行本程序")
		os.Exit(0)
	}
}

//获取 section
func getSection(value string) map[string]string {
	account, err := config.GetSection(value)
	if err != nil {
		initSettings()
	}
	return account
}

func getAccountSettings() accountSettings {
	accountMap := getSection("account")
	var aStruct accountSettings
	aStruct.username = accountMap["username"]
	aStruct.password = accountMap["password"]
	if len(aStruct.username) == 0 {
		fmt.Println("请于 conf.ini 中填写用户名")
		os.Exit(1)
	}
	if len(aStruct.password) == 0 {
		fmt.Println("请于 conf.ini 中填写密码")
		os.Exit(1)
	}
	return aStruct
}

func getPathSettings() {
	pathMap := getSection("path")
	var aStruct pathSettings
	aStruct.tmp = pathMap["tmp"]
	aStruct.out = pathMap["out"]
	if len(aStruct.tmp) == 0 {
		fmt.Println("请于 conf.ini 中填写临时目录路径")
		os.Exit(1)
	}
	if len(aStruct.out) == 0 {
		fmt.Println("请于 conf.ini 中填写输出目录路径")
		os.Exit(1)
	}
	path = aStruct
}

func getToken() string {
	tokenMap := getSection("token")
	return tokenMap["token"]
}

func getMobi() {
	mobiMap := getSection("mobi")
	hintBool, err := strconv.ParseBool(mobiMap["hint"])
	if err != nil {
		mobi.hint = true
	} else {
		mobi.hint = hintBool
	}
	mobi.caliPath = mobiMap["calibre-path"]
}

func initSettings() {
	os.Create("conf.ini")
	cfg, err := goconfig.LoadConfigFile("conf.ini")
	check(err)
	cfg.SetValue("account", "username", "")
	cfg.SetKeyComments("account", "username", "# 用户名，必填")
	cfg.SetValue("account", "password", "")
	cfg.SetKeyComments("account", "password", "# 密码，必填")
	cfg.SetValue("token", "token", "")
	cfg.SetKeyComments("token", "token", "# 自动生成，请勿修改")
	cfg.SetValue("path", "tmp", "tmp")
	cfg.SetKeyComments("path", "tmp", "# 临时目录，必填")
	cfg.SetValue("path", "out", "output")
	cfg.SetKeyComments("path", "out", "# 输出目录，必填")
	cfg.SetValue("mobi", "hint", "true")
	cfg.SetKeyComments("mobi", "hint", "# 首次提示，只能为 true 或 false")
	cfg.SetValue("mobi", "calibre-path", "")
	cfg.SetKeyComments("mobi", "calibre-path", "# calibre 路径，请精确到 ebook-convert 可执行文件所在的目录")
	err1 := goconfig.SaveConfigFile(cfg, "conf.ini")
	check(err1)
}

func setSeparator() {
	if ostype == "windows" {
		pathSeparator = "\\"
	} else if ostype == "linux" {
		pathSeparator = "/"
	} else if ostype == "darwin" {
		pathSeparator = "/"
	}
}

func initConfig() {
	getConfig()
	getPathSettings()
	setSeparator()
	getMobi()
}
