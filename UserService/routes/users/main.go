package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Pmca96/API-Servers/UserService/connectors"
	"github.com/Pmca96/API-Servers/UserService/models/users"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/swag/example/celler/httputil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

var UsersCollection *mongo.Collection = connectors.MongoOpenCollection(connectors.MongoClient, "users")

// @Summary      Add User
// @Description  Add a user
// @Tags         Users
// @Param        user body users.UserWithoutId true "user"
// @Success      200  {object}  users.Users
// @Failure    	 500  {object}  object
// @Router       /users/AddUser [post]
func AddUser(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user users.Users

	if err := c.BindJSON(&user); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		httputil.NewError(c, http.StatusInternalServerError, validationErr)
		fmt.Println(validationErr)
		return
	}
	user.ID = primitive.NewObjectID()

	result, insertErr := UsersCollection.InsertOne(ctx, user)
	if insertErr != nil {
		httputil.NewError(c, http.StatusInternalServerError, insertErr)
		fmt.Println(insertErr)
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

// @Summary      Get Users
// @Description  Get AllUsers
// @Tags         Users
// @Success      200  {object}  users.Users
// @Failure    	 500  {object}  object
// @Router       /users/GetUsers [get]
func GetUsers(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var users []bson.M

	cursor, err := UsersCollection.Find(ctx, bson.M{})

	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &users); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(users)

	c.JSON(http.StatusOK, users)
}

// @Summary      Get User By Email
// @Description  Get User By Email
// @Tags         Users
// @Param        email query string true "Email"
// @Success      200  {object}  users.Users
// @Failure    	 500  {object}  gin.H
// @Router       /users/GetUserByEmail [get]
func GetUserByEmail(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var email users.UserEmail
	var users []bson.M

	if err := c.BindQuery(&email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	cursor, err := UsersCollection.Find(ctx, bson.M{"email": email})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(users)

	c.JSON(http.StatusOK, users)
}

// @Summary      Get User By Id
// @Description  Get User By Id
// @Tags         Users
// @Param        id path string true "Id"
// @Success      200  {object}  users.Users
// @Failure    	 500  {object}  object
// @Router       /users/GetUserById [get]
func GetUserById(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user bson.M

	if err := UsersCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&user); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(user)

	c.JSON(http.StatusOK, user)
}

// @Summary      Update Name
// @Description  Update Name
// @Tags         Users
// @Param        id path string true "Id"
// @Param        firstname body users.UserName true "FirstName"
// @Success      200  {object}  int
// @Failure    	 500  {object}  object
// @Router       /users/UpdateName [put]
func UpdateName(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var userName users.UserName

	if err := c.BindJSON(&userName); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	result, err := UsersCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{{"server", userName.FirstName}}},
		},
	)

	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.ModifiedCount)

}

// @Summary      Update User First Name
// @Description  Update User First Name
// @Tags         Users
// @Param        id path string true "Id"
// @Param        firstname body users.UserName true "FirstName"
// @Success      200  {object}  int
// @Failure    	 400  {object}  object
// @Router       /users/UpdateUser [put]
func UpdateUser(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user users.UserWithoutId

	if err := c.BindJSON(&user); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		httputil.NewError(c, http.StatusInternalServerError, validationErr)
		fmt.Println(validationErr)
		return
	}

	result, err := UsersCollection.ReplaceOne(
		ctx,
		bson.M{"_id": docID},
		bson.M{
			"firstname": user.FirstName,
			"lastname":  user.LastName,
			"age":       user.Age,
			"city":      user.City,
		},
	)

	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.ModifiedCount)
}

// @Summary      Delete User
// @Description  Delete User
// @Tags         Users
// @Param        id path string true "Id"
// @Success      200  {object}  int
// @Failure    	 400  {object}  object
// @Router       /users/DeleteUser [delete]
func DeleteUser(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := UsersCollection.DeleteOne(ctx, bson.M{"_id": docID})

	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.DeletedCount)
}
