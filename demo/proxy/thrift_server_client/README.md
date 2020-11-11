## 首先安装thrift，安装步骤请参照 [官网](https://thrift.apache.org/docs/install/)
manjaro　安装：

````bash
sudo pacman -S thrift
````

## 构建 thrift 测试 server 与 client
- 首先编写 IDL 就是 thrift 的配置文件`thrift_gen.thrift`
- 运行 IDL 生成命令 `thrift --gen go thrift_gen.thrift`
- 使用生成的 IDL 单独构建 server 与 client 即可