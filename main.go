package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dcnetio/dcevauth/command"
)

//go:embed  "version.json"
var configVersion []byte
var GetVersion = func() (verStr string) {
	var version struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(configVersion, &version); err != nil {
		fmt.Println("unmarshal version.json error:", err)
	}
	verStr = version.Version
	return
}()

func main() {
	fmt.Println("dcevauth version:", GetVersion)
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
