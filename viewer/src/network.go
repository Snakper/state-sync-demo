package src

import (
	"bytes"
	"encoding/json"
	"net"
)

func ConnectToServer() {
	conn, err := net.Dial("tcp", "127.0.0.1:9123")
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			var b = make([]byte, 1024)
			_, err := conn.Read(b)
			if err != nil {
				panic(err)
			}
			index := bytes.IndexByte(b, 0)
			bb := b[:index]
			lock.Lock()
			if len(msg) == 0 {
				err = json.Unmarshal(bb, &msg)
				if err != nil {
					panic(err)
				}
			} else {
				m := make([]ControlMsg, 0)
				err = json.Unmarshal(bb, &msg)
				if err != nil {
					panic(err)
				}
				msg = append(msg, m...)
			}
			lock.Unlock()
		}
	}()
}
