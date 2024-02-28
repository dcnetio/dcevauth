package command

import (
	"crypto/ed25519"
	"fmt"
	"os"
	"strings"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/dcnetio/go-substrate-rpc-client/types/codec"
	"github.com/edgelesssys/ego/ecrypto"
	"github.com/libp2p/go-libp2p/core/crypto"
	mbase "github.com/multiformats/go-multibase"
)

func ShowHelp() {
	fmt.Println("usage: dcevauth  [OPTIONS] [args]")
	fmt.Println("Options:")
	fmt.Println("")
	fmt.Println(" --config string              config mnemonic for signer")
	fmt.Println(" --sign string                create signature for input string,When the input")
	fmt.Println("                              string starts with '0x', it will be decoded in")
	fmt.Println("                              hexadecimal before signing")
	fmt.Println(" --signer                     show signer publickey coded with mbase")
}

func ConfigDeal() {
	if len(os.Args) < 3 {
		ShowHelp()
		return
	}
	mnemonic := os.Args[2]
	//tee 封装mnemonic，保存到文件中
	sealedMnemonicBytes, err := ecrypto.SealWithUniqueKey([]byte(mnemonic), nil)
	if err != nil {
		fmt.Println("failed to seal secret")
		return
	}
	//保存到文件中
	err = os.WriteFile("/opt/dcnetio/data/.mnemonic", sealedMnemonicBytes, 0644)
	if err != nil {
		fmt.Println("failed to write secret to file,try with sudo or root")
		return
	}
	_, codePubkey, err := loadPrivkey()
	if err != nil {
		fmt.Println("config  error:", err)
		return
	}
	fmt.Println("config success,public key:", codePubkey)
}

func SignDeal() {
	if len(os.Args) < 3 {
		ShowHelp()
		return
	}
	var needSign []byte
	var err error
	if strings.HasPrefix(os.Args[2], "0x") || strings.HasPrefix(os.Args[2], "0X") {
		needSign, err = codec.HexDecodeString(os.Args[2])
		if err != nil {
			fmt.Println("invalid input, error:", err)
			return
		}
	} else {
		needSign = []byte(os.Args[2])
	}

	privKey, codedPubkey, err := loadPrivkey()
	if err != nil {
		fmt.Println("load private key failed,please confirm mnemonic configed, error:", err)
		return
	}
	signature, err := privKey.Sign(needSign)
	if err != nil {
		fmt.Println("sign enclaveid failed, error:", err)
		return
	}
	encodedSignature, err := mbase.Encode(mbase.Base32, signature)
	if err != nil {
		fmt.Println("encode signature failed, error:", err)
		return
	}
	//打印签名结果
	fmt.Println("input:", os.Args[2])
	fmt.Println("public key:", codedPubkey)
	fmt.Println("signature:", encodedSignature)
}

// 获取程序的运行状态
func ShowSigner() {
	_, codePubkey, err := loadPrivkey()
	if err != nil {
		fmt.Println("load private key failed,please confirm mnemonic configed, error:", err)
		return
	}
	fmt.Println("public key:", codePubkey)
}

func loadPrivkey() (privkey crypto.PrivKey, codedPubkey string, err error) {
	//从文件中读取mnemonic
	sealedMnemonicBytes, err := os.ReadFile("/opt/dcnetio/data/.mnemonic")
	if err != nil {
		return
	}
	//解封
	mnemonic, err := ecrypto.Unseal(sealedMnemonicBytes, nil)
	if err != nil {
		return
	}
	//导入私钥
	privkey, err = ImportTestMnemonic(string(mnemonic))
	if err != nil {
		return
	}
	pubkey := privkey.GetPublic()
	pub, err := pubkey.Raw()
	if err != nil {
		return
	}
	codedPubkey, _ = mbase.Encode(mbase.Base32, pub)
	return
}

func ImportTestMnemonic(mnemonic string) (prikey crypto.PrivKey, err error) {
	seed, err := schnorrkel.SeedFromMnemonic(mnemonic, "")
	if err != nil {
		return
	}

	secret := ed25519.NewKeyFromSeed(seed[:32])
	prikey, err = crypto.UnmarshalEd25519PrivateKey([]byte(secret))
	if err != nil {
		return
	}
	return
}
