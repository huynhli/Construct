package handlers

import (
	"backend/database"
	"context"
	"errors"
	"log"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// signup
func Signup(c *fiber.Ctx) error {
	username := c.Query("username")

	err := usernameValidityChecker(username)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	err = passwordValidityChecker(username)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	// creds good, check if user exists
	var res bson.M
	err = database.ConstructDatabase.Collection("users").FindOne(context.TODO(), bson.M{"name": username}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		// user does not exist --> add to db
		// company is guaranteed to exist bc of dropdown
		insertUserIntoDB(username, c.Query("company"))
	}
	// user exists --> return msg user already exists

	return c.SendString("woohoo")
}

// returns names of all companies in db
func getAllCompanies() []string {
	var res bson.M
	cursor, err := database.ConstructDatabase.Collection("companies").Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"name": 1, "_id": 0})).Decode(&res)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.TODO())

	var companyNames []string
	for cursor.Next(context.TODO()) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			log.Fatal(err)
		}
		companyNames = append(companyNames, result["name"].(string))
	}
	return companyNames
}

func getCompanyIDWithCompanyName(companyName string) error {
	var res bson.M
	err := database.ConstructDatabase.Collection("companies").FindOne(context.TODO(), bson.M{"name": companyName}).Decode(&res)
	if err == mongo.ErrNoDocuments {
		// company does not exist
	}
	return nil
}

func insertCompanyIntoDB(companyName string) error {
	companyObj, err := database.ConstructDatabase.Collection("companies").
		InsertOne(context.TODO(), bson.M{
			"name": companyName,
		})
	if err != nil {
		return err
	}
	return nil
}

func insertUserIntoDB(username string, companyName string, hashedpw string) error {
	userObj, err := database.ConstructDatabase.Collection("users").
		InsertOne(context.TODO(), bson.M{
			"name":      username,
			"companyId": getCompanyIDWithCompanyName(companyName),
			"password":  hashedpw, "projects": make([]string),
			"tasks": make([]string),
		})
	if err != nil {
		return err
	}
	return nil
}

// login
func Login(c *fiber.Ctx) error {
	username := c.Query("username")
	err := usernameValidityChecker(username)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Username cannot be empty",
		}
	}

	return c.SendString("woohoo")
}

// for a given admin, show projects they are related to
func AdminGetProjects(c *fiber.Ctx) error {

	return c.SendString("P")
}

// for a given admin, show tasks they can use
func AdminGetTasks(c *fiber.Ctx) error {

	return c.SendString("P")
}

// for a given user, show projects they are related to
func UserGetProjects(c *fiber.Ctx) error {

	return c.SendString("P")
}

// for a given user, show tasks they are assigned
func UserGetTasks(c *fiber.Ctx) error {
	// tasksCollection := database.ConstructDatabase.Collection("")
	return c.SendString("Updated!")
}

func usernameValidityChecker(username string) error {
	if username == "" {
		return errors.New("username cannot be empty")
	} else if len(username) > 25 || len(username) < 3 {
		return errors.New("username not within length requirement")
	} else if re := regexp.MustCompile(`^[a-zA-Z0-9]+$`); !re.MatchString(username) {
		return errors.New("username not aligning with guidelines")
	}
	return nil
}

func passwordValidityChecker(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	} else if len(password) > 25 || len(password) < 3 {
		return errors.New("password not within length requirement")
	} else if re := regexp.MustCompile(`^[a-zA-Z0-9]+$`); !re.MatchString(password) {
		return errors.New("password not aligning with guidelines")
	}
	return nil
}
