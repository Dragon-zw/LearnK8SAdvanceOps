package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	// 使用 -s 参数来控制输出的字符串，默认值为 "hello go container"
	message := flag.String("s", "hello go container", "要输出的字符串")
	flag.Parse()

	// 每隔一秒钟就打印指定的字符串
	// 可以使用 Ctrl + C 来停止程序
	fmt.Printf("开始每隔 1 秒钟打印: %s\n", *message)
	for {
		fmt.Println(*message)
		time.Sleep(1 * time.Second)
	}
}