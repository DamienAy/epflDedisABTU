package Timestamp

import "testing"
import (
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestamp"
	"log"
)


func TestTheStruct(t *testing.T) {
	t1 := NewTimestamp(3)
	log.Println(t1.Size())
	s := t1.Size()
	s++
	log.Println("lol")
	log.Println(t1.Size())
}
/*func TestHappenedBeforeReturnsFalseWithIdenticalTimestamps(t *testing.T) {
	var t1 Timestamp = [3]int{1, 2, 3}
	var t2 Timestamp = [3]int{1, 2, 3}

	if t1.HappenedBefore(t2) {
		t.Error("Expected false")
	}
}

func TestHappenedBeforeReturnsCorrectlyWithDifferentTimestamps(t *testing.T) {
	var t1 Timestamp = [3]int{1, 2, 3}
	var t2 Timestamp = [3]int{4, 2, 3}
	var t3 Timestamp = [3]int{1, 4, 3}


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
	var t1 Timestamp = [3]int{1, 2, 3}
	var t2 Timestamp = [3]int{2, 2, 3}

	t1.Increment(0)

	if t1 != t2 {
		t.Error("Increment should increment at index i.")
	}
}

func TestIsContainedInShouldReturnCorrectly(t *testing.T) {
	var t1 Timestamp = [3]int{1, 2, 3}
	var t2 Timestamp = [3]int{4, 2, 3}
	var t3 Timestamp = [3]int{1, 4, 3}
	var t4 Timestamp = [3]int{1, 1, 3}


	tSlice := []Timestamp{t1, t2, t3}

	if !t1.IsContainedIn(tSlice) || !t2.IsContainedIn(tSlice){
		t.Error("IsContainedIn Should Return True When Timestamp Contained")
	}

	if t4.IsContainedIn(tSlice) {
		t.Error("IsContainedIn Should return false when timestamp not contained in.")
	}


}*/
