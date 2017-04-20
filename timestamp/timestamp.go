package timestamp

import . "github.com/DamienAy/epflDedisABTU/constants"
import . "github.com/DamienAy/epflDedisABTU/singleTypes"

//import "operation"

type Timestamp [N]uint64

// Increments the operation counter for site with id sId.
func (t *Timestamp) increment(sId SiteId) {
	t[sId]++
}

func (t1 *Timestamp) HappenedBefore(t2 Timestamp) bool {
	happenedBefore := false
	for index, element := range t1 {
		if element > t2[index] {
			return false
		} else if element < t2[index] {
			happenedBefore = true
		}
	}

	return happenedBefore
}
