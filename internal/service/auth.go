package service

import (
	"MyStonks-go/internal/common/redisclient"
	"MyStonks-go/internal/service/schema"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
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

func GenerateNonce() (string, error) {
	for range 5 {
		nonce := make([]byte, 32)
		_, err := rand.Read(nonce)
		if err != nil {
			return "", err
		}
		nonceStr := base64.StdEncoding.EncodeToString(nonce)

		success, err := redisclient.GetClient().SetNX(context.Background(), "nonce:"+nonceStr, 1, 5*time.Minute).Result()
		if err != nil {
			return "", err
		}

		if success {
			return nonceStr, nil
		}
	}
	return "", errors.New("failed to generate nonce")
}

func VerifySolanaWalletSignature(req *schema.LoginReq) (bool, error) {
	pubKey, err := solana.PublicKeyFromBase58(req.Address)
	if err != nil {
		log.Error().Err(err).Send()
		return false, ErrInvalidAddress
	}

	signature, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		log.Error().Err(err).Send()
		return false, ErrInvalidSignature
	}

	if len(signature) != 64 {
		log.Error().Msg("signature length is not 64")
		return false, ErrInvalidSignature
	}

	_, err = redisclient.GetClient().Get(context.Background(), "nonce:"+req.Nonce).Result()
	if err != nil {
		log.Error().Err(err).Send()
		return false, ErrInvalidNonce
	}
	err = redisclient.GetClient().Del(context.Background(), "nonce:"+req.Nonce).Err()
	if err != nil {
		log.Error().Err(err).Send()
		return false, ErrInvalidNonce
	}

	message := []byte(req.Nonce)
	var sig solana.Signature
	copy(sig[:], signature)
	valid := pubKey.Verify(message, sig)
	if !valid {
		return false, ErrInvalidSignature
	}

	return true, nil
}

func GenerateTokenPair(walletAddress string) (*schema.TokenPair, error) {
	accessToken, err := generateToken(walletAddress, TokenTypeAccess, 24*time.Hour)
	if err != nil {
		return nil, err
	}

	refreshToken, err := generateToken(walletAddress, TokenTypeRefresh, 7*24*time.Hour)
	if err != nil {
		return nil, err
	}

	return &schema.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
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

func RefreshToken(refreshToken string) (*schema.TokenPair, error) {
	claims, err := ValidateToken(refreshToken, TokenTypeRefresh)
	if err != nil {
		return nil, err
	}

	accessToken, err := generateToken(claims.WalletAddress, TokenTypeAccess, 24*time.Hour)
	if err != nil {
		return nil, err
	}
	return &schema.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func SetTokenBlacklist(token string, expiration time.Duration) error {
	ctx := context.Background()
	err := redisclient.GetClient().Set(ctx, "blacklist:"+token, 1, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置token黑名单失败: %w", err)
	}
	return nil
}

func IsTokenInBlacklist(token string) bool {
	ctx := context.Background()
	_, err := redisclient.GetClient().Get(ctx, "blacklist:"+token).Result()
	return err == nil
}
