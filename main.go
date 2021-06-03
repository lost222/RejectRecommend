package main

import (
	"ginrss/model"
	"ginrss/redismoon"
	"ginrss/routes"
)

func main() {
	model.InitDB()
	routes.InitRouter()
	Redismoon.Redisinit()
}
