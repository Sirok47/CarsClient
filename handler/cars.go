package handler

import (
	"context"
	"github.com/Sirok47/CarsServer/model"
	protocol "github.com/Sirok47/CarsServer/protocol"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Cars struct {
	client protocol.CarsClient
}

func NewCars(client protocol.CarsClient) *Cars {
	return &Cars{client: client}
}

var valid *validator.Validate

// SignUp godoc
// @Summary Create new user
// @Description Creates new user using nick and password
// @ID sign-up
// @Accept json
// @Produce plain
// @Param user body model.User true "User data"
// @Success 201 {object} string
// @Failure 500 {object} string
// @Router /user/signup [post]
func (h *Cars) SignUp(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(&user); err != nil {
		return err
	}
	if err := valid.Struct(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err, _ := h.client.SignUp(context.Background(), &protocol.Userdata{Nick: user.Nick, Password: user.Password})
	if err.Error != "" {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusCreated, "New user added")
}

//LogIn godoc
//@Summary Get token
//@Description Returns token using your nick and password
//@ID log-in
//@Accept json
//@Produce plain
//@Param user body model.User true "User data"
//@Success 201 {object} string
//@Failure 500 {object} string
//@Router /user/login [get]
func (h *Cars) LogIn(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return err
	}
	if err := valid.Struct(user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	token, err := h.client.LogIn(context.Background(), &protocol.Userdata{Nick: user.Nick, Password: user.Password})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusOK, token.Token)
}

//Create godoc
//@Summary Create new car
//@Description Creates new car using number, brand, type and mileage
//@ID create-car
//@Accept json
//@Produce plain
//@Param car body model.Car true "Car data"
//@Success 201 {object} string
//@Failure 500 {object} string
//@Router /car/create [post]
func (h *Cars) Create(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
	}
	if err := valid.Struct(car); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err, _ := h.client.Create(context.Background(), &protocol.Carparams{CarBrand: car.CarBrand, Mileage: int32(car.Mileage), CarType: car.CarType, CarNumber: int32(car.CarNumber)})
	if err.Error != "" {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusCreated, "Car have been created")
}

//Get godoc
//@Summary Get car data
//@Description Gets car data by number
//@ID get-car
//@Accept json
//@Produce json
//@Param car body model.Car true "Car data"
//@Success 201 {object} model.Car
//@Failure 500 {object} string
//@Router /car/get [get]
func (h *Cars) Get(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
	}
	if err := valid.Struct(car); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	carInfo, err := h.client.Get(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(200, carInfo)
}

//Update godoc
//@Summary Update car data
//@Description Replaces mileage with new by number
//@ID update-car
//@Accept json
//@Produce plain
//@Param car body model.Car true "Car data"
//@Success 201 {object} string
//@Failure 500 {object} string
//@Router /car/update [put]
func (h *Cars) Update(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
	}
	if err := valid.Struct(car); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err, _ := h.client.Update(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber), Mileage: int32(car.Mileage)})
	if err.Error != "" {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusOK, "Car updated")
}

//Delete godoc
//@Summary Delete car
//@Description Deletes car by its number
//@ID delete-car
//@Accept json
//@Produce plain
//@Param car body model.Car true "Car data"
//@Success 201 {object} string
//@Failure 500 {object} string
//@Router /car/delete [delete]
func (h *Cars) Delete(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
	}
	if err := valid.Struct(car); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	err, _ := h.client.Delete(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber)})
	if err.Error != "" {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusOK, "Car deleted")

}
