package main

import (
	"github.com/mahalichev/WB-L0/app"
	"github.com/mahalichev/WB-L0/inits"
)

func init() {
	inits.LoadEnvironment()
}

func main() {
	app.RunService()
}
