package routes

import "github.com/gin-gonic/gin"

func AddRoutes(service *gin.Engine) {

	UsersRoutes(service)

}
