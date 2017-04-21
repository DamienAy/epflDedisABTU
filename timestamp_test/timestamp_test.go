package timestamp_test

import "testing"
import . "github.com/DamienAy/epflDedisABTU/timestamp"


func TestHappenedBeforeReturnsFalseWithIdenticalTimestamps(t *testing.T) {
	var t1 Timestamp = [3]uint64{1, 2, 3}
	var t2 Timestamp = [3]uint64{1, 2, 3}

	if t1.HappenedBefore(t2) {
		t.Error("Expected false")
	}
}

func TestHappenedBeforeReturnsCorrectlyWithDifferentTimestamps(t *testing.T) {
	var t1 Timestamp = [3]uint64{1, 2, 3}
	var t2 Timestamp = [3]uint64{4, 2, 3}
	var t3 Timestamp = [3]uint64{1, 4, 3}


	if !t1.HappenedBefore(t2) {
		t.Error("Timestamp t1 happens before t2.")
	}

	if t2.HappenedBefore(t1){
		t.Error("Timestamp t1 happens before t2.")
	}

	if t2.HappenedBefore(t3) || t3.HappenedBefore(t2) {
		t.Error("Timestamp t2 and t3 are concurrent.")
	}
}

func TestIncrementWorks(t *testing.T) {
	var t1 Timestamp = [3]uint64{1, 2, 3}
	var t2 Timestamp = [3]uint64{2, 2, 3}

	t1.Increment(0)

	if t1 != t2 {
		t.Error("Increment should increment at index i.")
	}
}

func TestIsContainedInShouldReturnCorrectly(t *testing.T) {
	var t1 Timestamp = [3]uint64{1, 2, 3}
	var t2 Timestamp = [3]uint64{4, 2, 3}
	var t3 Timestamp = [3]uint64{1, 4, 3}
	var t4 Timestamp = [3]uint64{1, 1, 3}


	tSlice := []Timestamp{t1, t2, t3}

	if !t1.IsContainedIn(tSlice) || !t2.IsContainedIn(tSlice){
		t.Error("IsContainedIn Should Return True When Timestamp Contained")
	}

	if t4.IsContainedIn(tSlice) {
		t.Error("IsContainedIn Should return false when timestamp not contained in.")
	}


}
