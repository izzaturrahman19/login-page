package controller

import (
	"context"
	"net/http"
	"log"
	"fmt"

	"github.com/Izzaturrahman19/login-page/config/db"
	"github.com/Izzaturrahman19/login-page/model"
	"github.com/Izzaturrahman19/login-page/apierror"
    "github.com/labstack/echo"
	
	jwt "github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(c echo.Context) error {

	user := new(model.User)
	err := c.Bind(user)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, "Failed to get signin data. Probably content-type is not match with actual body type", err)
	}

	collection, err := db.GetDbCollection()

	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, "Failed to get database", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	user.Password = string(hash)
	_, err = collection.Collection("users").InsertOne(context.Background(), user)
	   if err != nil {
			return apierror.NewError(http.StatusUnprocessableEntity, "Failed to insert to database", err)
		}

	return c.JSON(http.StatusCreated, user)
}

func LoginHandler(c echo.Context) error {

	user := new(model.User)
 	err := c.Bind(user)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, "Failed to get signin data. Probably content-type is not match with actual body type", err)
	}

	collection, err := db.GetDbCollection()

	if err != nil {
		log.Fatal(err)
	}
	var result model.User
	var bearer model.Token
	
	err = collection.Collection("users").FindOne(context.Background(), bson.D{{"username", user.Username}}).Decode(&result)

	fmt.Println("result : ", result)
	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, "Collection not found", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(user.Password))

	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, "Password missmatch", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username":  result.Username,
		"firstname": result.FirstName,
		"lastname":  result.LastName,
	})

	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return apierror.NewError(http.StatusUnprocessableEntity, "Create token failed", err)
	}

	bearer.Token = tokenString

	return c.JSON(http.StatusCreated, bearer)
}

func ProfileHandler(c echo.Context) error {

	tokenString := c.Request().Header.Get("Authorization")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return []byte("secret"), nil
	})
	var result model.User
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		result.Username = claims["username"].(string)
		result.FirstName = claims["firstname"].(string)
		result.LastName = claims["lastname"].(string)

		return c.JSON(http.StatusOK, result)
	} else {
		return apierror.NewError(http.StatusUnprocessableEntity, "Failed to get data", err)
	}

}