package Timestamp

import (
	. "github.com/DamienAy/epflDedisABTU/singleTypes"
	"errors"
)

type Timestamp struct {
	time []uint64
	size uint64
}

// Returns a new Timestamp of size size.
func NewTimestamp(size uint64) Timestamp {
	return Timestamp{make([]uint64, size), size}
}

// Returns a copy of the time of the Timestamp
func (t *Timestamp) Time() []uint64 {
	time := make([]uint64, t.size)
	copy(time, t.time)
	return time
}

// Returns the size of the Timestamp
func (t *Timestamp) Size() uint64 {
	return t.size
}


// Increments the operation counter for site with id sId.
// Returns an error if
func (t *Timestamp) Increment(siteId SiteId) error {
	//type conversion possible, siteId is a uint64
//-------- This implementation can change in the future, SiteID might become a string
	if uint64(siteId) >= t.size {
		return errors.New("Cannot increment for siteId bigger than the size of the Timestamp.")
	} else {
		t.time[siteId]++
		return nil
	}
}

// Returns true if and only if the Timestamp t1 happened before the timmestamp t2.
// Returns an error if t1.size != t2.size.
func (t1 *Timestamp) HappenedBefore(t2 Timestamp) (bool, error) {
	if t1.size != t2.size {
		return false, errors.New("Two Timestamps of different lenght cannot be compared")
	} else {
		happenedBefore := false
		for index, element := range t1.time {
			if element > t2.time[index] {
				return false, nil
			} else if element < t2.time[index] {
				happenedBefore = true
			}
		}

		return happenedBefore, nil
	}
}

// Returns true if and only if t1.time is equal to t2.time.
// Returns an error if t1.size is not equal to t2.size.
func (t1 *Timestamp) Equals(t2 Timestamp) (bool, error) {
	if t1.size!=t2.size {
		return false, errors.New("Two Timestamps of different lenght cannot be compared")
	} else {
		p := true
		for index, el1 := range t1.time {
			p = p && (el1 == t2.time[index])
		}

		return p, nil
	}
}

// Returns true if and only if the Timestamp t is contained in the Timestamp slice tSlice.
// Returns an error if two compared Timestamps have different length.
func (t *Timestamp) IsContainedIn(tSlice []Timestamp) (bool, error) {
	for _, t2 := range tSlice {
		equals, err := t.Equals(t2)

		if err != nil {
			return false, err
		} else if equals {
			return true, nil
		}
	}

	return false, nil
}

// Returns true if and only if the intersection between the two slices of Timestamps is not empty.
// Returns an error if two compared Timestamps have different length.
func IntersectionIsNotEmpty(tSlice1 []Timestamp, tSlice2 []Timestamp) (bool, error) {
	for _, t1 := range tSlice1 {
		isContainedIn, err := t1.IsContainedIn(tSlice2)

		if err!=nil {
			return false, err
		} else if isContainedIn {
			return true, nil
		}
	}

	return false, nil
}
