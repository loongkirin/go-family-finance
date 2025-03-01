package oauth

type OAuthMaker interface {
	GenerateAccessToken(email, mobile, username string) (string, *OAuthClaims, error)
	GenerateRefreshToken(email, mobile, username string) (string, *OAuthClaims, error)
	VerifyToken(token string) (*OAuthClaims, error)
}
