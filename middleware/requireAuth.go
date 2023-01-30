package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mdokusV/go-Pracownia/initializers"
	"github.com/mdokusV/go-Pracownia/models"
)

func CheckAuthCookie(c *fiber.Ctx) error {
	//Get the cookie of req
	jwtCookie := c.Request().Header.Cookie("Authorization")
	if jwtCookie == nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	//Validate cookie
	token, err := jwt.Parse(string(jwtCookie), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error: Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//check the exp
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		//Find the user with token sub
		var user models.User
		IdUser := claims["sub"]
		err = initializers.DB.First(&user, IdUser).Error

		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		//attach to req
		c.Locals("userCookie", user)

		//continue
		return c.Next()
	} else {
		fmt.Println(err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}
}
