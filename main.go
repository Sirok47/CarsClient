package main

import (
	"fmt"
	"github.com/Sirok47/CarsClient/handler"
	protocol "github.com/Sirok47/CarsServer/protocol"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
	"github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
)

func main() {
	TokenValidation := middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey: []byte("sirok"),
	})
	clientConnection, err := grpc.Dial(":8080", grpc.WithInsecure())
	client := protocol.NewCarsClient(clientConnection)
	if err != nil {
		fmt.Print(err)
		return
	}

	hndl := handler.NewCars(client)

	e := echo.New()

	e.POST("/user/signup", hndl.SignUp)

	e.GET("/user/login", hndl.LogIn)

	e.POST("/car/create", hndl.Create, TokenValidation)

	e.GET("/car/get", hndl.Get, TokenValidation)

	e.PUT("/car/update", hndl.Update, TokenValidation)

	e.DELETE("/car/delete", hndl.Delete, TokenValidation)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
