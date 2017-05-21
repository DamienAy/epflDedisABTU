package Timestamp

import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypes"
	"errors"
)

// Represents a timestamp as described in the ABTU paper.
type Timestamp struct {
	time []uint64
	size uint64
}

// Returns a new Timestamp of size size.
func NewTimestamp(size uint64) Timestamp {
	return Timestamp{make([]uint64, size), size}
}

// Returns a deep copy of the timestamp t
func DeepCopyTimestamp(t Timestamp) Timestamp {
	time := make([]uint64, t.size)
	copy(time, t.time)
	return Timestamp{time, t.size}
}

// Returns a deep copy of the slice of Timestamp timestamps.
func DeepCopyTimestamps(timestamps []Timestamp) []Timestamp {
	timestampsCopy := make([]Timestamp, len(timestamps))

	for i, t := range timestamps {
		timestampsCopy[i] = DeepCopyTimestamp(t)
	}

	return timestampsCopy
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
// Returns an error if siteId is >= than the size of the timestamp.
// TODO adjust implementation for siteId and return error (siteIds might as well be strings)
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
		return false, errors.New("Two Timestamps of different length cannot be compared")
	} else {
		p := true
		for index, el1 := range t1.time {
			p = p && (el1 == t2.time[index])
		}

		return p, nil
	}
}

// Returns true if and only if the Timestamp t1 is causally ready at time currentSV and site siteId.
// Returns and error if t1 ans currentSV have different sizes.
// Does not modify t1 nor currentSV.

func (t1 *Timestamp) IsCausallyReady(currentSV Timestamp, siteId SiteId) (bool, error) {
	if t1.size!=currentSV.size {
		return false, errors.New("Two Timestamps of different length cannot be compared")
	} else {
		t := t1.time
		sv := currentSV.time

		isCausallyReady := t[siteId] == sv[siteId]

		for k:=uint64(0); k<t1.size; k++ {
			if k != uint64(siteId) {
				isCausallyReady = isCausallyReady && t[k] <= sv[k]
			}
		}

		return isCausallyReady, nil
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

// Same as Timestamp but with public fields.
// Used for encoding into Json objects.
type PublicTimestamp struct {
	Time []uint64
	Size uint64
}

// Returns the PublicTimestamp corresponding to the Timestamp t.
func TimestampToPublicTimestamp(t Timestamp) PublicTimestamp {
	timestamp := DeepCopyTimestamp(t)
	return PublicTimestamp{timestamp.time, timestamp.size}
}

// Returns the slice of PublicTimestamp corresponding to the slice of Timestamp timestamps.
func TimestampsToPublicTimestamps(timestamps []Timestamp) []PublicTimestamp {
	publicTimestamps := make([]PublicTimestamp, len(timestamps))
	for i, t := range timestamps {
		publicTimestamps[i] = TimestampToPublicTimestamp(t)
	}

	return publicTimestamps
}

// Returns the Timestamp corresponding to the PublicTimestamp publicT
func PublicTimestampToTimestamp(publicT PublicTimestamp) Timestamp {
	return DeepCopyTimestamp(Timestamp{publicT.Time, publicT.Size})
}

// Returns the slice of Timestamp corresponding to the slice of PublicTimestamp publicTimestamps.
func PublicTimestampsToTimestamps(publicTimestamps []PublicTimestamp) []Timestamp {
	timestamps := make([]Timestamp, len(publicTimestamps))
	for i, t := range publicTimestamps {
		timestamps[i] = PublicTimestampToTimestamp(t)
	}

	return timestamps
}
