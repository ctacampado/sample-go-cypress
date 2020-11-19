package main

import (
	"react-go-cypress/internal/service"
	"react-go-cypress/internal/todo"
)

func main() {
	svc, err := service.New()
	if err != nil {
		panic(err)
	}
	todoSvc := todo.Service{}
	todoSvc.Service = svc
	todoSvc.Service.Run(todoSvc.InitRoutes)
}
