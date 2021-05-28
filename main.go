package main

import (
	"ginrss/model"
	"ginrss/routes"
)

func main() {
	model.InitDB()
	routes.InitRouter()
}
