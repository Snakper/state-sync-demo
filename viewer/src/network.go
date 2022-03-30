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
			m := &ControlMsg{}
			err = json.Unmarshal(bb, m)
			if err != nil {
				panic(err)
			}
			lock.Lock()
			msg = append(msg, m)
			lock.Unlock()
		}
	}()
}
