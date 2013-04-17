package main

import (
	"fmt"
	"github.com/jonstout/ogo/core"
)

func main() {
	fmt.Println("Ogo 2013")
	ctrl := ogo.NewController()
	ctrl.Start(":6633")
}
