package main

import "github.com/x627234313/zinx-test/znet"

func main() {
	s := znet.NewServer("Zinx v0.1")

	s.Serve()
}
