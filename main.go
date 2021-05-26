package main

import (
	_ "CarsClient/docs"
	"fmt"
	"github.com/Sirok47/CarsClient/handler"
	protocol "github.com/Sirok47/CarsServer/protocol"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"
)

//@title Cars shop API
//@version 1.1
//@description This API lets managing cars DB
//@host localhost:1323
//@BasePath
//@schemes http

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

	e.Validator = &handler.CustomValidator{Valid: validator.New()}

	e.POST("/user/signup", hndl.SignUp)

	e.GET("/user/login", hndl.LogIn)

	e.POST("/car/create", hndl.Create, TokenValidation)

	e.GET("/car/get", hndl.Get, TokenValidation)

	e.PUT("/car/update", hndl.Update, TokenValidation)

	e.DELETE("/car/delete", hndl.Delete, TokenValidation)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
