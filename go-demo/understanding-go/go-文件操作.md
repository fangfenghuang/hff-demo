
# 判断文件是否存在/创建空文件
```
emptyFile, err := os.Create("empty.txt")
defer emptyFile.Close()
```
```
_, err := os.Stat("test")
 
if os.IsNotExist(err) {
    errDir := os.MkdirAll("test", 0755)
    if errDir != nil {
        log.Fatal(err)
    }

}
```
# 重命名文件: os.Rename
```
oldName := "test.txt"
newName := "testing.txt"
err := os.Rename(oldName, newName)
```

# 移动文件: os.Rename
```
oldLocation := "/var/www/html/test.txt"
newLocation := "/var/www/html/src/test.txt"
err := os.Rename(oldLocation, newLocation)
```

# 复制文件: io.Copy
```
sourceFile, err := os.Open("/var/www/html/src/test.txt")
defer sourceFile.Close()
newFile, err := os.Create("/var/www/html/test.txt")
defer newFile.Close()
bytesCopied, err := io.Copy(newFile, sourceFile)
```

# 获取文件的metadata信息: os.Stat()
```
fileStat, err := os.Stat("test.txt")
fmt.Println("File Name:", fileStat.Name())        // Base name of the file
fmt.Println("Size:", fileStat.Size())             // Length in bytes for regular files
fmt.Println("Permissions:", fileStat.Mode())      // File mode bits
fmt.Println("Last Modified:", fileStat.ModTime()) // Last modification time
fmt.Println("Is Directory: ", fileStat.IsDir())   // Abbreviation for Mode().IsDir()
```

# 删除文件: os.Remove()

# 读取文件字符: bufio.NewScanner()
```
filebuffer, err := ioutil.ReadFile(filename)
inputdata := string(filebuffer)
data := bufio.NewScanner(strings.NewReader(inputdata))
data.Split(bufio.ScanRunes)
for data.Scan() {
    fmt.Print(data.Text())
}
```

# 清除文件: os.Truncate() 
裁剪一个文件到100个字节.
如果文件本来就少于100个字节,则文件中原始内容得以保留,剩余的字节以null字节填充.
如果文件本来超过100个字节,则超过的字节会被抛弃.
这样我们总是得到精确的100个字节的文件.
传入0则会清空文件.

# 向文件添加内容
```
f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
fmt.Fprintf(f, "%s\n", message)
```

# 修改文件权限,时间戳
```
err = os.Chmod("test.txt", 0777)
err = os.Chown("test.txt", os.Getuid(), os.Getgid())
err = os.Chtimes("test.txt", lastAccessTime, lastModifyTime)
```

# zip操作 "archive/zip"
- 压缩文件到ZIP格式
- 读取ZIP文件里面的文件
- 解压ZIP文件
