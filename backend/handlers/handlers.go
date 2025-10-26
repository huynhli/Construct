package handlers

import (
	"backend/database"
	"context"
	"database/sql"
	"encoding/json"
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
	Username  string `json:"username"`
	Password  string `json:"password"`
	CompanyID int    `json:"company_id"`
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
		_, err := database.DB.ExecContext(ctx, "INSERT INTO users (company_id, username, password_hash) VALUES ($1, $2, $3)", signInReq.CompanyID, signInReq.Username, string(hashedBytePW))
		if err != nil {
			return &fiber.Error{
				Code:    fiber.ErrBadGateway.Code,
				Message: err.Error(),
			}
		}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "user created"})

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
}

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
func UserGetProjects(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "userID missing")
	}

	// for a given user/admin, show projects they are related to
	// using local.userID, get project IDs from project assignees table
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
	CreatorID   int        `json:"creator_id" db:"creator_id"`
	ProjectID   int        `json:"project_id" db:"project_id"`
	Title       string     `json:"title" db:"title"`
	Description *string    `json:"description,omitempty" db:"description"`
	Status      string     `json:"status" db:"status"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	Assignees   []User     `json:"assignees" db:"assignees"`
}

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
}

func UserGetTasks(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "userID missing")
	}
	isAdmin := c.Locals("isAdmin")
	if isAdmin == false {
		// for user, show tasks they are assigned, and who else has been assigned to the same task
		projectID := c.Query("projectID")

		// want all tasks for this project
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := database.DB.QueryContext(ctx,
			`
			SELECT 
				t.id,
				t.creator_id,
				t.project_id,
				t.title,
				t.description,
				t.status,
				t.due_date,
				t.created_at,
				t.updated_at,
				COALESCE(
					json_agg(
						DISTINCT jsonb_build_object(
							'id', u.id,
							'username', u.username
						)
					) FILTER (WHERE u.id IS NOT NULL),
					'[]'
				) AS assignees
			FROM tasks t
			JOIN task_assignees ta ON t.id = ta.task_id
			LEFT JOIN users u ON ta.user_id = u.id
			WHERE ta.user_id = $1 AND t.project_id = $2
			GROUP BY t.id
			`, userID, projectID,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "query error: "+err.Error())
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			var assigneesJSON []byte

			err := rows.Scan(
				&task.ID,
				&task.CreatorID,
				&task.ProjectID,
				&task.Title,
				&task.Description,
				&task.Status,
				&task.DueDate,
				&task.CreatedAt,
				&task.UpdatedAt,
				&assigneesJSON,
			)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "scan error: "+err.Error())
			}

			err = json.Unmarshal(assigneesJSON, &task.Assignees)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "unmarshal error: "+err.Error())
			}
			tasks = append(tasks, task)
		}
		err = rows.Err()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(tasks)
	} else {
		// for a given admin, show all tasks they can use
		projectID := c.Query("projectID")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		rows, err := database.DB.QueryContext(ctx,
			`
			SELECT 
				t.id,
				t.creator_id,
				t.project_id,
				t.title,
				t.description,
				t.status,
				t.due_date,
				t.created_at,
				t.updated_at,
				COALESCE(
					json_agg(
						DISTINCT jsonb_build_object(
							'id', u.id,
							'username', u.username
						)
					) FILTER (WHERE u.id IS NOT NULL),
					'[]'
				) AS assignees
			FROM tasks t
			LEFT JOIN task_assignees ta ON t.id = ta.task_id
			LEFT JOIN users u ON ta.user_id = u.id
			WHERE t.project_id = $1
			GROUP BY t.id
			`, projectID,
		)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "query error: "+err.Error())
		}
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var task Task
			var assigneesJSON []byte

			err := rows.Scan(
				&task.ID,
				&task.CreatorID,
				&task.ProjectID,
				&task.Title,
				&task.Description,
				&task.Status,
				&task.DueDate,
				&task.CreatedAt,
				&task.UpdatedAt,
				&assigneesJSON,
			)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "scan error: "+err.Error())
			}

			err = json.Unmarshal(assigneesJSON, &task.Assignees)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "unmarshal error: "+err.Error())
			}
			tasks = append(tasks, task)
		}
		err = rows.Err()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.JSON(tasks)
	}
}

func AdminAddOrEditProject(c *fiber.Ctx) error {
	return c.SendString("meow")
}

func AdminDeleteProject(c *fiber.Ctx) error {
	return c.SendString("meow")
}

func AdminAddOrEditTask(c *fiber.Ctx) error {
	return c.SendString("meow")
}

func DeleteTask(c *fiber.Ctx) error {
	return c.SendString("meow")
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
