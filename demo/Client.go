package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ...")
	time.Sleep(time.Second * 1)

	conn, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("Dial TCP addr error: ", err)
		return
	}

	for {
		_, err := conn.Write([]byte("Hello Zinx v0.2 !"))
		if err != nil {
			fmt.Println("client write error: ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("client read error: ", err)
			return
		}

		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)

		time.Sleep(time.Second * 1)
	}
}
