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

	e.POST("/car/create", hndl.CreateCar)

	e.GET("/car/get", hndl.GetCar)

	e.PUT("/car/update", hndl.UpdateCar)

	e.DELETE("/car/delete", hndl.DeleteCar)

	e.Logger.Fatal(e.Start(":1323"))
}
