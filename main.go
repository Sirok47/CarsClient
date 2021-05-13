package main

import (
	"fmt"
	"github.com/Sirok47/CarsClient/handler"
	"github.com/Sirok47/CarsServer/protocol"
	"github.com/labstack/echo"
	"google.golang.org/grpc"
)

func main() {
	clientConnection, err := grpc.Dial(":8080", grpc.WithInsecure())
	client := protocol.NewCarsClient(clientConnection)
	if err != nil {
		fmt.Print(err)
		return
	}

	hndl := handler.NewCars(client)

	e := echo.New()

	e.POST("/car/create", hndl.Create)

	e.GET("/car/get", hndl.Get)

	e.PUT("/car/update", hndl.Update)

	e.DELETE("/car/delete", hndl.Delete)

	e.Logger.Fatal(e.Start(":1323"))
}
