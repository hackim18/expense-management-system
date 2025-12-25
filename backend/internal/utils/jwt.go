package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type JWTHelper struct {
	secret    []byte
	issuer    string
	audience  string
	expiresIn time.Duration
}

func NewJWT(v *viper.Viper) *JWTHelper {
	expiresMinutes := v.GetInt("JWT_EXPIRES_MINUTES")
	if expiresMinutes <= 0 {
		expiresMinutes = 60
	}

	return &JWTHelper{
		secret:    []byte(v.GetString("JWT_SECRET")),
		issuer:    v.GetString("JWT_ISSUER"),
		audience:  v.GetString("JWT_AUDIENCE"),
		expiresIn: time.Duration(expiresMinutes) * time.Minute,
	}
}

type AccessClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func (j *JWTHelper) GenerateAccessToken(id uuid.UUID, email string) (string, error) {
	if len(j.secret) == 0 {
		return "", errors.New("jwt secret is empty")
	}

	now := time.Now()
	claims := AccessClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   id.String(),
			Issuer:    j.issuer,
			Audience:  jwt.ClaimStrings{j.audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func (j *JWTHelper) DecodeAccessToken(tokenStr string) (*AccessClaims, error) {
	if len(j.secret) == 0 {
		return nil, errors.New("jwt secret is empty")
	}

	claims := &AccessClaims{}
	parserOptions := []jwt.ParserOption{
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	}
	if j.audience != "" {
		parserOptions = append(parserOptions, jwt.WithAudience(j.audience))
	}
	if j.issuer != "" {
		parserOptions = append(parserOptions, jwt.WithIssuer(j.issuer))
	}

	parser := jwt.NewParser(parserOptions...)
	token, err := parser.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if token == nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
