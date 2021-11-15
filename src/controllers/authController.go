package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/sdblg/hoa-auth/src/dao"
	"github.com/sdblg/hoa-auth/src/lib"
	"github.com/sdblg/hoa-auth/src/models"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

// UserRegister users
func UserRegister(ctx *fiber.Ctx) error {

	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(fiber.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"data":    nil,
			"message": "request is invalid",
			"track":   err.Error(),
		})
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest)

		return ctx.JSON(fiber.Map{
			"data":    nil,
			"message": "request is invalid",
			"track":   err.Error(),
		})
	}
	u := &models.User{
		Model:     gorm.Model{},
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
		Role:      "user",
	}

	//var u *models.User
	//if err := json.Unmarshal(ctx.Body(), &u); err != nil {
	//	ctx.Status(fiber.StatusBadRequest)
	//	return ctx.JSON(fiber.Map{
	//		"data":    nil,
	//		"message": "request is invalid",
	//		"track":   err.Error(),
	//	})
	//}

	if lib.IsEmailValid(u.Email) == false {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"data":    nil,
			"message": "email is invalid format",
		})
	}

	var dbUser *models.User
	dao.DB.Where("email = ?", u.Email).Find(&dbUser)
	if dbUser.ID != 0 {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"data":    nil,
			"message": "email is already exist",
		})
	}

	dao.DB.Debug().Create(u)
	//if err != nil {
	//	ctx.Status(fiber.StatusInternalServerError)
	//	log.Errorln("Create=", err)
	//	return ctx.JSON(
	//		map[string]interface{}{
	//			"data":    nil,
	//			"message": err,
	//		},
	//	)
	//}

	return ctx.JSON(u)
}

func HealthCheck(ctx *fiber.Ctx) error {
	ctx.Status(fiber.StatusOK)
	return ctx.SendString("app is up and run")
}

func UserLogin(ctx *fiber.Ctx) error {
	var data map[string]string

	if err := ctx.BodyParser(&data); err != nil {
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"data":    nil,
			"message": err.Error(),
		})
	}
	var user *models.User
	dao.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		ctx.Status(fiber.StatusNotFound)
		return ctx.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		log.Errorln(err)
		ctx.Status(fiber.StatusBadRequest)
		return ctx.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}
	cookieDurationStr := os.Getenv("COOKIE_EXPIRE_DURATION")
	cookieDuration, _ := time.ParseDuration(cookieDurationStr)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(cookieDuration).Unix(), //1 hour
	})

	token, err := claims.SignedString([]byte(os.Getenv("TOKEN_SHARED_KEY")))

	if err != nil {
		ctx.Status(fiber.StatusInternalServerError)
		return ctx.JSON(fiber.Map{
			"message": "could not login",
		})
	}

	cookie := fiber.Cookie{
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    token,
		Expires:  time.Now().Add(cookieDuration),
		HTTPOnly: true,
	}

	ctx.Cookie(&cookie)

	return ctx.JSON(fiber.Map{
		"message": "success",
	})

}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies(os.Getenv("COOKIE_NAME"))

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SHARED_KEY")), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user models.User

	dao.DB.Where("id = ?", claims.Issuer).First(&user)

	return c.JSON(user)
}

func UserLogout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
