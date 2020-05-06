package erutan

type ClientPacket struct {
	ClientToken string
	Packet Packet
}

func NewClientPacket(token string, packet Packet) *ClientPacket {
	return &ClientPacket{token,packet}
}