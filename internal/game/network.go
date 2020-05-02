package game

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/mon"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"github.com/golang/protobuf/ptypes"
	"math"
	"time"
)

type networkObject struct {
	// clientsAction will apply a specific action for every action related to this object
	// for example filtering out an object if too far away ...
	clientsAction map[string]networkSettings
	components    []*erutan.Component
}

type NetworkSystem struct {
	objects        octree.Octree
	lastUpdateTime float64
}

type networkAction int

const (
	ignore networkAction = iota // Don't show to client
	update // Show to client
	hide // Tell client to hide (destroy client's locally)
)

type networkSettings struct {
	networkAction         // Action to take regarding to a client (ignore, update, destroy ...
	lastUpdate    float64 // Last time an object has been synced with client, unused yet
	// Can be useful to prevent cheating for objects owned by the client
}

func NewNetworkSystem(lastUpdateTime float64) *NetworkSystem {
	return &NetworkSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Global.Logic.GroundSize*1000)),
		lastUpdateTime: lastUpdateTime}
}

func (n *NetworkSystem) Priority() int {
	return 1
}

func (n *NetworkSystem) Add(object octree.Object,
	components []*erutan.Component) {
	// Create the network object with its "mask" of network actions & components
	no := &networkObject{clientsAction: make(map[string]networkSettings), components: components}
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
			no.clientsAction[key.(string)] = networkSettings{networkAction: ignore, lastUpdate: utils.GetProtoTime()}
		} else {
			no.clientsAction[key.(string)] = networkSettings{networkAction: update, lastUpdate: utils.GetProtoTime()}
		}
		return true
	})

	// All objects even non-physics should have a position even 0;0;0, since our data structures are spatially ordered
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
	//ManagerInstance.Broadcast <- erutan.Packet{
	//	Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
	//	Type: &erutan.Packet_UpdateObject{
	//		UpdateObject: &erutan.Packet_UpdateObjectPacket{
	//			ObjectId:   object.ID(),
	//			Components: components,
	//		},
	//	},
	//}
	mon.NetworkAddCounter.Inc()
	//x :=n.objects.GetAllObjects()
	//for _, i := range x {
	//	utils.DebugLogf("tree add %v", i.ID())
	//}
}

// Remove removes the Object from the System. This is what most Remove methods will look like
func (n *NetworkSystem) Remove(object octree.Object) {
	// Notify every clients of the removal of this object
	ManagerInstance.Broadcast <- erutan.Packet{
		Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
		Type: &erutan.Packet_DestroyObject{
			DestroyObject: &erutan.Packet_DestroyObjectPacket{
				ObjectId: object.ID(),
			},
		},
	}
	if !n.objects.Remove(object) {

	} else {
		utils.DebugLogf("Failed to remove %d, data: %T", object.ID(), object.Data)
		//x :=n.objects.GetAllObjects()
		//for _, i := range x {
		//	utils.DebugLogf("tree rm %v", i.ID())
		//}
	}
	mon.NetworkRemoveCounter.Inc()
}

// NetworkSystem Update function role is to synchronise what we want to be synchronised with clients
// Thus only objects that have been added to this system will be synchronised
// Plus it should take into account some client-specific preferences and filter in/out some objects
// For example client X might want to be synchronised only on the area in: sphere centered 12;0;-100 of radius 100
// Or client Y might want to see the Octree data structure ...
func (n *NetworkSystem) Update(dt float64) {
	// TODO: should it be better to update only when there is a change ... (observer ..) ?
	// Limit synchronisation to specific fps to avoid burning your computer
	if (utils.GetProtoTime()-n.lastUpdateTime)/math.Pow(10, 9) > cfg.Global.NetworkRate /**float64(len(objects))*/ { // times len obj = more obj = less net updates
		for _, object := range n.objects.GetAllObjects() {
		//n.objects.Range(func(object *octree.Object) bool {
			if no, ok := object.Data.(*networkObject); ok {
				for keyClient, clientValue := range no.clientsAction { // TODO: well done its probably possible to handle each client concurrently
					if channel, hasChannel := ManagerInstance.ClientsOut[keyClient]; hasChannel {
						// Get channel
						switch clientValue.networkAction {
						case ignore: // Ignore is continuous
							//utils.DebugLogf("ignore")
						case update: // Update is continuous too, don't change it
							var isOwner bool
							for _, cType := range no.components {
								switch c := cType.Type.(type) {
								case *erutan.Component_NetworkBehaviour:
									// Let's check if this client is owner of this object
									if keyClient == c.NetworkBehaviour.OwnerToken {
										isOwner = true
									}
									break
								}
							}
							now := utils.GetProtoTime()
							// If he is owner of the object, we don't send him updates (easy cheating)

							dif := (now - clientValue.lastUpdate) / math.Pow(10, 9)
							//utils.DebugLogf("clientValue.lastUpdate %v, now - lastUpdate %v", clientValue.lastUpdate,
							//	dif)
							if isOwner {
							} else if !isOwner && dif > 5 { // Quick hack to reduce throughput
								channel <- erutan.Packet{
									Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
									Type: &erutan.Packet_UpdateObject{
										UpdateObject: &erutan.Packet_UpdateObjectPacket{
											ObjectId:   object.ID(),
											Components: no.components,
										},
									},
								}
								no.clientsAction[keyClient] = networkSettings{networkAction: update, lastUpdate: now}
							}

						case hide: // Destroy is discrete, only do it once then ignore
							channel <- erutan.Packet{
								Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
								Type: &erutan.Packet_DestroyObject{
									DestroyObject: &erutan.Packet_DestroyObjectPacket{
										ObjectId: object.ID(),
									},
								},
							}
							no.clientsAction[keyClient] = networkSettings{networkAction: ignore,
								lastUpdate: no.clientsAction[keyClient].lastUpdate}
						}
					}

				}
			}
			//return true
		//})
		}
		n.lastUpdateTime = utils.GetProtoTime()
	}
}

func (n *NetworkSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	case obs.PhysicsUpdateResponse:
		// No collision here
		if e.Other == nil {
			var me *octree.Object
			var asNo *networkObject
			// If the event was coming with an object from another system
			if no, ok := e.Me.Data.(*networkObject); !ok { // FIXME
				me = Find(n.objects, *e.Me)
				if me == nil {
					utils.DebugLogf("Unable to find %v in system %T", e.Me.ID(), n)
					return
				}
				asNo = me.Data.(*networkObject)
			} else { // Otherwise just use the event object
				me = e.Me
				asNo = no
			}

			for i := range asNo.components {
				switch c := asNo.components[i].Type.(type) {
				case *erutan.Component_Space:
					*c.Space.Position = e.NewPosition
					break
				}
			}
			// Need to reinsert in the octree
			if !n.objects.Move(me, e.NewPosition.X, e.NewPosition.Y, e.NewPosition.Z) {
				//utils.DebugLogf("objects %v", n.objects.GetAllObjects())
				utils.DebugLogf("Failed to move %v; %v to %v", me.ID(), me.Bounds.GetCenter(), e.NewPosition)
			}
			//utils.DebugLogf("net move %v %v", center, asNo.components)

			// Over
			return
		}
	case obs.ClientConnection:
		n.updateClientAction(e.ClientToken, isClientDebugging(e.Settings.UpdateParameters.Parameters),
			getCullingArea(e.Settings.UpdateParameters.Parameters))
		utils.DebugLogf("ClientConnection %v", e.ClientToken)
		//p := protometry.RandomSpherePoint(*protometry.NewVector3(0, 100, 0), 20)
		utils.DebugLogf("Spawning player for client [%s] at %v", e.ClientToken, cfg.Global.Logic.Player.Spawn)
		id, data := ManagerInstance.AddPlayer(&cfg.Global.Logic.Player.Spawn, e.ClientToken)
		// Notify everyone of the creation of this player object
		// Somehow need to wait a little bit before sending
		go time.AfterFunc(1*time.Second, func() {
			ManagerInstance.Broadcast <- erutan.Packet{
				Metadata: &erutan.Metadata{Timestamp: ptypes.TimestampNow()},
				Type: &erutan.Packet_CreatePlayer{
					CreatePlayer: &erutan.Packet_CreatePlayerPacket{
						ObjectId: id,
						Components: []*erutan.Component{
							{Type: &erutan.Component_Space{Space: data.Component_SpaceComponent}},
							{Type: &erutan.Component_Render{Render: data.Component_RenderComponent}},
							{Type: &erutan.Component_NetworkBehaviour{NetworkBehaviour: data.Component_NetworkBehaviourComponent}},
						},
					},
				},
			}
		})
	case obs.ClientDisconnection:
		// Currently default behaviour will remove all objects owned by this client
		for _, object := range n.objects.GetAllObjects() {
			// Cast to network object
			if no, isNo := object.Data.(*networkObject); isNo {
				for _, c := range no.components {
					switch c := c.Type.(type) {
					case *erutan.Component_NetworkBehaviour:
						if c.NetworkBehaviour.OwnerToken == e.ClientToken {
							ManagerInstance.World.RemoveObject(object)
						}
						break
					}
				}
			}
		}


	// Depending on settings, network system will "tag" every objects with an action to do for each clients
	case obs.ClientSettingsUpdate:
		n.updateClientAction(e.ClientToken, isClientDebugging(e.Settings.UpdateParameters.Parameters),
			getCullingArea(e.Settings.UpdateParameters.Parameters))
	}
}

func (n *NetworkSystem) updateClientAction(clientToken string, isClientDebugging bool, cullingArea *protometry.Box) {
	// Now tag every object with a network action
	//n.objects.Range(func(object *octree.Object) bool {
	for _, object := range n.objects.GetAllObjects() {
		// Cast to network object
		if no, isNo := object.Data.(*networkObject); isNo {
			for _, c := range no.components {
				switch c := c.Type.(type) {
				case *erutan.Component_NetworkBehaviour:
					// Owned objects default behaviour is to be shown anyway
					if c.NetworkBehaviour.OwnerToken == clientToken {
						no.clientsAction[clientToken] = networkSettings{networkAction: update,
							lastUpdate: no.clientsAction[clientToken].lastUpdate}
						return
					}

					// Otherwise filter out objects outside client's culling area
					if cullingArea != nil && object.Bounds.Fit(*cullingArea) {
						// Filter out / in debug objects
						if isClientDebugging && // First easy condition filter most of things
							c.NetworkBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_DEBUG { // Object is a debug thing
							if no.clientsAction[clientToken].networkAction != ignore { // Hide once !!
								no.clientsAction[clientToken] = networkSettings{networkAction: hide,
									lastUpdate: no.clientsAction[clientToken].lastUpdate}
							}
						} else if c.NetworkBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_ALL {
							// Object isn't tagged with debug network behaviour, just update
							no.clientsAction[clientToken] = networkSettings{networkAction: update,
								lastUpdate: no.clientsAction[clientToken].lastUpdate}
						}
						//utils.DebugLogf("fit")
					} else if no.clientsAction[clientToken].networkAction != ignore { // Destroy client's objects that have just been out of culling area
						no.clientsAction[clientToken] = networkSettings{networkAction: hide,
							lastUpdate: no.clientsAction[clientToken].lastUpdate}
					}
					break
				}
			}

		}
	}
		//return true
	//})

}

func isClientDebugging(params []*erutan.Packet_UpdateParametersPacket_Parameter) bool {
	// Get the network action
	for _, paramInterface := range params {
		switch p := paramInterface.Type.(type) {
		case *erutan.Packet_UpdateParametersPacket_Parameter_Debug:
			// So the client asked for debug mode just now
			return p.Debug
		}
	}
	return false
}

func getCullingArea(params []*erutan.Packet_UpdateParametersPacket_Parameter) *protometry.Box {
	// Get the culling area
	for _, paramInterface := range params {
		switch p := paramInterface.Type.(type) {
		case *erutan.Packet_UpdateParametersPacket_Parameter_CullingArea:
			return p.CullingArea
		}
	}
	return nil
}