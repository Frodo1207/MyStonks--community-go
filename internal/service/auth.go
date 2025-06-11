package service

import (
	"MyStonks-go/internal/common/redisclient"
	"context"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("无效的token")
	ErrExpiredToken     = errors.New("token已过期")
	ErrInvalidSignature = errors.New("无效的签名")
	ErrInvalidAddress   = errors.New("无效的钱包地址")
	ErrInvalidNonce     = errors.New("无效的nonce")
	secretKey           = []byte(os.Getenv("JWT_SECRET"))
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	WalletAddress string    `json:"wallet_address"`
	Type          TokenType `json:"type"`
	jwt.RegisteredClaims
}

func generateToken(walletAddress string, tokenType TokenType, expiration time.Duration) (string, error) {
	claims := Claims{
		WalletAddress: walletAddress,
		Type:          tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if IsTokenInBlacklist(tokenString) {
		return nil, ErrInvalidToken
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {

		if claims.Type != tokenType {
			return nil, ErrInvalidToken
		}
		return claims, nil
	}
	return nil, ErrInvalidToken
}

func IsTokenInBlacklist(token string) bool {
	ctx := context.Background()
	_, err := redisclient.GetClient().Get(ctx, "blacklist:"+token).Result()
	return err == nil
}
