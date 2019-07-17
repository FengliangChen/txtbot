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

## 帮助：
```bash
txtbot -h
```
