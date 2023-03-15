package main

import (
	"os"

	"github.com/dcnetio/dcevauth/command"
)

func main() {
	//读取命令行参数，并解析响应
	if len(os.Args) == 1 { //显示帮助
		command.ShowHelp()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "--config": //配置助记词，一次配置，一个月有效
		command.ConfigDeal()
	case "--sign": //对输入的enclaveid进行签名
		command.SignDeal()
	case "--signer": //显示签名人的公钥
		command.ShowSigner()
	default:
		command.ShowHelp()
	}
	os.Exit(1)
}
