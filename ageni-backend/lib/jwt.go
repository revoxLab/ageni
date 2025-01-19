package lib

import (
	"github.com/golang-jwt/jwt/v4"
	xjwt "github.com/readonme/open-studio/common/auth/jwt"
	"github.com/readonme/open-studio/dal/model"
	"time"
)

func GenJwtToken(user *model.User, config *xjwt.Config) (string, error) {
	userClaim := &xjwt.AccountClaims{
		Account: xjwt.Account{
			Id:            user.ID,
			Name:          user.Name,
			WalletAddress: user.WalletAddress,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(config.Expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token, err := xjwt.Jwt.Encode(userClaim, config)
	if err != nil {
		return "", err
	}
	return token, nil
}
