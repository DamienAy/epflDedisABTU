package Timestamp

import "testing"
import (
	"log"
)

func TestHappenedBeforeReturnsFalseWithIdenticalTimestamps(t *testing.T) {
	var t1 Timestamp = NewTimestamp(3)
	t1.Increment(0)
	t1.Increment(1)
	t1.Increment(1)
	t1.Increment(2)
	t1.Increment(2)
	t1.Increment(2)

	var t2 Timestamp = NewTimestamp(3)
	t2.Increment(0)
	t2.Increment(1)
	t2.Increment(1)
	t2.Increment(2)
	t2.Increment(2)
	t2.Increment(2)

	test, _ := t1.HappenedBefore(t2)
	if test {
		t.Error("Expected false")
	}
}

func TestHappenedBeforeReturnsCorrectlyWithDifferentTimestamps(t *testing.T) {
	var t1 Timestamp = NewTimestamp(3)
	t1.Increment(0)
	t1.Increment(1)
	t1.Increment(1)
	t1.Increment(2)
	t1.Increment(2)
	t1.Increment(2)

	var t2 Timestamp = NewTimestamp(3)
	t2.Increment(0)
	t2.Increment(0)
	t2.Increment(0)
	t2.Increment(0)
	t2.Increment(1)
	t2.Increment(1)
	t2.Increment(2)
	t2.Increment(2)
	t2.Increment(2)

	var t3 Timestamp = NewTimestamp(3)
	t3.Increment(1)
	t3.Increment(2)
	t3.Increment(2)
	t3.Increment(2)
	t3.Increment(2)
	t3.Increment(3)
	t3.Increment(3)
	t3.Increment(3)

	test, _ := t1.HappenedBefore(t2)
	if !test {
		t.Error("Timestamp t1 happens before t2.")
	}
	test, _ = t2.HappenedBefore(t1)
	if test {
		t.Error("Timestamp t1 happens before t2.")
	}
	test, _= t2.HappenedBefore(t3)
	test2, _:= t3.HappenedBefore(t2)
	if test || test2 {
		t.Error("Timestamp t2 and t3 are concurrent.")
	}
}

func TestIsContainedInShouldReturnCorrectly(t *testing.T) {
	var t1 Timestamp = NewTimestamp(3)
	t1.Increment(0)
	t1.Increment(1)
	t1.Increment(1)
	t1.Increment(2)
	t1.Increment(2)
	t1.Increment(2)

	var t2 Timestamp = NewTimestamp(3)
	t2.Increment(0)
	t2.Increment(0)
	t2.Increment(0)
	t2.Increment(0)
	t2.Increment(1)
	t2.Increment(1)
	t2.Increment(2)
	t2.Increment(2)
	t2.Increment(2)

	var t3 Timestamp = NewTimestamp(3)
	t3.Increment(0)
	t3.Increment(1)
	t3.Increment(1)
	t3.Increment(1)
	t3.Increment(1)
	t3.Increment(2)
	t3.Increment(2)
	t3.Increment(2)

	var t4 Timestamp = NewTimestamp(2)
	t4.Increment(0)
	t4.Increment(1)
	t4.Increment(2)
	t4.Increment(2)
	t4.Increment(2)

	tSlice := []Timestamp{t1, t2, t3}

	test1, _ := t1.IsContainedIn(tSlice)
	test2, _:= t2.IsContainedIn(tSlice)

	if !test1 || !test2{
		t.Error("IsContainedIn Should Return True When Timestamp Contained")
	}

	test, _ := t4.IsContainedIn(tSlice)
	if test {
		t.Error("IsContainedIn Should return false when timestamp not contained in.")
	}


}
