// 命名空间 thrift_gen 对应 go 代码的生成目录
namespace go thrift_gen

// Data 结构体
struct Data {
    1: string text // 第一个属性，标号为1 string 类型 名称是 text
}

// 定义一个服务 名称是 format_data
service format_data {
    // 服务的方法 Data 是返回值 do_format 是名称　1:Data data 表示第一个参数类型是 Data 名为data
    Data do_format(1:Data data),
}