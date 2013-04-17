package main

import (
	"fmt"
	"github.com/ogo/ogo"
	"github.com/demo"
)

func main() {
	fmt.Println("Ogo 2013")
	ctrl := ogo.NewController()
	ctrl.RegisterApplication( new(demo.DemoApp) )
	ctrl.Start(":6633")
}
