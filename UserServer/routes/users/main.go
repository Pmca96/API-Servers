package users

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"UserServer/connectors"
	"UserServer/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/swaggo/swag/example/celler/httputil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var validate = validator.New()

var UsersCollection *mongo.Collection = connectors.MongoOpenCollection(connectors.MongoClient, "users")

func AddUsers(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.Users

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		fmt.Println(validationErr)
		return
	}
	user.ID = primitive.NewObjectID()

	result, insertErr := UsersCollection.InsertOne(ctx, user)
	if insertErr != nil {
		msg := fmt.Sprintf("order item was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println(insertErr)
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

// GetUsers godoc
// @Summary      Get Users
// @Description  Get AllUsers
// @Tags         Users
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Users
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

//get all users by the email
func GetUsersByEmail(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var users []bson.M

	var email string

	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

//get an users by its id
func GetUsersById(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user bson.M

	if err := UsersCollection.FindOne(ctx, bson.M{"_id": docID}).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	fmt.Println(user)

	c.JSON(http.StatusOK, user)
}

//update a waiter's name for an order
func UpdateName(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	type UserName struct {
		FirstName *string `json:"firstname"`
	}

	var userName UserName

	if err := c.BindJSON(&userName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	result, err := UsersCollection.UpdateOne(ctx, bson.M{"_id": docID},
		bson.D{
			{"$set", bson.D{{"server", userName.FirstName}}},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.ModifiedCount)

}

//update the user
func UpdateUser(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var user models.Users

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(user)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.ModifiedCount)
}

//delete an user given the id
func DeleteUser(c *gin.Context) {

	userID := c.Params.ByName("id")
	docID, _ := primitive.ObjectIDFromHex(userID)

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	result, err := UsersCollection.DeleteOne(ctx, bson.M{"_id": docID})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	defer cancel()

	c.JSON(http.StatusOK, result.DeletedCount)

}
