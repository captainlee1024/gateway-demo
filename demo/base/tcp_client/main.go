package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// 1. 连接服务器
	conn, err := net.Dial("tcp", "0.0.0.0:9090")
	//defer conn.Close() // 这里不关闭会出现什么问题
	if err != nil {
		fmt.Printf("connect failed, err: %v\n", err)
		return
	}

	// 2. 读取命令行输入
	inputReader := bufio.NewReader(os.Stdin)
	for {
		// 3. 一直读取到 \n
		input, err := inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("read from console failed, err: %v\n", err)
			break
		}

		// 4. 读取 Q 时停止
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "Q" {
			break
		}

		// 5. 回复服务器信息
		_, err = conn.Write([]byte(trimmedInput))
		if err != nil {
			fmt.Printf("write failed, err: %v\n", err)
			break
		}
	}

	fmt.Printf("client close")
	// 这里先不让它退出程序，睡眠一会儿，看看在睡眠期间连接有没有释放
	time.Sleep(100 * time.Second)
}
