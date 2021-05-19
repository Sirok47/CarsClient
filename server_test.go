package main

import (
	"CarsClient/handler"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	protocol "github.com/Sirok47/CarsServer/protocol"
	"github.com/Sirok47/CarsServer/service"
	"github.com/dgrijalva/jwt-go"
	"github.com/jackc/pgx/v4"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"google.golang.org/grpc"
	"io"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

var (
	dbconn *pgx.Conn
	body   []byte
	err    error
	token  string
	cli    *http.Client
)

func TestMain(m *testing.M) {
	dbconn, err = pgx.Connect(context.Background(), "postgres://maks:glazirovanniisirok@127.0.0.1:5432/cars")
	if err != nil {
		fmt.Print(err)
		return
	}
	srv := grpc.NewServer()
	srvobj := service.NewService(dbconn)
	protocol.RegisterCarsServer(srv, srvobj)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Print(err)
	}
	go func() { err = srv.Serve(l) }()
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

func TestSignUp(t *testing.T) {
	defer dbconn.Exec(context.Background(), "delete from users")
	reqbody, _ := json.Marshal(
		map[string]string{
			"Nick":     "keklik",
			"Password": "qpwoeirutyM123",
		})
	req, err := http.NewRequest("POST", "http://localhost:1323/user/signup", bytes.NewBuffer(reqbody))
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
	if strbody != "New user added" {
		t.Errorf("Got error: %v", strbody)
	}
}

func TestLogIn(t *testing.T) {
	dbconn.Exec(context.Background(), "insert into users (nick, password) values ($1,$2)", "keklik", "qpwoeirutyM123")
	defer dbconn.Exec(context.Background(), "delete from users")
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
		t.Errorf("Expected token, got %v", strbody)
	}
}

func TestCreate(t *testing.T) {
	defer dbconn.Exec(context.Background(), "delete from cars")
	reqbody, _ := json.Marshal(
		map[string]interface{}{
			"CarBrand":  "test",
			"CarNumber": 1234,
			"Type":      "test",
			"Mileage":   1000,
		})
	req, err := http.NewRequest("POST", "http://localhost:1323/car/create", bytes.NewBuffer(reqbody))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	strbody := string(body)
	if strbody != "Car have been created" {
		t.Errorf("Got error: %v", strbody)
	}
}

func TestGet(t *testing.T) {
	dbconn.Exec(context.Background(), "insert into cars (carbrand,carnumber,type,mileage) values ($1,$2,$3,$4)", "test", 1234, "test", 1000)
	defer dbconn.Exec(context.Background(), "delete from cars")
	reqbody, _ := json.Marshal(
		map[string]interface{}{
			"CarNumber": 1234,
		})
	req, err := http.NewRequest("GET", "http://localhost:1323/car/get", bytes.NewBuffer(reqbody))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	strbody := string(body)
	if strbody != "{\"CarBrand\":\"test\",\"CarNumber\":1234,\"Mileage\":1000,\"CarType\":\"test\"}\n" {
		t.Errorf("Expected {\"CarBrand\":\"test\",\"CarNumber\":1234,\"Mileage\":1000,\"CarType\":\"test\"}\n got: %v.", strbody)
	}
}

func TestUpdate(t *testing.T) {
	dbconn.Exec(context.Background(), "insert into cars (carbrand,carnumber,type,mileage) values ($1,$2,$3,$4)", "test", 1234, "test", 1111)
	defer dbconn.Exec(context.Background(), "delete from cars")
	reqbody, _ := json.Marshal(
		map[string]interface{}{
			"CarNumber": 1234,
			"mileage":   1112,
		})
	req, err := http.NewRequest("PUT", "http://localhost:1323/car/update", bytes.NewBuffer(reqbody))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	strbody := string(body)
	if strbody != "Car updated" {
		t.Errorf("Got error %v", strbody)
	}
}
func TestDelete(t *testing.T) {
	dbconn.Exec(context.Background(), "insert into cars (carbrand,carnumber,type,mileage) values ($1,$2,$3,$4)", "brand", 1234, "type", 1111)
	defer dbconn.Exec(context.Background(), "delete from cars")
	reqbody, _ := json.Marshal(
		map[string]interface{}{
			"CarNumber": 1234,
		})
	req, err := http.NewRequest("DELETE", "http://localhost:1323/car/delete", bytes.NewBuffer(reqbody))
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", "Bearer", token))
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	strbody := string(body)
	if strbody != "Car deleted" {
		t.Errorf("Got error %v", strbody)
	}
}
