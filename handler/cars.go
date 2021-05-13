package handler

import (
	"context"
	"github.com/Sirok47/CarsServer/model"
	"github.com/Sirok47/CarsServer/protocol"
	"github.com/labstack/echo"
	"net/http"
)

type Cars struct {
	client protocol.CarsClient
}

func NewCars(client protocol.CarsClient) *Cars {
	return &Cars{client: client}
}

func (h Cars) Create(c echo.Context) error {
	car := &model.CarParams{}
	if err := c.Bind(car); err != nil {
		return err
	}
	err, _ := h.client.Create(context.Background(), &protocol.Carparams{CarBrand: car.CarBrand, Mileage: int32(car.Mileage), CarType: car.CarType, CarNumber: int32(car.CarNumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusCreated, "Car have been created")
}

func (h Cars) Get(c echo.Context) error {
	car := &model.CarParams{}
	if err := c.Bind(car); err != nil {
		return err
	}
	carInfo, err := h.client.Get(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, carInfo)
}

func (h Cars) Update(c echo.Context) error {
	car := &model.CarParams{}
	if err := c.Bind(car); err != nil {
		return err
	}
	err, _ := h.client.Update(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber), Mileage: int32(car.Mileage)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusCreated, "Car updated")
}

func (h Cars) Delete(c echo.Context) error {
	car := &model.CarParams{}
	if err := c.Bind(car); err != nil {
		return err
	}
	err, _ := h.client.Delete(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusCreated, "Car deleted")

}
