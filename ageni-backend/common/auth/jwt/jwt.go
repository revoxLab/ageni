package jwt

import (
	xjwt "github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

var Jwt *_jwt

type _jwt struct{}

type Config struct {
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Expire    int64  `json:"expire" yaml:"expire"`
}

func (*_jwt) Encode(claim xjwt.Claims, config *Config) (string, error) {
	token := xjwt.NewWithClaims(xjwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(config.SecretKey))
}

func (j *_jwt) Secret(c *Config) xjwt.Keyfunc {
	return func(token *xjwt.Token) (interface{}, error) {
		return []byte(c.SecretKey), nil
	}
}

func (*_jwt) Decode(tokenString string, config *Config, claim *AccountClaims) error {
	if tokenString == "" {
		return errors.New("Token Empty")
	}

	_, err := xjwt.ParseWithClaims(tokenString, claim, func(token *xjwt.Token) (interface{}, error) {
		return []byte(config.SecretKey), nil
	})

	if err != nil {
		var ve *xjwt.ValidationError
		if errors.As(err, &ve) {
			if ve.Errors&xjwt.ValidationErrorExpired != 0 {
				return errors.New("Token Expired")
			} else if ve.Errors&xjwt.ValidationErrorMalformed != 0 {
				return errors.New("that's not even a token")
			} else if ve.Errors&xjwt.ValidationErrorNotValidYet != 0 {
				return errors.New("token not active yet")
			}
		}

		return errors.Wrap(err, "couldn't handle this token")
	}

	return nil
}
