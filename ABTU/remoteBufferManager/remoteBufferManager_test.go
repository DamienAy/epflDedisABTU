package remoteBufferManager

import (
	"testing"
	. "github.com/DamienAy/epflDedisABTU/ABTU/operation"
	_ "github.com/DamienAy/epflDedisABTU/ABTU/remoteBufferManager"
	"fmt"
	_"encoding/json"
	. "github.com/DamienAy/epflDedisABTU/ABTU/singleTypesTypes"
	. "github.com/DamienAy/epflDedisABTU/ABTU/timestampstamp"
)


type Response1 struct {
	Page   int
	Fruits []string
}

func TestJustTryIt ( *testing.T) {
	//operation := PartialOperation(0, 0, 2, 0)

	/*var rbm RemoteBufferManager

	rbm.Start(make([]Operation,0 ))

	answer := make(chan bool)
	answer2 := make(chan bool)
	rbm.Add <- AddOp{operation, answer}
	<-answer
	rbm.Add <- AddOp{operation, answer2}
	<-answer2

	ret := make(chan []Operation)
	rbm.Get <- GetOp{ret}
	fmt.Println(<-ret)
	rbm.RemoveRearrange <- RemoveRearrangeOp{1, answer }
	<-answer


	rbm.Get <- GetOp{ret}
	fmt.Println(<-ret)
	rbm.RemoveRearrange <- RemoveRearrangeOp{0, answer }
	<-answer

	rbm.Get <- GetOp{ret}
	fmt.Println(<-ret)*/

	/*js, _ := json.Marshal([]string{"lol", "lol2"})
	fmt.Println(string(js))
	fmt.Println(json.Unmarshal(json.Marshal(operation)))

	res1D := &Response1{
		Page:   1,
		Fruits: []string{"apple", "peach", "pear"}}
	res1B, _ := json.Marshal(res1D)

	fmt.Println(string(res1B))
	public := OperationToPublicOp(operation)
	js, _ = json.Marshal(&public)

	fmt.Println(string(js))
	var public3 map[string]interface{}
	json.Unmarshal(js, &public3)
	fmt.Println(public3["Id"])
	fmt.Println(public3)*/

	fmt.Println(INS)

}


