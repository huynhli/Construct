package handlers

import (
	"backend/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type SignupLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// signup, takes in POST req
func Signup(c *fiber.Ctx) error {
	var signInReq SignupLoginRequest
	err := c.BodyParser(&signInReq)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Invalid request body",
		}
	}

	// double check
	err = usernameValidityChecker(signInReq.Username)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	// query
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var resUsername string
	err = database.DB.QueryRowContext(ctx, "SELECT username FROM users WHERE username = $1", signInReq.Username).Scan(&resUsername)
	if err == sql.ErrNoRows {
		// no rows --> no exist, so make user
		err = passwordValidityChecker(signInReq.Password)
		if err != nil {
			return &fiber.Error{
				Code:    fiber.ErrBadRequest.Code,
				Message: err.Error(),
			}
		}

		// bcrypt
		var hashedBytePW []byte
		hashedBytePW, err = bcrypt.GenerateFromPassword([]byte(signInReq.Password), bcrypt.DefaultCost)
		if err != nil {
			return &fiber.Error{
				Code:    fiber.ErrBadGateway.Code,
				Message: err.Error(),
			}
		}

		// make user in table
		_, err := database.DB.ExecContext(ctx, "INSERT INTO users (username, password_hash) VALUES ($1, $2)", signInReq.Username, string(hashedBytePW))
		if err != nil {
			return &fiber.Error{
				Code:    fiber.ErrBadGateway.Code,
				Message: err.Error(),
			}
		}
		return c.JSON(fiber.Map{
			"Code": 201,
		})

	} else if err != nil {
		// other errors --> error querying
		return &fiber.Error{
			Code:    fiber.ErrBadGateway.Code,
			Message: err.Error(),
		}
	}
	// row match! --> exist, so return that user already exists
	return c.Status(fiber.StatusConflict).JSON(fiber.Map{
		"message": "username already exists",
	})

	// ensure table has only one row --> enforced by schema anyways

	// defer rows.Close() 9=gh
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	err := rows.Scan(&id, &name)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	log.Println(id, name)
	// }

	// var res bson.M
	// err = database.ConstructDatabase.Collection("users").FindOne(context.TODO(), bson.M{"name": username}).Decode(&res)
	// if err == mongo.ErrNoDocuments {
	// 	// user does not exist --> add to db
	// 	// company is guaranteed to exist bc of dropdown
	// 	insertUserIntoDB(username, c.Query("company"))
	// }
	//  // user exists --> return msg user already exists
}

// // returns names of all companies in db
// func getAllCompanies() []string {
// 	var res bson.M
// 	cursor, err := database.ConstructDatabase.Collection("companies").Find(context.TODO(), bson.M{}, options.Find().SetProjection(bson.M{"name": 1, "_id": 0})).Decode(&res)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer cursor.Close(context.TODO())

// 	var companyNames []string
// 	for cursor.Next(context.TODO()) {
// 		var result bson.M
// 		if err := cursor.Decode(&result); err != nil {
// 			log.Fatal(err)
// 		}
// 		companyNames = append(companyNames, result["name"].(string))
// 	}
// 	return companyNames
// }
//
// func getCompanyIDWithCompanyName(companyName string) error {
// 	var res bson.M
// 	err := database.ConstructDatabase.Collection("companies").FindOne(context.TODO(), bson.M{"name": companyName}).Decode(&res)
// 	if err == mongo.ErrNoDocuments {
// 		// company does not exist
// 	}
// 	return nil
// }
//
// func insertCompanyIntoDB(companyName string) error {
// 	companyObj, err := database.ConstructDatabase.Collection("companies").
// 		InsertOne(context.TODO(), bson.M{
// 			"name": companyName,
// 		})
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// submitted login data
func Login(c *fiber.Ctx) error {
	var loginReq SignupLoginRequest
	err := c.BodyParser(&loginReq)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: "Invalid request body",
		}

	}

	err = usernameValidityChecker(loginReq.Username)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	err = passwordValidityChecker(loginReq.Password)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: err.Error(),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var resUserID int
	var isAdmin bool
	var storedPWHash string
	err = database.DB.QueryRowContext(ctx, `SELECT id, password_hash, is_admin FROM users WHERE username = $1`, loginReq.Username).Scan(&resUserID, &storedPWHash, &isAdmin)
	if err == sql.ErrNoRows {
		return fiber.NewError(fiber.StatusUnauthorized, "invalid username or password")
	} else if err != nil {
		// other errors --> error querying
		return &fiber.Error{
			Code:    fiber.ErrBadGateway.Code,
			Message: err.Error(),
		}
	}

	// bcrypt check
	err = bcrypt.CompareHashAndPassword([]byte(storedPWHash), []byte(loginReq.Password))
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrUnauthorized.Code,
			Message: "Incorrect password",
		}
	}

	// row match! --> exist, so return generate jwt with username and postgres obj id
	jwtString, err := GenerateJWT(loginReq.Username, resUserID, isAdmin)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.ErrBadGateway.Code,
			Message: err.Error(),
		}
	}
	// send jwt to frontend as json
	return c.JSON(fiber.Map{
		"token": jwtString,
	})

}

func GenerateJWT(username string, userID int, isAdmin bool) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userID,
		"username": username,
		"isAdmin":  isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("could not generate JWT")
	}

	return tokenString, nil
}

type Project struct {
	ID          int        `json:"id" db:"id"`
	CompanyID   int        `json:"company_id" db:"company_id"`
	Name        string     `json:"name" db:"name"`
	Description *string    `json:"description,omitempty" db:"description"` // nullable
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`       // nullable
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// for a given admin, show projects they are related to
func AdminGetProjects(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "userID missing")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx,
		`
		SELECT DISTINCT p.*
		FROM projects p
		JOIN project_assignees pa ON p.id = pa.project_id
		JOIN users u ON pa.user_id = u.id
		WHERE u.id = $1 AND u.is_admin = TRUE AND p.company_id = u.company_id
		`, userID,
	)
	if err != nil {
		return &fiber.Error{
			Message: "error when getting projects for admin",
		}
	}
	defer rows.Close()

	projects := make([]Project, 0)
	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.ID,
			&project.CompanyID,
			&project.Name,
			&project.Description,
			&project.DueDate,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "error scanning project row")
		}
		projects = append(projects, project)
	}
	err = rows.Err()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error scanning project rows")
	}

	return c.JSON(projects)
}

// for a given user, show projects they are related to
// using local.userID, get project IDs from project assignees table
func UserGetProjects(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "userID missing")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := database.DB.QueryContext(ctx,
		`
		SELECT DISTINCT p.*
		FROM projects p 
		JOIN project_assignees pa ON p.id = pa.project_id
		JOIN users u ON pa.user_id = u.id
		WHERE u.id = $1 AND p.company_id = u.company_id
		`,
		userID)
	if err != nil {
		return &fiber.Error{
			Message: "error when getting projects for user",
		}
	}
	defer rows.Close()

	projects := make([]Project, 0)
	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.ID,
			&project.CompanyID,
			&project.Name,
			&project.Description,
			&project.DueDate,
			&project.CreatedAt,
			&project.UpdatedAt,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "error scanning project row")
		}
		projects = append(projects, project)
	}
	err = rows.Err()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "error scanning project rows")
	}
	return c.JSON(projects)
}

type Task struct {
	ID          int        `json:"id" db:"id"`
	CreatorID   int        `json:"company_id" db:"company_id"`
	ProjectID   int        `json:"project_id" db:"project_id"`
	Title       string     `json:"title" db:"title"`
	Description *string    `json:"description,omitempty" db:"description"`
	Status      string     `json:"status" db:"status"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// for a given admin, show tasks they can use
func AdminGetTasks(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "userID missing")
	}

	return c.SendString("")
}

// for a given user, show tasks they are assigned
func UserGetTasks(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "userID missing")
	}
	// var tasks Tasks
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// err := database.DB.QueryContext(ctx, `SELECT FROM Tasks WHERE`).Scan(&tasks)
	// if err != nil {

	// }
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
