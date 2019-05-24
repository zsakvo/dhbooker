package main

import (
	"fmt"

	"github.com/tidwall/gjson"
)

func login() {
	// fmt.Println("正在检查凭证")
	token := getToken()
	if len(token) == 0 {
		// fmt.Println("无凭证，尝试使用用户密码登入")
		account := getAccountSettings()
		loginByPass(account)
	} else {
		// fmt.Println("尝试使用上次凭证登入")
		localToken := getSection("token")["token"]
		loginByToken(localToken)
	}
}

//凭证登入
func loginByToken(localToken string) {
	token.readerID = gjson.Get(localToken, "reader_info.reader_id").String()
	token.loginToken = gjson.Get(localToken, "login_token").String()
	token.account = gjson.Get(localToken, "reader_info.account").String()
	url := "https://app.hbooker.com/reader/get_my_info"
	paramsMap := map[string]string{"reader_id": token.readerID, "app_version": appVersion, "login_token": token.loginToken, "account": token.account}
	res, err := httpGet(url, paramsMap)
	if err != nil {
		panic(err)
	}
	body, err := getBody(res)
	check(err)
	body = decode(body, initEncryptKey)
	code := gjson.Get(body, "code").String()
	if code != "100000" {
		// fmt.Println(gjson.Get(body, "tip"))
		// fmt.Println("凭证登入失败，尝试使用用户密码登入")
		account := getAccountSettings()
		loginByPass(account)
	} else {
		// println("登入成功，准备下载……")
	}
}

//账号登入
func loginByPass(account accountSettings) {
	name := account.username
	pass := account.password
	url := "https://app.hbooker.com/signup/login"
	paramsMap := map[string]string{"login_name": name, "app_version": appVersion, "passwd": pass}
	res, err := httpGet(url, paramsMap)
	if err != nil {
		panic(err)
	}
	body, err := getBody(res)
	check(err)
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
		// println("登入成功，准备下载……")
	}
}
