package jwt

import (
	"fmt"
	"testing"
	"time"
)

type CustomData struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Value    int64
}

func TestGenerateJWTToken(t *testing.T) {
	data := CustomData{
		User:     "test",
		Password: "123456",
	}
	claims := NewCustomClaims(data, 5*time.Minute)

	str, err := GenerateJWTToken(claims)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(str)
}

func TestVerifyJWTToken(t *testing.T) {
	claims := NewCustomClaims(nil, 4*time.Second)
	str, err := GenerateJWTToken(claims)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(3 * time.Second)
	valid, err := VerifyJWTToken(str)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fail()
	}
}

func TestGetCustomData(t *testing.T) {
	data := CustomData{
		User:     "test",
		Password: "123456",
		Value:    time.Now().UnixNano(),
	}
	claims := NewCustomClaims(data, 5*time.Second)

	str, err := GenerateJWTToken(claims)
	if err != nil {
		t.Fatal(err)
	}
	dataVal := CustomData{}
	err = GetCustomData(str, &dataVal)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%v", dataVal)
}
