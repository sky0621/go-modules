package main

import (
	"fmt"

	"github.com/sky0621/go-modules"

	uuid "github.com/satori/go.uuid"
)

func main() {
	fmt.Println("Hello, World!")
	modules.SayBye()
	u1 := uuid.Must(uuid.NewV4())
	fmt.Printf("UUIDv4: %s\n", u1)
}
