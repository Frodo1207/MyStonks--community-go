package service

import (
	"MyStonks-go/internal/common/redisclient"
	"MyStonks-go/internal/models"
	"MyStonks-go/internal/routers/schema"

	"MyStonks-go/internal/store"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gagliardetto/solana-go"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"math/rand"
	"time"
)

type UserSrv struct {
	us store.UserStore
}

func NewUserSrv(store store.UserStore) *UserSrv {
	return &UserSrv{us: store}
}

func (u *UserSrv) CreateUserIfNotExists(solAddress string) error {
	return u.us.CreateUserIfNotExists(solAddress)
}

func (u *UserSrv) GetUserInfo(solAddress string) (*models.User, error) {
	user, err := u.us.GetUserBySolAddress(solAddress)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserSrv) GenerateNonce() (string, error) {
	for i := 0; i < 5; i++ {
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

func (u *UserSrv) VerifySolanaWalletSignature(req *schema.LoginReq) (bool, error) {
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

func (u *UserSrv) GenerateTokenPair(walletAddress string) (*schema.TokenPair, error) {
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

func (u *UserSrv) generateToken(walletAddress string, tokenType TokenType, expiration time.Duration) (string, error) {
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

func (u *UserSrv) ValidateToken(tokenString string, tokenType TokenType) (*Claims, error) {
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

func (u *UserSrv) RefreshToken(refreshToken string) (*schema.TokenPair, error) {
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

func (u *UserSrv) SetTokenBlacklist(token string, expiration time.Duration) error {
	ctx := context.Background()
	err := redisclient.GetClient().Set(ctx, "blacklist:"+token, 1, expiration).Err()
	if err != nil {
		return fmt.Errorf("设置token黑名单失败: %w", err)
	}
	return nil
}

func (u *UserSrv) IsTokenInBlacklist(token string) bool {
	ctx := context.Background()
	_, err := redisclient.GetClient().Get(ctx, "blacklist:"+token).Result()
	return err == nil
}
