package main

import (
	"fmt"
	"naboobase/proto_struct"
)

func main() {

	fmt.Println(&proto_struct.Person{
		Email: "lodjim5@gmail.com",
		Id:    23,
		Name:  "Hello",
	})

}
