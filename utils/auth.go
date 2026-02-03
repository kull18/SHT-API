package utils

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    Rol    string `json:"rol"`
    jwt.RegisteredClaims
}

// HashPassword genera un hash de la contraseña
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// CheckPasswordHash verifica si la contraseña coincide con el hash
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// GenerateJWT genera un token JWT
func GenerateJWT(userID int, email, rol string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Rol:    rol,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// ValidateJWT valida un token JWT
func ValidateJWT(tokenString string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })

    if err != nil {
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("token inválido")
    }

    return claims, nil
}
