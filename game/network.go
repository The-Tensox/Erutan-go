package game

import (
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"math"

	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/erutan/utils"
	"github.com/golang/protobuf/ptypes"
)

type networkObject struct {
	Id         uint64
	components []*erutan.Component
}

type NetworkSystem struct {
	objects        octree.Octree
	lastUpdateTime float64
}

func NewNetworkSystem(lastUpdateTime float64) *NetworkSystem {
	return &NetworkSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(),
		utils.Config.GroundSize*10)),
		lastUpdateTime: lastUpdateTime}
}

func (n *NetworkSystem) Add(id uint64, components []*erutan.Component) {
	no := networkObject{Id: id, components: components}
	var position *protometry.VectorN
	for _, c := range components {
		if s := c.GetSpace(); s != nil {
			position = s.Position
		}
	}
	if position != nil {
		if ok := n.objects.Insert(*octree.NewObjectCube(no, position.Get(0),
			position.Get(1),
			position.Get(2),
			0.5)); !ok {
		 	utils.DebugLogf("Failed to insert %v", no.components)
		 }
	} else {
		n.objects.Insert(*octree.NewObjectCube(no, 0, 0, 0, 1))
	}
	// Broadcast on network the add
	ManagerInstance.Broadcast <- erutan.Packet{
		Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
		Type: &erutan.Packet_CreateEntity{
			CreateEntity: &erutan.Packet_CreateEntityPacket{
				EntityId:   id,
				Components: components,
			},
		},
	}
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (n *NetworkSystem) Remove(object octree.Object) {
	n.objects.Remove(object)
}

func (n *NetworkSystem) Update(dt float64) {
	// TODO: should probably be better to update only when there is a change ... (observer ..)
	// FOR NOW BRUTE FORCE DUMB ALGO

	objects := n.objects.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVector3Zero(), utils.Config.GroundSize*10))
	//utils.DebugLogf("len %v", len(objects))

	if (utils.GetProtoTime()-n.lastUpdateTime)/math.Pow(10, 9) > 0.0001/**float64(len(objects))*/ {
		for _, entity := range objects {

			if no, ok := entity.Data.(networkObject); ok {
				// Broadcast on network the update
				ManagerInstance.Broadcast <- erutan.Packet{
					Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
					Type: &erutan.Packet_UpdateEntity{
						UpdateEntity: &erutan.Packet_UpdateEntityPacket{
							EntityId:   no.Id,
							Components: no.components,
						},
					},
				}
			}
		}
		n.lastUpdateTime = utils.GetProtoTime()
	}
}

func (n *NetworkSystem) SyncNewClient(tkn string) {
	objects := n.objects.GetColliding(*protometry.NewBoxOfSize(*protometry.NewVector3Zero(), utils.Config.GroundSize*10))
	for _, entity := range objects {
		if no, ok := entity.Data.(networkObject); ok {
			c, _ := ManagerInstance.ClientStreams.Load(tkn)
			if res, ok := c.(chan erutan.Packet); ok {
				res <- erutan.Packet{
					Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
					Type: &erutan.Packet_CreateEntity{
						CreateEntity: &erutan.Packet_CreateEntityPacket{
							EntityId:   no.Id,
							Components: no.components,
						},
					},
				}
			}
		}
	}
}
