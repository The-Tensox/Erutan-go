package game

import (
	"github.com/The-Tensox/erutan/internal/cfg"
	"github.com/The-Tensox/erutan/internal/mon"
	"github.com/The-Tensox/erutan/internal/obs"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"math"

	"github.com/The-Tensox/erutan/internal/utils"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/golang/protobuf/ptypes"
)


type networkObject struct {
	Id uint64
	// clientsAction will apply a specific action for every action related to this object
	// for example filtering out an object if too far away ...
	clientsAction map[string]networkAction
	components    []*erutan.Component
}

type NetworkSystem struct {
	objects        octree.Octree
	lastUpdateTime float64
}

type networkAction int

const (
	ignore networkAction = iota
	update
	destroy
)

func NewNetworkSystem(lastUpdateTime float64) *NetworkSystem {
	return &NetworkSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(*protometry.NewVector3Zero(),
		cfg.Global.Logic.GroundSize*1000)),
		lastUpdateTime: lastUpdateTime}
}

func (n *NetworkSystem) Priority() int {
	return 1
}

func (n *NetworkSystem) Add(id uint64,
	size float64,
	components []*erutan.Component) {
	no := networkObject{Id: id, clientsAction: make(map[string]networkAction), components: components}
	ManagerInstance.ClientsSettings.Range(func(key, value interface{}) bool {
		no.clientsAction[key.(string)] = update
		return true
	})
	var position *protometry.VectorN
	for _, c := range components {
		if s := c.GetSpace(); s != nil {
			position = s.Position
		}
	}
	if position != nil {
		if ok := n.objects.Insert(*octree.NewObjectCube(no,
			position.Get(0),
			position.Get(1),
			position.Get(2),
			size)); !ok {
			utils.DebugLogf("Failed to insert, tree size: %v, object: %v", n.objects.GetSize(), no.components)
			return
		}
	} else { // Collision-less objects ?
		n.objects.Insert(*octree.NewObjectCube(no, 0, 0, 0, 1))
	}
	// Broadcast on network the add
	ManagerInstance.Broadcast <- erutan.Packet{
		Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
		Type: &erutan.Packet_UpdateEntity{
			UpdateEntity: &erutan.Packet_UpdateEntityPacket{
				EntityId:   id,
				Components: components,
			},
		},
	}
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (n *NetworkSystem) Remove(object octree.Object) {
	// Notify every clients of the removal of this object
	if no, ok := object.Data.(networkObject); ok {
		ManagerInstance.Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_DestroyEntity{
				DestroyEntity: &erutan.Packet_DestroyEntityPacket{
					EntityId: no.Id,
				},
			},
		}
	}

	n.objects.Remove(object)
}

// NetworkSystem Update function role is to synchronise what we want to be synchronised with clients
// Thus only entities that have been added to this system will be synchronised
// Plus it should take into account some client-specific preferences and filter in/out some entities
// For example client X might want to be synchronised only on the area in: sphere centered 12;0;-100 of radius 100
// Or client Y might want to see the Octree data structure ...
func (n *NetworkSystem) Update(dt float64) {
	// TODO: should it be better to update only when there is a change ... (observer ..) ?
	// Limit synchronisation to specific fps to avoid burning your computer
	if (utils.GetProtoTime()-n.lastUpdateTime)/math.Pow(10, 9) > 0.1 /**float64(len(objects))*/ { // times len obj = more obj = less net updates
		objects := n.objects.GetObjects()
		for _, entity := range objects {
			if no, ok := entity.Data.(networkObject); ok {
				for keyClient, clientValue := range no.clientsAction {
					if streamInterface, hasStream := ManagerInstance.ClientsOut.Load(keyClient); hasStream {
						// Get channel
						if channel, isChannel := streamInterface.(chan erutan.Packet); isChannel {
							switch clientValue {
							case ignore: // Ignore is continuous
								mon.NetworkActionIgnoreCounter.Inc()

							case update: // Update is continuous too, don't change it
								mon.NetworkActionUpdateCounter.Inc()

								channel <- erutan.Packet{
									Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
									Type: &erutan.Packet_UpdateEntity{
										UpdateEntity: &erutan.Packet_UpdateEntityPacket{
											EntityId:   no.Id,
											Components: no.components,
										},
									},
								}
							case destroy: // Destroy is discrete, only do it once then ignore
								mon.NetworkActionDestroyCounter.Inc()

								channel <- erutan.Packet{
									Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
									Type: &erutan.Packet_DestroyEntity{
										DestroyEntity: &erutan.Packet_DestroyEntityPacket{
											EntityId: no.Id,
										},
									},
								}
								no.clientsAction[keyClient] = ignore
							}
						}
					}
				}
			}
		}
		n.lastUpdateTime = utils.GetProtoTime()
	}
}

func (n *NetworkSystem) Handle(event obs.Event) {
	switch settings := event.Value.(type) {
	case obs.ClientConnected:
		objects := n.objects.GetObjects()
		debugAction := isDebug(settings.Settings.UpdateParameters.Parameters)

		// Now tag every object with a network action
		for i := range objects {
			// Cast to network object
			if no, isNo := objects[i].Data.(networkObject); isNo {
				var netBehaviour *erutan.Component_NetworkBehaviourComponent
				for _, c := range no.components {
					netBehaviour = c.GetNetworkBehaviour()
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_DEBUG {
					no.clientsAction[settings.ClientToken] = debugAction
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_ALL {
					// Object isn't tagged with debug network behaviour, just update
					no.clientsAction[settings.ClientToken] = update
				}
				//utils.DebugLogf("%v", no.clientsAction[settings.ClientToken])
			}
		}

	// Depending on settings, network system will "tag" every objects with an action to do for each clients
	case obs.ClientSettingsUpdated:
		objects := n.objects.GetObjects()
		debugAction := isDebug(settings.Settings.UpdateParameters.Parameters)

		// Now tag every object with a network action
		for i := range objects {
			// Cast to network object
			if no, isNo := objects[i].Data.(networkObject); isNo {
				var netBehaviour *erutan.Component_NetworkBehaviourComponent
				for _, c := range no.components {
					netBehaviour = c.GetNetworkBehaviour()
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_DEBUG {
					no.clientsAction[settings.ClientToken] = debugAction
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_ALL {
					// Object isn't tagged with debug network behaviour, just update
					no.clientsAction[settings.ClientToken] = update
				}
				//utils.DebugLogf("%v", no.clientsAction[settings.ClientToken])
			}
		}
	}
}

func isDebug(params []*erutan.Packet_UpdateParametersPacket_Parameter) networkAction {
	var debugAction networkAction
	// Get the network action
	for _, paramInterface := range params {
		switch p := paramInterface.Type.(type) {
		case *erutan.Packet_UpdateParametersPacket_Parameter_Debug:
			// So the client asked for debug mode just now
			if p.Debug { // Debug objects need to be created
				debugAction = update
			} else { // Otherwise he turned it off, debug objects need to be destroyed
				debugAction = destroy
			}
		}
	}
	return debugAction
}
