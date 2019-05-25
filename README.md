# dhbooker

`本程序是命令行工具，请在终端内运行，不要直接双击！`

## 简介：

刺猬猫（欢乐书客）小说网站书籍下载工具。

- 程序本体使用 golang 编写，所以编译后体积较大

- 使用其 Android 客户端接口，所以需要自行配置账户密码

- 支持已购买的付费章节内容获取

- 支持导出格式：txt，epub

- 支持自定义导出目录

- 支持自定义超时时间（单位 ms）

## 使用方法

- 程序需要在终端（terminal/powershell/cmd）内运行，请不要直接双击运行。

- 第一次使用时请先直接运行一下程序初始化配置文件 `conf.ini` ，或者直接下载 [此文件]("https://raw.githubusercontent.com/zsakvo/dhbooker/master/conf.example.ini") 放入和程序的同级目录下。

- 按提示填写好 `conf.ini` 中的必要字段（用户名 & 密码）

- 在当前目录启动终端，执行

```
./dhbooker -b bookid -t type -p timeout

#dhbooker 代表具体的可执行文件名，请酌情修改

#bookid 可以在书籍网页 url 中找到

#type 有两种，分为 txt 以及 epub，如果不填写则默认 txt

#timeout 为超时时间，单位毫秒，默认值 5000
```

- 等待提示下载完毕即可

## 关于 mobi 格式

经评估后感觉没有必要专门添加，如果有需要，请自行安装 Calibre，然后使用本程序下载对应书籍的 epub 以后使用 `ebook-convert xxx.epub xxx.mobi` 命令自行转换

## 注意事项

- 程序运行时会清空缓存目录，请务必不要将其它文件放入其中

- 本程序仅支持获取免费与您已付费订阅的内容，不具备任何白嫖付费章节的作用
