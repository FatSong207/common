package jwt

import (
	"github.com/golang-jwt/jwt/v4"
)

// GetJwtToken 獲取Token secretKey為加密字串. key為唯一標示
func GetJwtToken(secretKey string, iat, seconds int64, key string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["key"] = key
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

// ParseToken 解析 JWT Token 並回傳 Key
func ParseToken(tokenString, secretKey string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// jwt.MapClaims 會把 number 解析成 float64
		if k, ok := claims["key"].(string); ok {
			return k, nil
		}
	}

	return "", jwt.ErrInvalidKey
}
