package game

import (
	"github.com/golang/protobuf/ptypes"
	erutan "github.com/user/erutan/protos/realtime"
	"github.com/user/erutan/utils"

	"github.com/user/erutan/ecs"
)

type networkEntity struct {
	*ecs.BasicEntity
	*erutan.Component_SpaceComponent
	*erutan.Component_SpaceTimeComponent
}

type NetworkSystem struct {
	entities []networkEntity
}

func (n *NetworkSystem) Add(basic *ecs.BasicEntity,
	space *erutan.Component_SpaceComponent) {
	n.entities = append(n.entities, networkEntity{basic, space,
		&erutan.Component_SpaceTimeComponent{
			Timestamp: ptypes.TimestampNow(),
			Space: &erutan.Component_SpaceComponent{ // Last position initialized as first position
				Position: space.Position,
				Rotation: space.Rotation,
				Scale:    space.Scale,
			},
		}})

	// Broadcast on network the add
	ManagerInstance.Broadcast <- erutan.Packet{
		Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
		Type: &erutan.Packet_CreateEntity{
			CreateEntity: &erutan.Packet_CreateEntityPacket{
				EntityId: n.entities[len(n.entities)-1].ID(),
				Components: []*erutan.Component{
					&erutan.Component{
						Type: &erutan.Component_Space{
							Space: n.entities[len(n.entities)-1].Component_SpaceComponent,
						},
					},
				},
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
	for _, entity := range n.entities {
		// If moved, rotated, rescaled, sync network
		if utils.Distance(*entity.Position, *entity.Component_SpaceTimeComponent.Space.Position) > 1 ||
			entity.Rotation != entity.Component_SpaceTimeComponent.Space.Rotation ||
			entity.Scale != entity.Component_SpaceTimeComponent.Space.Scale {

			//utils.DebugLogf("Network space update at %v", ptypes.TimestampNow())
			// Broadcast on network the update
			ManagerInstance.Broadcast <- erutan.Packet{
				Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
				Type: &erutan.Packet_UpdatePosition{
					UpdatePosition: &erutan.Packet_UpdatePositionPacket{
						EntityId: entity.ID(),
						Position: entity.Position, // Refer to the Space component position
					},
				},
			}
			// Refresh last space
			entity.Component_SpaceTimeComponent = &erutan.Component_SpaceTimeComponent{Space: entity.Space,
				Timestamp: ptypes.TimestampNow()}
		}
	}
}

func (n *NetworkSystem) SyncNewClient(tkn string) {
	ManagerInstance.StreamsMtx.Lock()
	for _, entity := range n.entities {
		ManagerInstance.ClientStreams[tkn] <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_CreateEntity{
				CreateEntity: &erutan.Packet_CreateEntityPacket{
					EntityId: entity.ID(),
					Components: []*erutan.Component{
						&erutan.Component{
							Type: &erutan.Component_Space{
								Space: entity.Component_SpaceComponent,
							},
						},
					},
				},
			},
		}
	}
	ManagerInstance.StreamsMtx.Unlock()
}
