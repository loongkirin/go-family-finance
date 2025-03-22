package oauth

type OAuthMaker interface {
	GenerateAccessToken(email, phone, username string) (string, *OAuthClaims, error)
	GenerateRefreshToken(email, phone, username string) (string, *OAuthClaims, error)
	VerifyToken(token string) (*OAuthClaims, error)
}
