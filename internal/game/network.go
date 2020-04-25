package game

import (
	"github.com/The-Tensox/erutan/internal/cfg"
	"github.com/The-Tensox/erutan/internal/mon"
	"github.com/The-Tensox/erutan/internal/obs"
	"github.com/The-Tensox/erutan/internal/utils"
	erutan "github.com/The-Tensox/erutan/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"github.com/golang/protobuf/ptypes"
	"math"
)

type networkObject struct {
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
	return &NetworkSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0,
		cfg.Global.Logic.GroundSize*1000)),
		lastUpdateTime: lastUpdateTime}
}

func (n *NetworkSystem) Priority() int {
	return 1
}

func (n *NetworkSystem) Add(object octree.Object,
	components []*erutan.Component) {
	// Create the network object with its "mask" of network actions & components
	no := &networkObject{clientsAction: make(map[string]networkAction), components: components}
	var position *protometry.Vector3
	debug := false
	for _, c := range components {
		switch t := c.Type.(type) {
		case *erutan.Component_Space:
			position = t.Space.Position
		case *erutan.Component_NetworkBehaviour:
			if t.NetworkBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_DEBUG {
				debug = true
			}
		}
	}

	// We want to tag this new object with appropriate action to take regarding network
	// FIXME: clients still see new runtime added debug objects ?
	ManagerInstance.ClientsSettings.Range(func(key, value interface{}) bool {
		if debug {
			no.clientsAction[key.(string)] = ignore
		} else {
			no.clientsAction[key.(string)] = update
		}
		return true
	})

	if position == nil {
		utils.DebugLogf("Failed to insert, no space component provided: %v", no.components)
		return
	}
	object.Data = no
	if ok := n.objects.Insert(object); !ok {
		utils.DebugLogf("Failed to insert, tree size: %v, object: %v", n.objects.GetSize(), no.components)
		return
	}

	// Broadcast on network the add
	ManagerInstance.Broadcast <- erutan.Packet{
		Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
		Type: &erutan.Packet_UpdateEntity{
			UpdateEntity: &erutan.Packet_UpdateEntityPacket{
				EntityId:   object.ID(),
				Components: components,
			},
		},
	}
}

// Remove removes the Entity from the System. This is what most Remove methods will look like
func (n *NetworkSystem) Remove(object octree.Object) {
	if n.objects.Remove(object) {
		// Notify every clients of the removal of this object
		ManagerInstance.Broadcast <- erutan.Packet{
			Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
			Type: &erutan.Packet_DestroyEntity{
				DestroyEntity: &erutan.Packet_DestroyEntityPacket{
					EntityId: object.ID(),
				},
			},
		}
	} else {
		utils.DebugLogf( "Failed to remove")
	}
}

// NetworkSystem Update function role is to synchronise what we want to be synchronised with clients
// Thus only entities that have been added to this system will be synchronised
// Plus it should take into account some client-specific preferences and filter in/out some entities
// For example client X might want to be synchronised only on the area in: sphere centered 12;0;-100 of radius 100
// Or client Y might want to see the Octree data structure ...
func (n *NetworkSystem) Update(dt float64) {
	// TODO: should it be better to update only when there is a change ... (observer ..) ?
	// Limit synchronisation to specific fps to avoid burning your computer
	if (utils.GetProtoTime()-n.lastUpdateTime)/math.Pow(10, 9) > cfg.Global.NetworkRate /**float64(len(objects))*/ { // times len obj = more obj = less net updates
		n.objects.Range(func(object *octree.Object) bool {
			if no, ok := object.Data.(*networkObject); ok {
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
											EntityId:   object.ID(),
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
											EntityId: object.ID(),
										},
									},
								}
								no.clientsAction[keyClient] = ignore
							}
						}
					}
				}
			}
			return true
		})
		n.lastUpdateTime = utils.GetProtoTime()
	}
}

func (n *NetworkSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	case obs.OnPhysicsUpdateResponse:
		// No collision here
		if e.Other == nil {
			me := Find(n.objects, *e.Me)
			if me == nil {
				utils.DebugLogf("Unable to find %v in system %T", e.Me.ID(), n)
				return
			}
			asNo := me.Data.(*networkObject)
			for i := range asNo.components {
				switch c := asNo.components[i].Type.(type) {
				case *erutan.Component_Space:
					*c.Space.Position = e.NewPosition
				}
			}
			// Need to reinsert in the octree
			if !n.objects.Move(me, e.NewPosition.X, e.NewPosition.Y, e.NewPosition.Z) {
				utils.DebugLogf("Failed to move %v", me)
			}
			//utils.DebugLogf("net move %v %v", center, asNo.components)

			// Over
			return
		}
	case obs.OnClientConnection:
		debugAction := isDebug(e.Settings.UpdateParameters.Parameters)
		// Now tag every object with a network action
		n.objects.Range(func(object *octree.Object) bool {
			// Cast to network object
			if no, isNo := object.Data.(*networkObject); isNo {
				var netBehaviour *erutan.Component_NetworkBehaviourComponent
				for _, c := range no.components {
					netBehaviour = c.GetNetworkBehaviour()
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_DEBUG {
					no.clientsAction[e.ClientToken] = debugAction
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_ALL {
					// Object isn't tagged with debug network behaviour, just update
					no.clientsAction[e.ClientToken] = update
				}
			}
			return true
		})

	// Depending on settings, network system will "tag" every objects with an action to do for each clients
	case obs.OnClientSettingsUpdate:
		debugAction := isDebug(e.Settings.UpdateParameters.Parameters)
		// Now tag every object with a network action
		n.objects.Range(func(object *octree.Object) bool {
			// Cast to network object
			if no, isNo := object.Data.(*networkObject); isNo {
				var netBehaviour *erutan.Component_NetworkBehaviourComponent
				for _, c := range no.components {
					netBehaviour = c.GetNetworkBehaviour()
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_DEBUG {
					no.clientsAction[e.ClientToken] = debugAction
				}
				if netBehaviour != nil && netBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_ALL {
					// Object isn't tagged with debug network behaviour, just update
					no.clientsAction[e.ClientToken] = update
				}
			}
			return true
		})
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
