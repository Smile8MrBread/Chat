package jwt

import (
	"github.com/Smile8MrBread/Chat/auth_service/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, duration time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["login"] = user.Login
	claims["user_id"] = user.Id
	claims["avatar"] = user.Avatar
	claims["first_ame"] = user.FirstName
	claims["last_name"] = user.LastName
	claims["exp"] = time.Now().Add(duration).Add(time.Hour * 3).Unix()

	tokenString, err := token.SignedString([]byte(models.ReturnSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
