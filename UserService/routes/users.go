package routes

import (
	"github.com/Pmca96/API-Servers/UserService/routes/users"

	"github.com/gin-gonic/gin"
)

const controllerName = "/User/"

func UsersRoutes(router *gin.Engine) {
	//C
	router.POST(controllerName+"AddUser", users.AddUser)
	//R
	router.GET(controllerName+"GetUserByEmail", users.GetUserByEmail)
	router.GET(controllerName+"GetUsers", users.GetUsers)
	router.GET(controllerName+"GetUserById/:id", users.GetUserById)
	//U
	router.PUT(controllerName+"UpdateUser/:id", users.UpdateUser)
	router.PUT(controllerName+"UpdateName/:id", users.UpdateName)
	//D
	router.DELETE(controllerName+"DeleteUser/:id", users.DeleteUser)
}
