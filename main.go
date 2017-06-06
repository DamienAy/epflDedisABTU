package main

import (
	"github.com/DamienAy/epflDedisABTU/management"
)

func main() {
	/* Create an instance of Management and establish control communication with Frontend*/
	mgmt := management.NewManagement()
	mgmt.Run()
}
