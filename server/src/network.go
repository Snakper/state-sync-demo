package src

import (
	"encoding/json"
	"net"
)

var c = make([]*net.Conn, 0)

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
			c = append(c, &conn)
		}
	}()
}

func SendToNetwork(msg ControlMsg) {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	for _, cn := range c {
		if cn != nil {
			_, err = (*cn).Write(b)
			if err != nil {
				panic(err)
			}
		}
	}
}
