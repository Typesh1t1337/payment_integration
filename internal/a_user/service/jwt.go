package service

import (
	"crypto/rsa"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type Token struct {
	Sub  string `json:"sub"`
	Exp  int64  `json:"exp"`
	Type string `json:"type"`
}

type JwtService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJwtService(privateKeyPath string, publicKeyPath string, accessTTL time.Duration, refreshTTL time.Duration) (*JwtService, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}
	return &JwtService{privateKey: privateKey, publicKey: publicKey, accessTTL: accessTTL, refreshTTL: refreshTTL}, nil
}

func (j *JwtService) GenerateAccessToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.accessTTL).Unix(),
		"type": "access",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j *JwtService) GenerateRefreshToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.refreshTTL).Unix(),
		"type": "refresh",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j *JwtService) ParseToken(tokenString string) (Token, error) {
	token ,err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return j.publicKey, nil
	})
	if err != nil {
		return Token{}, err
	}
	return Token{
		Sub: token.Claims.(jwt.MapClaims)["sub"].(string),
		Exp: int64(token.Claims.(jwt.MapClaims)["exp"].(float64)),
		Type: token.Claims.(jwt.MapClaims)["type"].(string),
	}, nil
}