package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/as/mqtt"
	"github.com/as/wire9/varint"
)

func str(s string) mqtt.String {
	return mqtt.String{uint16(len(s)), []byte(s)}
}

func main() {
	if len(os.Args) <= 4{
		log.Fatalln("usage: ex server:port clientid topic msg")
	}
	dst := os.Args[1]
	name := os.Args[2]
	topic := os.Args[3]
	msg := os.Args[4]

	buf := new(bytes.Buffer)
	cp := mqtt.PackConn{}
	cp.Head.Name = str("MQTT")
	cp.Head.Level = 4
	cp.Head.Flags |= mqtt.ConnFlag
	cp.Payload.Client = str(name)
	cp.Payload.Topic = str(topic)
	cp.Payload.Message = str(msg)

	conn, err := net.Dial("tcp", dst)
	if err != nil {
		log.Fatalln(err)
	}
	cp.WriteBinary(buf)
	hdr := mqtt.Head{
		Flags: byte(mqtt.ContConn << 4),
		Len:   varint.V(buf.Len()),
	}
	buf.Reset()
	hdr.WriteBinary(buf)
	cp.WriteBinary(buf)	// lazy double write here
	io.Copy(conn, buf)

	hdr = mqtt.Head{}
	hdr.ReadBinary(conn)
	ack := &mqtt.PackConnAck{}
	ack.ReadBinary(conn)
	fmt.Printf("head: %#v\nack: %#v\n", hdr, ack)
}
