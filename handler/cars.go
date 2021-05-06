package handler

import (
	"CarsServer/model"
	"CarsServer/protocol"
	"context"
	"github.com/labstack/echo"
	"net/http"
)

type Cars struct {
	client protocol.CarsClient
}

func NewCars(client protocol.CarsClient) *Cars {
	return &Cars{client: client}
}

func (h Cars) CreateCar(c echo.Context) error {
	car:=&model.CarParams{}
	if err := c.Bind(car); err != nil {
		return err
	}
	err, _ := h.client.CreateCar(context.Background(), &protocol.Carparams{Carbrand: car.Carbrand,Mileage: int32(car.Mileage),Cartype: car.Cartype,Carnumber: int32(car.Carnumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusCreated, "Car have been created")
}


func (h Cars) GetCar(c echo.Context) error {
	car:=&model.CarParams{}
	if err := c.Bind(car); err != nil {
		return err
	}
	carInfo, err := h.client.GetCar(context.Background(), &protocol.Carparams{Carnumber: int32(car.Carnumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, carInfo)
}

func (h Cars) UpdateCar(c echo.Context) error {
		car:=&model.CarParams{}
		if err := c.Bind(car); err != nil {
			return err
		}
		err, _ := h.client.UpdateCar(context.Background(), &protocol.Carparams{Carnumber: int32(car.Carnumber),Mileage: int32(car.Mileage)})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error)
		}
		return c.String(http.StatusCreated, "Car updated")
}

func (h Cars) DeleteCar(c echo.Context) error {
	car:=&model.CarParams{}
	if err := c.Bind(car); err != nil {
		return err
	}
	err, _ := h.client.DeleteCar(context.Background(), &protocol.Carparams{Carnumber: int32(car.Carnumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusCreated, "Car deleted")

}