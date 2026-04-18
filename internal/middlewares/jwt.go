package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	appcontext "tourism-backend/internal/context"
)

var (
	accessTokenSecret  []byte
	refreshTokenSecret []byte
)

const (
	AccessTokenTTL   = 15 * time.Minute
	RefreshTokenTTL  = 7 * 24 * time.Hour
	accessTokenType  = "access"
	refreshTokenType = "refresh"
)

var (
	errTokenSecretNotConfigured = errors.New("token secret is not configured")
	errInvalidToken             = errors.New("invalid token")
	errInvalidTokenType         = errors.New("invalid token type")
)

func InitSecret(secret string) {
	InitSecrets(secret, secret)
}

func InitSecrets(accessSecret, refreshSecret string) {
	accessTokenSecret = []byte(accessSecret)
	refreshTokenSecret = []byte(refreshSecret)
}

type RefreshToken struct {
	UserID int    `json:"user_id"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type AccessToken struct {
	UserID    int     `json:"user_id"`
	Email     string  `json:"email"`
	Name      string  `json:"name"`
	Role      string  `json:"role"`
	AvatarURL *string `json:"avatar_url,omitempty"`
	Balance   float64 `json:"balance"`
	Type      string  `json:"type"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	TokenType        string `json:"token_type"`
	AccessExpiresIn  int64  `json:"access_expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

func GenerateToken(userID int, email, role, name string, avatarURL *string, balance float64) (string, error) {
	return GenerateAccessToken(userID, email, role, name, avatarURL, balance)
}

func GenerateAccessToken(userID int, email, role, name string, avatarURL *string, balance float64) (string, error) {
	now := time.Now()
	claims := &AccessToken{
		UserID:    userID,
		Email:     email,
		Name:      name,
		Role:      role,
		AvatarURL: avatarURL,
		Balance:   balance,
		Type:      accessTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
		},
	}
	return signToken(claims, accessTokenSecret)
}

func GenerateRefreshToken(userID int) (string, error) {
	now := time.Now()
	claims := &RefreshToken{
		UserID: userID,
		Type:   refreshTokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.Itoa(userID),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenTTL)),
		},
	}
	return signToken(claims, refreshTokenSecret)
}

func GenerateTokenPair(userID int, email, role, name string, avatarURL *string, balance float64) (*TokenPair, error) {
	accessToken, err := GenerateAccessToken(userID, email, role, name, avatarURL, balance)
	if err != nil {
		return nil, err
	}
	refreshToken, err := GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		TokenType:        "Bearer",
		AccessExpiresIn:  int64(AccessTokenTTL.Seconds()),
		RefreshExpiresIn: int64(RefreshTokenTTL.Seconds()),
	}, nil
}

func ValidateToken(tokenStr string) (*AccessToken, error) {
	return ValidateAccessToken(tokenStr)
}

func ValidateAccessToken(tokenStr string) (*AccessToken, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&AccessToken{},
		keyFunc(accessTokenSecret),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AccessToken)
	if !ok || !token.Valid {
		return nil, errInvalidToken
	}
	if claims.Type != accessTokenType {
		return nil, errInvalidTokenType
	}

	return claims, nil
}

func ValidateRefreshToken(tokenStr string) (*RefreshToken, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&RefreshToken{},
		keyFunc(refreshTokenSecret),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*RefreshToken)
	if !ok || !token.Valid {
		return nil, errInvalidToken
	}
	if claims.Type != refreshTokenType {
		return nil, errInvalidTokenType
	}

	return claims, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := bearerTokenFromHeader(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusUnauthorized)
			return
		}

		claims, err := ValidateAccessToken(tokenStr)
		if err != nil {
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		appClaims := &appcontext.Claims{
			UserID: claims.UserID,
			Email:  claims.Email,
			Role:   claims.Role,
		}
		ctx := appcontext.WithClaims(r.Context(), appClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetClaims(r *http.Request) *appcontext.Claims {
	claims, _ := appcontext.GetClaims(r.Context())
	return claims
}

func RoleMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r)
			if claims == nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			for _, role := range roles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
		})
	}
}

func signToken(claims jwt.Claims, secret []byte) (string, error) {
	if len(secret) == 0 {
		return "", errTokenSecretNotConfigured
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func keyFunc(secret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if len(secret) == 0 {
			return nil, errTokenSecretNotConfigured
		}
		alg := "<nil>"
		if token.Method != nil {
			alg = token.Method.Alg()
		}
		if alg != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %s", alg)
		}
		return secret, nil
	}
}

func bearerTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("unauthorized")
	}

	parts := strings.Fields(authHeader)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid token format")
	}

	return parts[1], nil
}
