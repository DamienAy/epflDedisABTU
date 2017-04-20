package main

import (
	t "github.com/DamienAy/epflDedisABTU/timestamp"
	"fmt"
)


func main() {

	t1 := t.Timestamp{0, 1, 2};
	t2 := t.Timestamp{1, 2, 3};

	fmt.Println(t2.HappenedBefore(t1))
}
