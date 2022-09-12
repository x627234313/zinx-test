package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	fmt.Println("Client Test ... Start")
	time.Sleep(time.Second * 1)

	conn, err := net.Dial("tcp", "0.0.0.0:9999")
	if err != nil {
		fmt.Println("Dial TCP addr error: ", err)
		return
	}

	for {
		_, err := conn.Write([]byte("Hello Zinx v0.3!"))
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

		fmt.Printf("Server Call Back:\n%s, cnt = %d\n", buf, cnt)

		time.Sleep(time.Second * 1)
	}
}
