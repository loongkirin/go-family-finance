package oauth

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	util "github.com/loongkirin/go-family-finance/pkg/util"
)

var (
	ErrInvalidKey   = errors.New("key is invalid")
	ErrTokenExpired = errors.New("token is expired")
	ErrTokenInvalid = errors.New("token is invalid")
)

type OAuthClaims struct {
	Id        string    `json:"id"`
	Email     string    `json:"email"`
	Mobile    string    `json:"mobile"`
	UserName  string    `json:"user_name"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
	NotBefore time.Time `json:"not_before"`
	Issuer    string    `json:"issuer,omitempty"`
}

func (o OAuthClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(o.ExpiredAt), nil
}

func (o OAuthClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(o.IssuedAt), nil
}

func (o OAuthClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(o.NotBefore), nil
}

func (o OAuthClaims) GetIssuer() (string, error) {
	return o.Issuer, nil
}

func (o OAuthClaims) GetSubject() (string, error) {
	return "subject", nil
}

func (o OAuthClaims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{"audience"}, nil
}

func (o OAuthClaims) Valid() error {
	if o.ExpiredAt.Before(time.Now()) {
		return ErrTokenExpired
	}

	return nil
}

func NewOAuthClaims(email, mobile, userName, issuer string, duration time.Duration) *OAuthClaims {
	claims := &OAuthClaims{
		Id:        util.GenerateId(),
		Email:     email,
		Mobile:    mobile,
		UserName:  userName,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
		NotBefore: time.Now().Add(time.Second * -10),
		Issuer:    issuer,
	}

	return claims
}
