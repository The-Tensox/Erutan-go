package game

import (
	"github.com/The-Tensox/Erutan-go/internal/cfg"
	"github.com/The-Tensox/Erutan-go/internal/mon"
	"github.com/The-Tensox/Erutan-go/internal/obs"
	"github.com/The-Tensox/Erutan-go/internal/utils"
	erutan "github.com/The-Tensox/Erutan-go/protobuf"
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
	"math"
)

type networkObject struct {
	// clientsAction will apply a specific action for every action related to this object
	// for example filtering out an object if too far away ...
	clientsAction map[string]networkAction
	lastUpdate    float64
	components    []*erutan.Component
}

type NetworkSystem struct {
	objects   octree.Octree
	keepAlive float64
}

type networkAction int

const (
	ignore networkAction = iota // Don't show to client
	update                      // Show to client
	hide                        // Tell client to hide (destroy client's locally)
)

func NewNetworkSystem(lastUpdateTime float64) *NetworkSystem {
	return &NetworkSystem{objects: *octree.NewOctree(protometry.NewBoxOfSize(0, 0, 0, cfg.Global.Logic.OctreeSize)),
		keepAlive: lastUpdateTime}
}

func (n NetworkSystem) Priority() int {
	return math.MaxInt64 - 1
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
		if updParam, ok := value.(erutan.Packet_UpdateParameters); ok {
			for _, p := range updParam.UpdateParameters.Parameters {
				switch t := p.Type.(type) {
				case *erutan.Packet_UpdateParametersPacket_Parameter_Debug:
					// If client is in debug mode and it's a debug object, display
					if debug && t.Debug {
						no.clientsAction[key.(string)] = update
					} else {
						no.clientsAction[key.(string)] = ignore
					}
				}
			}
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

	//utils.DebugLogf("Inserted %v at %v", object.ID(), object.Bounds.GetCenter())
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
	if !n.objects.Remove(object) {
		utils.DebugLogf("Failed to remove %d, data: %T", object.ID(), object.Data)
		//x :=n.objects.GetAllObjects()
		//for _, i := range x {
		//	utils.DebugLogf("tree rm %v", i.ID())
		//}
	} else {
		// Notify every clients of the removal of this object
		ManagerInstance.Broadcast(erutan.Packet{
			Type: &erutan.Packet_DestroyObject{
				DestroyObject: &erutan.Packet_DestroyObjectPacket{
					ObjectId: object.ID(),
				},
			},
		})
		mon.NetworkRemoveCounter.Inc()
	}
}

// NetworkSystem Update function role is to synchronise what we want to be synchronised with clients
// Thus only objects that have been added to this system will be synchronised
// Plus it should take into account some client-specific preferences and filter in/out some objects
// For example client X might want to be synchronised only on the area in: sphere centered 12;0;-100 of radius 100
// Or client Y might want to see the Octree data structure ...
func (n *NetworkSystem) Update(dt float64) {
	// Limit synchronisation to specific fps to avoid burning your computer
	if (utils.GetProtoTime()-n.keepAlive)/math.Pow(10, 9) > cfg.Global.NetworkRate /**float64(len(objects))*/ { // times len obj = more obj = less net updates
		n.syncWholeState()
		n.keepAlive = utils.GetProtoTime()
	}
}

func (n NetworkSystem) Handle(event obs.Event) {
	switch e := event.Value.(type) {
	case obs.PhysicsUpdateResponse:
		// Update position of every objects, if there was a collision or not
		for i := range e.Objects { // TODO: stuff when player collide refuse moving
			//utils.DebugLogf("need to move %v; %v to %v", obj.ID(), obj.Bounds.GetCenter(), obj.Vector3)

			// Super ugly: we need to check if the incoming object has Data as *networkObject or not
			var me *octree.Object
			var asNo *networkObject
			if _, ok := e.Objects[i].Data.(*networkObject); !ok {
				//me = Find(n.objects, obj.Object)
				me = n.objects.Get(e.Objects[i].Object.ID(), e.Objects[i].Object.Bounds)
			} else {
				me = &e.Objects[i].Object
			}
			if me == nil {
				utils.DebugLogf("Unable to find %v in system %T", e.Objects[i].Object.ID(), n)
				//for _, o := range n.objects.GetAllObjects() {
				//	utils.DebugLogf("my objs %v", o.ID())
				//}
				return
			}
			//distanceDone := obj.Vector3.Distance(me.Bounds.GetCenter())
			asNo, ok := me.Data.(*networkObject)
			if ok {
				for _, c := range asNo.components {
					switch t := c.Type.(type) {
					case *erutan.Component_Space:
						t.Space.Position = &e.Objects[i].Vector3
						break
					}
				}
			}

			// Need to reinsert in the octree
			if !n.objects.Move(me, e.Objects[i].Vector3.X, e.Objects[i].Vector3.Y, e.Objects[i].Vector3.Z) {
				utils.DebugLogf("Failed to move %v to %v", me.ID(), e.Objects[i].Vector3)
				//for _, o := range n.objects.GetAllObjects() {
				//	utils.DebugLogf("my objs %v", o.ID())
				//}
				continue
			} else {
				//utils.DebugLogf("moved %v to %v", me.ID(), me.Bounds.GetCenter())
			}
			//utils.DebugLogf("lastupdate %v", (utils.GetProtoTime()-asNo.lastUpdate)/math.Pow(10, 9))
			// Sync object that moved (obviously won't sync static objects)
			if (utils.GetProtoTime()-asNo.lastUpdate)/math.Pow(10, 9) > 0.01 {
				n.syncSingleObject(*me) // TODO: maybe some objects need lower sync than other i.e players vs mobs
			}
		}

	case obs.ClientConnection:
		n.updateClientAction(e.ClientToken, isClientDebugging(e.Settings.UpdateParameters.Parameters),
			getCullingArea(e.Settings.UpdateParameters.Parameters))
		utils.DebugLogf("ClientConnection %v", e.ClientToken)
		//p := protometry.RandomSpherePoint(*protometry.NewVector3(0, 100, 0), 20)
		p := protometry.NewVector3(cfg.Global.Logic.Player.Spawn[0], cfg.Global.Logic.Player.Spawn[1], cfg.Global.Logic.Player.Spawn[2])
		id, data := ManagerInstance.AddPlayer(p, e.ClientToken)
		utils.DebugLogf("Spawning player for client [%s] with id %d at %v", e.ClientToken, id, cfg.Global.Logic.Player.Spawn)

		// Notify everyone of the creation of this player object
		// Somehow need to wait a little bit before sending
		//n.syncWholeState()

		//time.Sleep(3*time.Second)
		utils.DebugLogf("Send client connection [%s]", e.ClientToken)
		//utils.DebugLogf("blocked")

		ManagerInstance.Send(e.ClientToken, erutan.Packet{
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
		})
		//utils.DebugLogf("unblocked")

	case obs.ClientDisconnection: // TODO: super inefficient ?
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
		n.syncWholeState()

	// Depending on settings, network system will "tag" every objects with an action to do for each clients
	case obs.ClientSettingsUpdate:
		n.updateClientAction(e.ClientToken, isClientDebugging(e.Settings.UpdateParameters.Parameters),
			getCullingArea(e.Settings.UpdateParameters.Parameters))
		//n.syncWholeState()
	}
}

// Request a sync of a single object to clients
func (n *NetworkSystem) syncSingleObject(object octree.Object) {
	//utils.DebugLogf("block")

	if no, ok := object.Data.(*networkObject); ok {
		for keyClient, clientValue := range no.clientsAction { // TODO: well done its probably possible to handle each client concurrently
			// Get channel
			switch clientValue {
			case ignore: // Ignore is continuous
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
				// If he is owner of the object, we don't send him updates (easy cheating)
				if isOwner {
				} else if !isOwner {
					//utils.DebugLogf("sync %v", object.ID())
					ManagerInstance.Send(keyClient, erutan.Packet{
						Type: &erutan.Packet_UpdateObject{
							UpdateObject: &erutan.Packet_UpdateObjectPacket{
								ObjectId:   object.ID(),
								Components: no.components,
							},
						},
					})
					no.clientsAction[keyClient] = update
				}

			case hide: // Destroy is discrete, only do it once then ignore
				//for _, c := range no.components {
				//	switch c := c.Type.(type) {
				//	case *erutan.Component_NetworkBehaviour:
				//		if c.NetworkBehaviour.OwnerToken == keyClient {
				//			utils.DebugLogf("Hiding my own object is forbidden !!!!!")
				//		}
				//	}
				//}
				ManagerInstance.Send(keyClient, erutan.Packet{
					Type: &erutan.Packet_DestroyObject{
						DestroyObject: &erutan.Packet_DestroyObjectPacket{
							ObjectId: object.ID(),
						},
					},
				})
				no.clientsAction[keyClient] = ignore
			}
			no.lastUpdate = utils.GetProtoTime()
		}
	}
	//utils.DebugLogf("block")

}

// Request a sync of the whole state to clients
func (n *NetworkSystem) syncWholeState() {
	for _, object := range n.objects.GetAllObjects() {
		//n.objects.Range(func(object *octree.Object) bool {
		n.syncSingleObject(object)
		//return true
		//})
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
					//utils.DebugLogf("obj owner %v, client %v", c.NetworkBehaviour.OwnerToken, clientToken)
					if c.NetworkBehaviour.OwnerToken == clientToken {
						no.clientsAction[clientToken] = update
					} else {
						// Otherwise filter out objects outside client's culling area
						if cullingArea != nil && object.Bounds.Fit(*cullingArea) {
							// Filter out / in debug objects
							if c.NetworkBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_ALL {
								// Object isn't tagged with debug network behaviour, just update
								no.clientsAction[clientToken] = update
							} else if c.NetworkBehaviour.Tag == erutan.Component_NetworkBehaviourComponent_DEBUG { // Object is a debug thing
								if isClientDebugging { // Is the client debugging ?
									no.clientsAction[clientToken] = update
								} else if no.clientsAction[clientToken] != ignore {
									no.clientsAction[clientToken] = hide // Hide once
								}
							}
						} else if no.clientsAction[clientToken] != ignore { // Destroy client's objects that have just been out of culling area
							no.clientsAction[clientToken] = hide
						}
					}
					//utils.DebugLogf("display %v", no.clientsAction[clientToken])
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
			//utils.DebugLogf("new culling area %v", p.CullingArea)
			return p.CullingArea
		}
	}
	return nil
}
