package erutan

// ClientPacket (is there a better name) holds a tuple client identification token - packet
type ClientPacket struct {
	Token string
	Packet Packet
}

func NewClientPacket(token string, packet Packet) *ClientPacket {
	return &ClientPacket{
		Token: token,
		Packet: packet,
	}
}