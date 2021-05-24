package handler

import (
	"context"
	"github.com/Sirok47/CarsServer/model"
	protocol "github.com/Sirok47/CarsServer/protocol"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Cars struct {
	client protocol.CarsClient
}

func NewCars(client protocol.CarsClient) *Cars {
	return &Cars{client: client}
}

// SignUp godoc
// @Summary Create new user
// @Description Creates new user using nick and password
// @ID sign-up
// @Accept json
// @Produce plain
// @Param userobj body model.Userdata true "User data"
// @Router /user/signup [post]
// @Success 201 {object} c.String(http.StatusCreated, "New user added")
// @Failure 500 {object} c.String(http.StatusInternalServerError, err.Error)

func (h *Cars) SignUp(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return err
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
//@ID sign-up
//@Accept json
//@Produce plain
//@Param userobj body model.Userdata true "User data"
//@Router /user/signup [post]
//@Success 201 {object} c.String(http.StatusOK, token.Token)
//@Failure 500 {object} c.String(http.StatusInternalServerError, err.Error)

func (h *Cars) LogIn(c echo.Context) error {
	user := &model.User{}
	if err := c.Bind(user); err != nil {
		return err
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
//@Param userobj body model.Car true "Car data"
//@Router /car/create [post]
//@Success 201 {object} c.String(http.StatusCreated, "Car have been created")
//@Failure 500 {object} c.String(http.StatusInternalServerError, err.Error)

func (h *Cars) Create(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
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
//@Param userobj body model.Car true "Car data"
//@Router /car/get [get]
//@Success 201 {object} c.JSON(http.StatusOK, carInfo)
//@Failure 500 {object} c.String(http.StatusInternalServerError, err.Error)

func (h *Cars) Get(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
	}
	carInfo, err := h.client.Get(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber)})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, carInfo)
}

//Update godoc
//@Summary Update car data
//@Description Replaces mileage with new by number
//@ID update-car
//@Accept json
//@Produce plain
//@Param userobj body model.Car true "Car data"
//@Router /car/update [put]
//@Success 201 {object} c.String(http.StatusOK, "Car updated")
//@Failure 500 {object} c.String(http.StatusInternalServerError, err.Error)

func (h *Cars) Update(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
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
//@Tags root
//@Param userobj body model.Car true "Car data"
//@Router /car/delete [delete]
//@Success 201 {object} c.String(http.StatusOK, "Car deleted")
//@Failure 500 {object} c.String(http.StatusInternalServerError, err.Error)

func (h *Cars) Delete(c echo.Context) error {
	car := &model.Car{}
	if err := c.Bind(car); err != nil {
		return err
	}
	err, _ := h.client.Delete(context.Background(), &protocol.Carparams{CarNumber: int32(car.CarNumber)})
	if err.Error != "" {
		return c.String(http.StatusInternalServerError, err.Error)
	}
	return c.String(http.StatusOK, "Car deleted")

}
