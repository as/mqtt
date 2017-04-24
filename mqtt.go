package mqtt

import(
	 "io"
)
// TODO: default endianness

type utf8 string // TODO: Conformance

//wire9 String N[2] Data[N]
//wire9 Head Flags[1] Len[varint.V]

//wire9 PackConnPayload Client[,String] Topic[,String] Message[,String] User[,String] Pass[,String]
//wire9 PackConnHead Name[,String] Level[1] Flags[1] KeepAlive[2]
//wire9 PackConn Head[,PackConnHead] Payload[,PackConnPayload]

//wire9 PackConnHead Name[,String] Level[1] Flags[1] KeepAlive[2]

//!wire9 PackPub Flags[1] Len[varint.V] N[2] Topic[N] PID[2] Data[Len-2-N]

const (
	ConnNone byte = 1 << iota
	ConnClean
	ConnFlag
	ConnQoS1
	ConnQoS2
	ConnRetain
	ConnUser
)

//wire9 PackConnAck Flags[1] Return[1]

type AckRet byte

const (
	AckOk AckRet = iota
	AckBadID
	AckBadServer
	AckBadCreds
	AckBadAuth
)

func (z *PackConn) ReadBinary(r io.Reader) error {
	err := z.Head.ReadBinary(r)
	if err != nil {
		return err
	}
	z.Payload = PackConnPayload{
		Client: String{},
	}
	(&(z.Payload.Client)).ReadBinary(r)
	
	f :=  z.Head.Flags
	if f >>= 1; f&1 == 1 { // Will Topic
		// Topic and Message
		var t, m String
		if err = (&t).ReadBinary(r); err != nil {
			return err
		}
		if err = (&m).ReadBinary(r); err != nil {
			return err
		}
		z.Payload.Message = m
		z.Payload.Topic = t
	}
	if f >>= 1; f&1 == 1 { //  User Name and Password
		var u, p String
		if err = (&u).ReadBinary(r); err != nil {
			return err
		}
		if err = (&p).ReadBinary(r); err != nil {
			return err
		}
		z.Payload.User = u
		z.Payload.Pass = p
	}
	return nil
}
