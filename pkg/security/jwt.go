package security

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTConfig struct {
	SigningKey      []byte
	SigningMethod   jwt.SigningMethod
	ExpirationTime  time.Duration
	RefreshDuration time.Duration
	Issuer          string
	Audience        []string
}

type JWTManager struct {
	config JWTConfig
}

type Claims struct {
	jwt.StandardClaims
	UserID   string                 `json:"uid"`
	Roles    []string               `json:"roles"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func NewJWTManager(config JWTConfig) *JWTManager {
	if config.SigningMethod == nil {
		config.SigningMethod = jwt.SigningMethodHS256
	}
	return &JWTManager{config: config}
}

func (j *JWTManager) GenerateToken(claims Claims) (string, error) {
	now := time.Now()
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(j.config.ExpirationTime).Unix()
	claims.Issuer = j.config.Issuer
	claims.Audience = j.config.Audience

	token := jwt.NewWithClaims(j.config.SigningMethod, claims)
	return token.SignedString(j.config.SigningKey)
}
