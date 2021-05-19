package main

import (
	"CarsClient/handler"
	"bytes"
	"encoding/json"
	"fmt"
	protocol "github.com/Sirok47/CarsServer/protocol"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"google.golang.org/grpc"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	body []byte
	//err error
	token string
	cli   *http.Client
)

//req.Header.Set("Authorization",fmt.Sprintf("%s %s","Bearer",token))
func TestMain(m *testing.M) {
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
	go func() { e.Logger.Fatal(e.Start(":1323")) }()
	time.Sleep(100 * time.Millisecond)
	cli = &http.Client{}
	t := jwt.New(jwt.SigningMethodHS256)
	claims := t.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour).Unix()
	token, err = t.SignedString([]byte("sirok"))
	if err != nil {
		os.Exit(1)
	}
	exitCode := m.Run()

	os.Exit(exitCode)
}

func TestLogIn(t *testing.T) {
	reqbody, _ := json.Marshal(
		map[string]string{
			"Nick":     "keklik",
			"Password": "qpwoeirutyM123",
		})
	req, err := http.NewRequest("GET", "http://localhost:1323/user/login", bytes.NewBuffer(reqbody))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	strbody := string(body)
	if strbody == "rpc error: code = Unknown desc = code=401, message=Unauthorized" {
		t.Errorf("Expected token, got %v", string(body))
	}
}

//for{
//bs := make([]byte, 1014)
//n, err := resp.Body.Read(bs)
//strbody=fmt.Sprintf("%s%s",strbody,string(bs[:n]))
//if n == 0 || err != nil{
//break
//}
//}
