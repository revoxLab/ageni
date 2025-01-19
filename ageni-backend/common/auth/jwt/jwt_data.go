package jwt

import (
	"encoding/json"
	xjwt "github.com/golang-jwt/jwt/v4"
)

type AccountClaims struct {
	Account Account
	Ticket  Ticket
	xjwt.RegisteredClaims
}

func (i *AccountClaims) String() string {
	b, _ := json.Marshal(i)
	return string(b)
}

type Account struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	WalletAddress string `json:"wallet_address"`
}

type Ticket struct {
	Key      string `json:"key"`
	Platform string `json:"platform"`
}

type AccountType int

const (
	AccountTypeBegin AccountType = iota
	AccountTypeReadOn
	AccountTypeTwitter
)
