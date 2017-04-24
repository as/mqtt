package mqtt

type Cont byte

const (
	ContForbidden0 Cont = iota
	ContConn
	ContConnACK
	ContPub
	ContPubACK
	ContPubREC
	ContPubREL
	ContPubCOMP
	ContSub
	ContSubACK
	ContUnsub
	ContUnsubACK
	ContPingTX
	ContPingRX
	ContDisconnect
	ContForbidden1
)

func (c Cont) Valid() bool {
	return c > ContForbidden0 && c < ContForbidden1
}
func (c Cont) HasVarHead() bool {
	if !c.Valid() {
		return false
	}
	switch c {
	case ContConn, ContConnACK, ContPingTX, ContPingRX:
		return false
	case ContPub:
		panic("TODO: HasVarHead: Need to check if QoS is set")
	}
	return true
}
