package routes

import (
	"UserServer/routes/users"

	"github.com/gin-gonic/gin"
)

func UsersRoutes(router *gin.Engine) {
	//C
	router.POST("/users/AddUsers", users.AddUsers)
	//R
	router.GET("/users/GetUsersByEmail", users.GetUsersByEmail)
	router.GET("/users/GetUsers", users.GetUsers)
	router.GET("/users/GetUsersById/:id", users.GetUsersById)
	//U
	router.PUT("/users/UpdateUser/:id", users.UpdateUser)
	router.PUT("/users/UpdateName/:id", users.UpdateName)
	//D
	router.DELETE("/users/DeleteUser/:id", users.DeleteUser)
}
