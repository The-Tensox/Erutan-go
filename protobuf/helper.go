package erutan

/*
import (
	"github.com/The-Tensox/erutan/utils"
)
// TODO!!!!!!!!!!!!!!!
func (sc *Component_SpaceComponent) Update(newSc Component_SpaceComponent) bool {
	ManagerInstance.Watch.Notify(utils.Event{Value: EntityPhysicsUpdated{id: entity.ID(), newSc: newSc, dt: dt}})
	return true
}
*/

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