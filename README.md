# txtbot

Handy tool.

## Run

```go
package main

import (
	"github.com/FengliangChen/txtbot"
	)

func main() {
	txtbot.Run()
}
```

## 安装
Mac打开终端：
```bash
$(curl https://raw.githubusercontent.com/FengliangChen/py/master/release/txtbot -o /usr/local/bin/txtbot) && cd /usr/local/bin && chmod +x txtbot && cd $HOME 

```
安装字典：
```bash
mkdir $HOME/Documents/txtbot && $(curl https://raw.githubusercontent.com/FengliangChen/py/master/release/clientcode.json -o $HOME/Documents/txtbot/clientcode.json)
```

## 运行：
默认 生成文字 + 打开文件夹：
```bash
txtbot 190707
```
或
```bash
txtbot 190707
```

## 取消打开文件夹：
```bash
txtbot -c
```
或
```bash
txtbot -c 190707
```

## 取消PF

对未开始和进行中都没有放有单号的，默认运行停止，可用-t取消停止。

```bash
txtbot -t 190707
```
或
```bash
txtbot -c -t 190707
```

## 帮助：
```bash
txtbot -h
```
