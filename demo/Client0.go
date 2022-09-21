package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/x627234313/zinx-test/znet"
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
		// 封包， 发送Message
		dp := znet.NewDataPack()
		msg, _ := dp.Pack(znet.NewMessage(0, []byte("Zinx V0.6 Client-0 Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("client write error: ", err)
			return
		}

		// 先读出流中head
		head := make([]byte, dp.GetHead())
		if _, err := io.ReadFull(conn, head); err != nil {
			fmt.Println("read head error")
			break
		}

		// 将headData拆包到msg中
		msgHead, err := dp.Unpack(head)
		if err != nil {
			fmt.Println("server unpack error:", err)
			return
		}

		if msgHead.GetMsgDataLen() > 0 {
			// msg 中有data数据，再次读取
			msg := msgHead.(*znet.Message)
			data := make([]byte, msg.GetMsgDataLen())

			if _, err := io.ReadFull(conn, data); err != nil {
				fmt.Println("server unpack data error:", err)
				return
			}

			msg.SetMsgData(data)

			fmt.Println("==> Recv Msg ID = ", msg.GetMsgId(), "Len = ", msg.GetMsgDataLen(), "Data = ", string(msg.GetMsgData()))
		}

		time.Sleep(time.Second * 1)
	}
}
