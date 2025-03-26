package oauth

type OAuthMaker interface {
	GenerateAccessToken(userId, email, phone, userName string) (string, *OAuthClaims, error)
	GenerateRefreshToken(userId, email, phone, userName string) (string, *OAuthClaims, error)
	VerifyToken(token string) (*OAuthClaims, error)
}
