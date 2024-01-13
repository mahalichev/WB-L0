package main

import (
	"github.com/mahalichev/WB-L0/api/app"
	"github.com/mahalichev/WB-L0/api/inits"
)

func init() {
	inits.LoadEnvironment()
}

func main() {
	app.RunService()
}
