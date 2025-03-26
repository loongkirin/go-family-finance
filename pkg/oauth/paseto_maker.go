package oauth

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto/v2"
	"golang.org/x/crypto/chacha20poly1305"
)

// PasetoMaker is a PASETO maker.
type PasetoMaker struct {
	paseto      *paseto.V2
	secretKey   []byte
	oauthConfig OAuthConfig
}

func (maker PasetoMaker) GenerateAccessToken(userId, email, phone, userName string) (string, *OAuthClaims, error) {
	duration, err := time.ParseDuration(maker.oauthConfig.AccessExpiresTime)
	if err != nil {
		duration = time.Hour * 8
	}

	claims := NewOAuthClaims(userId, email, phone, userName, maker.oauthConfig.Issuer, duration)
	token, err := maker.paseto.Encrypt(maker.secretKey, claims, nil)
	return token, claims, err
}

func (maker PasetoMaker) GenerateRefreshToken(userId, email, phone, userName string) (string, *OAuthClaims, error) {
	duration, err := time.ParseDuration(maker.oauthConfig.RefreshExpiresTime)
	if err != nil {
		duration = time.Hour * 24 * 7
	}

	claims := NewOAuthClaims(userId, email, phone, userName, maker.oauthConfig.Issuer, duration)
	token, err := maker.paseto.Encrypt(maker.secretKey, claims, nil)
	return token, claims, err
}

func (maker PasetoMaker) VerifyToken(token string) (*OAuthClaims, error) {
	claims := &OAuthClaims{}
	err := maker.paseto.Decrypt(token, maker.secretKey, claims, nil)
	if err != nil {
		return nil, ErrTokenInvalid
	}

	err = claims.Valid()
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// NewPasetoMaker creates a new PasetoMaker.
func NewPasetoMaker(oauthCfg OAuthConfig) (OAuthMaker, error) {
	if len(oauthCfg.SecretKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:      paseto.NewV2(),
		secretKey:   []byte(oauthCfg.SecretKey),
		oauthConfig: oauthCfg,
	}

	return maker, nil
}
