package game

import (
	"math"

	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"github.com/golang/protobuf/ptypes"

	"github.com/The-Tensox/erutan/ecs"
)

type networkEntity struct {
	*ecs.BasicEntity
	components []*erutan.Component
}

type NetworkSystem struct {
	entities       []networkEntity
	lastUpdateTime float64
}

func (n *NetworkSystem) Add(basic *ecs.BasicEntity, components []*erutan.Component) {
	n.entities = append(n.entities, networkEntity{BasicEntity: basic, components: components})

	// Broadcast on network the add
	ManagerInstance.Broadcast <- erutan.Packet{
		Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
		Type: &erutan.Packet_CreateEntity{
			CreateEntity: &erutan.Packet_CreateEntityPacket{
				EntityId:   n.entities[len(n.entities)-1].ID(),
				Components: n.entities[len(n.entities)-1].components,
			},
		},
	}

}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (n *NetworkSystem) Remove(basic ecs.BasicEntity) {
	var delete int = -1
	for index, entity := range n.entities {
		if entity.ID() == basic.ID() {
			delete = index
			break
		}
	}
	if delete >= 0 {
		n.entities = append(n.entities[:delete], n.entities[delete+1:]...)
	}
}

func (n *NetworkSystem) Update(dt float64) {
	if (utils.GetProtoTime()-n.lastUpdateTime)/math.Pow(10, 9) > 0.0002*float64(len(n.entities)) {
		for _, entity := range n.entities {
			// Broadcast on network the update
			ManagerInstance.Broadcast <- erutan.Packet{
				Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
				Type: &erutan.Packet_UpdateEntity{
					UpdateEntity: &erutan.Packet_UpdateEntityPacket{
						EntityId:   entity.ID(),
						Components: entity.components,
					},
				},
			}
		}
		n.lastUpdateTime = utils.GetProtoTime()
	}
}

func (n *NetworkSystem) SyncNewClient(tkn string) {
	for _, entity := range n.entities {
		c, _ := ManagerInstance.ClientStreams.Load(tkn)
		if res, ok := c.(chan erutan.Packet); ok {
			res <- erutan.Packet{
				Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
				Type: &erutan.Packet_CreateEntity{
					CreateEntity: &erutan.Packet_CreateEntityPacket{
						EntityId:   entity.ID(),
						Components: entity.components,
					},
				},
			}
		}
	}
}
