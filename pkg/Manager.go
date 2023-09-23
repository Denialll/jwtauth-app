package pkg

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	time "time"
)

type TokenManager interface {
	NewAccessToken(uuid string, timeCreate time.Time) (string, error)
	Parse(accessToken string) (jwt.MapClaims, error)
	ParseUnverified(accessToken string) (jwt.MapClaims, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	signingKeyAccess  string
	signingKeyRefresh string
	accessTokenTTL    time.Duration
	refreshTokenTTL   time.Duration
}

func NewManager(signingKeyAccess string, accessTokenTTL, refreshTokenTTL time.Duration) (*Manager, error) {
	if signingKeyAccess == "" {
		return nil, errors.New("empty signing key")
	}
	if accessTokenTTL >= refreshTokenTTL {
		return nil, errors.New("AccessTokenTTL more then RefreshTokenTTL")
	}

	return &Manager{
		signingKeyAccess: signingKeyAccess,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
	}, nil
}

func (m *Manager) NewAccessToken(userId string, timeCreate time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.accessTokenTTL).Unix(),
			NotBefore: timeCreate.Unix(),
			Subject:   userId,
		})

	return token.SignedString([]byte(m.signingKeyAccess))
}

func (m *Manager) Parse(accessToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.signingKeyAccess), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}

	fmt.Println(claims)

	return claims, nil
}

func (m *Manager) ParseUnverified(accessToken string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(accessToken, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error get user claims from token")
	}

	fmt.Println(claims)

	return claims, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
