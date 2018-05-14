
# disk_spider

一款简单的磁盘文件的遍历、打包、加密工具
分为 archiver 打包工具 与 unarchiver 解包工具

* [windows-archiver-64-1.0.0](https://github.com/naxiemolv/disk_spider/releases/download/1.0.0/archiver-64.exe)
* [windows-unarchiver-64-1.0.0](https://github.com/naxiemolv/disk_spider/releases/download/1.0.0/unarchiver-64.exe)
* [macOS-archiver-64-1.0.0](https://github.com/naxiemolv/disk_spider/releases/download/1.0.0/macOS-archiver)
* [macOS-unarchiver-64-1.0.0](https://github.com/naxiemolv/disk_spider/releases/download/1.0.0/macOS-unarchiver)

#配置文件
```
{
  "dir_path": [
    "C:\\Users"
  ],
  "suffix": [
    "gif"
  ],
  "name_contain": [],
  "max_file_count": 10000,
  "max_file_size": 51200,
  "max_output_size": 1024000
}
```
dir_path：为目标目录的数组

suffix： 所需要的文件后缀

name_contain：文件名包含（暂未实现）

max_file_count：最大文件数量

max_file_size：最大文件大小

max_output_size： 输出文件大小


上述配置将扫描 C:\Users 文件夹下所有的 gif文件，最多不超过10000个，最大单一文件不超过51.2 M 输出文件不超过 1024 M

#使用说明


* 打包
    
    1.创建如上的配置文件命名为 config.json

    2.置 archiver与config 于同目录下

    3.运行 archiver 程序将生成 arc.x 打包文件


* 解包

    1.置 unarchiver 与 arc.x 于同目录下
    
    2.运行 unarchiver



