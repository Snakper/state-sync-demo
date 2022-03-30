package src

import (
	"encoding/json"
	"net"
)

var c *net.Conn

func Listen() {
	go func() {
		l, err := net.Listen("tcp", "127.0.0.1:9123")
		if err != nil {
			panic(err)
		}
		for {
			conn, err := l.Accept()
			if err != nil {
				panic(err)
			}
			c = &conn
		}
	}()
}

func SendToNetwork(msg ControlMsg) {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	if c != nil {
		_, err = (*c).Write(b)
		if err != nil {
			panic(err)
		}
	}
}
