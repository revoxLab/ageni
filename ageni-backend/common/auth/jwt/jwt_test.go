package jwt

import (
	xjwt "github.com/golang-jwt/jwt/v4"
	"testing"
	"time"
)

var config = &Config{
	SecretKey: "098f6bcd4621d373cade4e832627b4f6",
	Expire:    31536000,
}

func Test_jwt_Encode(t *testing.T) {

	tokenAccount := Account{
		Id:   157,
		Name: "readon",
	}

	userClaim := &AccountClaims{
		Account: tokenAccount,
		RegisteredClaims: xjwt.RegisteredClaims{
			ExpiresAt: xjwt.NewNumericDate(time.Now().Add(time.Duration(config.Expire) * time.Second)),
			IssuedAt:  xjwt.NewNumericDate(time.Now()),
			NotBefore: xjwt.NewNumericDate(time.Now()),
		},
	}

	token, err := Jwt.Encode(userClaim, config)
	t.Log(token)
	t.Log(err)
}

func Test_jwt_Decode(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBY2NvdW50Ijp7ImlkIjoxNTcsIm5hbWUiOiJyZWFkb24iLCJlbWFpbCI6InRlc3RAcmVhZG9uLmNvbSIsImFjY291bnRfdHlwZSI6MX0sImV4cCI6MTY5MDg4MTI2MCwibmJmIjoxNjU5MzQ1MjYwLCJpYXQiOjE2NTkzNDUyNjB9._-obDYUyJMtwhnsGFtc5chWlApWh8p8WvVEXXnCOI3I"
	userClaim := &AccountClaims{}
	err := Jwt.Decode(tokenString, config, userClaim)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(userClaim)
	t.Log(userClaim.Account.Id)
	t.Log(userClaim.Account.Name)
}
