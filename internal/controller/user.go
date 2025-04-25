package controller

import (
	"crud/internal/model"
	"crud/internal/service"
	responseUtil "crud/internal/util/response"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
)

const DefaultOffset = 0
const DefaultLimit = 10

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{userService: userService}
}

func (controller *UserController) SetupRoutes(superRoute *gin.RouterGroup) {
	userRouter := superRoute.Group("user")
	{
		userRouter.Use(cors.Default())
		// Keeping it here for preflight requests because of https://github.com/gin-gonic/gin/issues/3546
		userRouter.OPTIONS("/*any", func(c *gin.Context) {})
		userRouter.GET("/", controller.GetUsers)
		userRouter.GET("/:id", controller.GetUserById)
		userRouter.POST("/", controller.CreateUser)
		userRouter.PUT("/:id", controller.UpdateUser)
		userRouter.DELETE("/:id", controller.DeleteUser)
	}
}

// GetUsers gets list of users
//
// @Summary		Gets list of users summary
// @Description	Gets list of users description
// @Produce		json
// @Success		200		{object}	model.UserResponse
// @Router		/user/ [get]
func (controller *UserController) GetUsers(context *gin.Context) {
	offset, err := responseUtil.GetIntParamOrDefault(context, "offset", DefaultOffset)
	if err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}

	limit, err := responseUtil.GetIntParamOrDefault(context, "limit", DefaultLimit)
	if err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}

	users, err := controller.userService.GetUsers(offset, limit)
	if err != nil {
		responseUtil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusOK, users)
}

// GetUserById gets user by id
//
// @Summary		Gets user by id summary
// @Description	Gets user by id description
// @Produce		json
// @Param		id 		path		int		true	"User ID"
// @Success		200		{object}	model.UserResponse
// @Failure		400		{object}	response.HTTPStatusMessage
// @Failure		404		{object}	response.HTTPStatusMessage
// @Router		/user/{id} [get]
func (controller *UserController) GetUserById(context *gin.Context) {
	id, err := responseUtil.GetIntParam(context, "id")
	if err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userService.GetById(id)
	if err != nil {
		responseUtil.NewError(context, http.StatusNotFound, err)
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
// @Param		user	body		model.CreateUserRequest	true	"Add user"
// @Success		201		{object}	model.UserResponse
// @Failure		400		{object}	response.HTTPStatusMessage
// @Router		/user/ [post]
func (controller *UserController) CreateUser(context *gin.Context) {
	request := model.CreateUserRequest{}
	if err := context.ShouldBindJSON(&request); err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}

	userResponse, err := controller.userService.Create(&request)
	if err != nil {
		responseUtil.NewError(context, http.StatusInternalServerError, err)
		return
	}
	context.JSON(http.StatusCreated, userResponse)
}

// UpdateUser updates a user in the user service
//
// @Summary		Updates a user
// @Description	Updates a user
// @Accept		json
// @Produce		json
// @Param		user	body		model.UpdateUserRequest	true	"User new data"
// @Success		200		{object}	model.UserResponse
// @Failure		400		{object}	response.HTTPStatusMessage
// @Failure		404		{object}	response.HTTPStatusMessage
// @Router		/user/{id} [put]
func (controller *UserController) UpdateUser(context *gin.Context) {
	updateUserRequest := model.UpdateUserRequest{}
	if err := context.ShouldBindJSON(&updateUserRequest); err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userService.Update(&updateUserRequest)
	if err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, user)
}

// DeleteUser updates a user in the user service
//
// @Summary		Deletes a user
// @Description	Deletes a user
// @Accept		json
// @Produce		json
// @Param		id		path		int			true	"User ID"
// @Success		200		{object}	response.HTTPStatusMessage
// @Failure		400		{object}	response.HTTPStatusMessage
// @Router		/user/{id} [delete]
func (controller *UserController) DeleteUser(context *gin.Context) {
	id, err := responseUtil.GetIntParam(context, "id")
	if err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}

	user, err := controller.userService.Delete(id)
	if err != nil {
		responseUtil.NewError(context, http.StatusBadRequest, err)
		return
	}
	context.JSON(http.StatusOK, user)
}
