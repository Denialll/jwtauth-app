package pkg

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// TokenManager provides logic for JWT & Refresh tokens generation and parsing.
type TokenManager interface {
	NewJWT(uuid string) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken(uuid2 string) (string, error)
}

type Manager struct {
	signingKey      string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewManager(signingKey string, accessTokenTTL, refreshTokenTTL time.Duration) (*Manager, error) {
	if signingKey == "" {
		return nil, errors.New("empty signing key")
	}
	if accessTokenTTL >= refreshTokenTTL {
		return nil, errors.New("AccessTokenTTL more then RefreshTokenTTL")
	}

	return &Manager{
		signingKey:      signingKey,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

func (m *Manager) NewJWT(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.accessTokenTTL).Unix(),
			//IssuedAt:  time.Now().Unix(),
			Subject: userId,
		})
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKey), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("error get user claims from token")
	}

	return claims["sub"].(string), nil
}

func (m *Manager) NewRefreshToken(userId string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.refreshTokenTTL).Unix(),
			//IssuedAt:  time.Now().Unix(),
			Subject: userId,
		})
	token.Header["typ"] = "REF"

	return token.SignedString([]byte(m.signingKey))
}

//func (m *Manager) NewRefreshToken() (models.Session, error) {
//	b := make([]byte, 32)
//
//	s := rand.NewSource(time.Now().Unix())
//	r := rand.New(s)
//
//	if _, err := r.Read(b); err != nil {
//		return models.Session{}, err
//	}
//
//	session := models.Session{
//		RefreshToken: fmt.Sprintf("%x", b),
//		ExpiresAt:    time.Now().Add(m.refreshTokenTTL),
//	}
//
//	return session, nil
//}
