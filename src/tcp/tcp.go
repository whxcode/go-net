package net

import (
	"fmt"
	"net"
	"strings"
)

func procces(conn net.Conn) {
	defer func() {
		fmt.Println("Closing connection from:", conn.RemoteAddr())
		conn.Close()
	}() // 处理完毕关闭

	// 接收数据
	buf := make([]byte, 1024)
	for {

		fmt.Println("Waiting for data...")

		n, err := conn.Read(buf)

		if n == 0 {
			fmt.Println("Connection closed by client")
			return
		}

		fmt.Println("read n:", n)

		if err != nil {
			fmt.Println("Error reading data:", err)
			return
		}

		msg := string(buf[:n])

		MSG := strings.ToUpper(msg)

		n, err = conn.Write([]byte(MSG))
		if err != nil {
			fmt.Println("Error writing data:", err)
		}

		fmt.Println("Received data:")
	}
}

func Test() {
	listner, err := net.Listen("tcp", "0.0.0.0:8081")
	println("Server started on")

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	for {
		fmt.Println("Waiting for connection...")
		conn, err := listner.Accept()
		println("Connection accepted from:", conn.RemoteAddr())

		if err != nil {
			fmt.Println("Error accepting connection:", err)
		}

		go procces(conn)

	}
}
