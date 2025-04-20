package controller

import (
	"crud/httputil"
	"crud/model"
	"errors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var counter = 0
var db = make(map[int]*model.User)

func UserRoutes(superRoute *gin.RouterGroup) {
	userRouter := superRoute.Group("user")
	{
		userRouter.Use(cors.Default())
		// Keeping it here for preflight requests because of https://github.com/gin-gonic/gin/issues/3546
		userRouter.OPTIONS("/*any", func(c *gin.Context) {})
		userRouter.GET("/", GetUsers)
		userRouter.GET("/:id", GetUserById)
		userRouter.POST("/", CreateUser)
		userRouter.PUT("/:id", UpdateUser)
		userRouter.DELETE("/:id", DeleteUser)
	}
}

// GetUsers gets list of users
//
// @Summary		Gets list of users summary
// @Description	Gets list of users description
// @Produce		json
// @Success		200		{object}	model.User
// @Router		/user/ [get]
func GetUsers(context *gin.Context) {
	users := make([]model.User, 0, len(db))
	for _, value := range db {
		users = append(users, *value)
	}
	context.JSON(http.StatusOK, users)
}

// GetUserById gets user by id
//
// @Summary		Gets user by id summary
// @Description	Gets user by id description
// @Produce		json
// @Param		id 		path		int		true	"User ID"
// @Success		200		{object}	model.User
// @Failure		400		{object}	httputil.HTTPStatusMessage
// @Failure		404		{object}	httputil.HTTPStatusMessage
// @Router		/user/{id} [get]
func GetUserById(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	user, found := db[id]
	if !found {
		httputil.NewError(context, http.StatusNotFound, errors.New("user not found"))
		return
	}
	context.JSON(http.StatusOK, user)
}

// CreateUser gets user by id
//
// @Summary		Updates a user summary
// @Description	Updates a user description
// @Accept		json
// @Produce		json
// @Param		user	body		model.CreateOrUpdateUser	true	"Add user"
// @Success		201		{object}	model.User
// @Failure		400		{object}	httputil.HTTPStatusMessage
// @Router		/user/ [post]
func CreateUser(context *gin.Context) {
	var createUser model.CreateOrUpdateUser
	if err := context.ShouldBind(&createUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	user := model.User{ID: counter, Name: createUser.Name, Email: createUser.Email, Age: createUser.Age}
	counter++
	db[user.ID] = &user
	context.JSON(http.StatusCreated, user)
}

// UpdateUser updates a user in the user service
//
// @Summary		Updates a user
// @Description	Updates a user
// @Accept		json
// @Produce		json
// @Param		id		path		int			true	"User ID"
// @Param		user	body		model.CreateOrUpdateUser	true	"User new data"
// @Success		200		{object}	model.User
// @Failure		400		{object}	httputil.HTTPStatusMessage
// @Failure		404		{object}	httputil.HTTPStatusMessage
// @Router		/user/{id} [put]
func UpdateUser(context *gin.Context) {
	var updateUser model.CreateOrUpdateUser
	if err := context.ShouldBind(&updateUser); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	user, found := db[id]
	if !found {
		httputil.NewError(context, http.StatusNotFound, errors.New("user not found"))
		return
	}
	user.Name = updateUser.Name
	user.Email = updateUser.Email
	user.Age = updateUser.Age
	context.JSON(http.StatusOK, user)
}

// DeleteUser updates a user in the user service
//
// @Summary		Deletes a user
// @Description	Deletes a user
// @Accept		json
// @Produce		json
// @Param		id		path		int			true	"User ID"
// @Success		200		{object}	httputil.HTTPStatusMessage
// @Failure		400		{object}	httputil.HTTPStatusMessage
// @Router		/user/{id} [delete]
func DeleteUser(context *gin.Context) {
	idStr := context.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		httputil.NewError(context, http.StatusBadRequest, err)
		return
	}
	delete(db, id)
	context.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
