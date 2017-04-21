package timestamp

import . "github.com/DamienAy/epflDedisABTU/constants"
import . "github.com/DamienAy/epflDedisABTU/singleTypes"

//import "operation"

type Timestamp [N]int

// Increments the operation counter for site with id sId.
func (t *Timestamp) Increment(siteId SiteId) {
	t[siteId]++
}

// Returns true if and only if the timestamp t1 happened before the timmestamp t2.
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

// Returns true if and only if the timestamp t is contained in the timestamp slice tSlice.
func (t *Timestamp) IsContainedIn(tSlice []Timestamp) bool {
	for _, element := range tSlice {
		if *t == element {
			return true
		}
	}

	return false
}
