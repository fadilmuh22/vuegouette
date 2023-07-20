package util

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"

	"github.com/fadilmuh22/restskuy/internal/model"
)

// 30 days
var LOGIN_EXPIRATION_DURATION = time.Duration(24 * 30 * time.Hour)

func GetJWTSecret() string {
	return viper.GetString("JWT_SECRET")
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time.
type Claims struct {
	ID      uuid.UUID `json:"uuid"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	IsAdmin bool      `json:"is_admin"`
	jwt.RegisteredClaims
}

// GenerateTokens generates jwt token and saves it to the http-only cookie.
func GenerateAccessToken(user *model.User, c echo.Context) (string, error) {
	expirationTime := time.Now().Add(LOGIN_EXPIRATION_DURATION)

	accessToken, _, err := generateToken(user, expirationTime, []byte(GetJWTSecret()))

	if err != nil {
		return accessToken, err
	}

	return accessToken, nil
}

// Pay attention to this function. It holds the main JWT token generation logic.
func generateToken(user *model.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &Claims{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(err error, c echo.Context) error {
	// Redirects to the signIn form.
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
