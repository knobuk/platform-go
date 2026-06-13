package platform

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID string   `json:"uid"`
	Email  string   `json:"email"`
	Roles  []string `json:"roles"`
	Type   string   `json:"typ"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

func NewJWTManager(accessSecret, refreshSecret string, accessMinutes, refreshDays int) *JWTManager {
	return &JWTManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     time.Duration(accessMinutes) * time.Minute,
		refreshTTL:    time.Duration(refreshDays) * 24 * time.Hour,
	}
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func (m *JWTManager) Generate(userID uuid.UUID, email string, roles []string) (TokenPair, error) {
	access, err := m.sign(userID, email, roles, "access", m.accessTTL, m.accessSecret)
	if err != nil {
		return TokenPair{}, err
	}
	refresh, err := m.sign(userID, email, roles, "refresh", m.refreshTTL, m.refreshSecret)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{AccessToken: access, RefreshToken: refresh, ExpiresIn: int64(m.accessTTL.Seconds())}, nil
}

func (m *JWTManager) sign(userID uuid.UUID, email string, roles []string, typ string, ttl time.Duration, secret []byte) (string, error) {
	claims := Claims{
		UserID: userID.String(),
		Email:  email,
		Roles:  roles,
		Type:   typ,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func (m *JWTManager) ParseAccess(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, "access", m.accessSecret)
}

func (m *JWTManager) ParseRefresh(tokenStr string) (*Claims, error) {
	return m.parse(tokenStr, "refresh", m.refreshSecret)
}

func (m *JWTManager) parse(tokenStr, expectedType string, secret []byte) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}
	if claims.Type != expectedType {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}
