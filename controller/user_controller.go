package controller

import (
	"fmt"
	"log"
	"net/http"
	"simple_crude/manager"
	"simple_crude/request"
	"strings"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type UserCont struct {
	manager *manager.UserMgr
	
}

func NewUserController(mn *manager.UserMgr) *UserCont {
	return &UserCont{manager: mn}
}
func (us *UserCont) CreateUser(c echo.Context) error {

	var req request.StdRequest
	log.Println("rq------>")
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	// log.Println("REQ-------->", req)
	if err := validate.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	createdUser, err := us.manager.CreateUser(req)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, createdUser)
}

func (us *UserCont) UpdateUser(c echo.Context) error {
	id := c.Param("id")

	var req request.StdRequest
	if err := c.Bind(&req); err != nil {
		log.Printf("Failed to bind request: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Invalid ID format: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	updatedUser, err := us.manager.UpdateUser(objectId, req)
	if err != nil {
		log.Printf("Error updating user with ID %s: %v", id, err)
		if strings.Contains(err.Error(), "no user found with the given ID") {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": fmt.Sprintf("No user found with the given ID: %s. Please check the ID and try again.", id),
			})
		}
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "ID is Not Found.",
		})
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func (us *UserCont) Getall(c echo.Context) error {
	var req request.StdRequest
	req.Name = c.QueryParam("name")
	req.Student = c.QueryParam("student") == "true"
	req.Tag = c.QueryParams()["tag"]
	userResponse, err := us.manager.Getall(req)
	if err != nil {
		return c.JSON(http.StatusFound, map[string]string{"error": err.Error()})

	}
	return c.JSON(http.StatusOK,userResponse )

}

func (us *UserCont) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	objectId, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		log.Printf("Invalid ObjectID format: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	if err := us.manager.DeleteUser(objectId); err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User successfully deleted"})

}
