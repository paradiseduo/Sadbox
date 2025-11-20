# Sadbox
专用于macOS的应用沙盒清理工具

```bash
git clone https://github.com/paradiseduo/Sadbox.git
cd Sadbox
go build
chmod +x Sadbox
./Sadbox -h

  -delete string
    	要删除的文件名（会删除包含该文件的容器目录），可以指定多个文件名，用空格分隔
  -system
    	显示以 com.apple. 开头的系统文件
```

删除多个文件：
```bash
./Sadbox -delete "com.taobao.taobao4iphone com.jd.store"
```

显示系统应用沙盒：
```bash
./Sadbox -system
```
