package entity

import jwt "github.com/dgrijalva/jwt-go"

type ResponseStatus struct {
	IsSucccess bool   `json:"issuccess"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Stack      string `json:"stack"`
}

type Color struct {
	R          int  `json:"r"`
	G          int  `json:"g"`
	B          int  `json:"b"`
	IsSelected bool `json:"-"`
}

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"fname"`
	LastName  string `json:"lname"`
	Email     string `json:"email"`
	jwt.StandardClaims
}

type UserRaundStepNumber struct {
	Step int `json:"step"`
}

type Token struct {
	AccessToken string `json:"token"`
}

type IEntity interface {
	Do()
}

func (color Color) Do() {}
