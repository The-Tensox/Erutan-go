package erutan

import (
	"github.com/The-Tensox/octree"
	"github.com/The-Tensox/protometry"
)

type (
	// ClientSettingsUpdate notify of a client updating its settings
	ClientSettingsUpdate struct {
		ClientToken string
		Settings    Packet_UpdateParameters
	}
	ClientConnection struct {
		ClientToken string
		Settings    Packet_UpdateParameters
	}
	ClientDisconnection struct {
		ClientToken string
	}
	PhysicsUpdateRequest struct {
		// Object is the object to apply physics on associated with a new position, later could add rotation, scale...
		Object struct {
			octree.Object
			protometry.Vector3
		}
		Dt float64
	}
	PhysicsUpdateResponse struct {
		// Objects is a slice of object and their new position to be updated, if len > 1 it is the objects collided with
		Objects []struct {
			octree.Object
			protometry.Vector3
		}
		//Me  *octree.Object
		//NewPosition  protometry.Vector3
		//// Other is nil if there is no collision
		//Other  *octree.Object
		Dt float64
	}
)
