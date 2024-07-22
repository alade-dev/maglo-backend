package handler

import (
	"errors"
	"net/mail"
	"time"

	"maglo/config"
	"maglo/database"
	"maglo/model"

	"gorm.io/gorm"

	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// CheckPasswordHash compare password with hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func isEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	type UserData struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(LoginInput)
	var userData UserData

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	if input.Identity == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Identity or password is empty", "data": nil})
	}

	identity := input.Identity
	pass := input.Password
	userModel, err := new(model.User), *new(error)

	if isEmail(identity) {
		userModel, err = getUserByEmail(identity)
	} else {
		userModel, err = getUserByUsername(identity)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if userModel == nil {
		CheckPasswordHash(pass, "")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		userData = UserData{
			ID:       userModel.ID,
			Username: userModel.Username,
			Email:    userModel.Email,
			Password: userModel.Password,
		}
	}

	if !CheckPasswordHash(pass, userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": nil})
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = userData.Username
	claims["email"] = userData.Email
	claims["user_id"] = userData.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}

func Logout(c *fiber.Ctx) error {
	// Destroy token
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(-time.Hour)

	_, err := token.SignedString([]byte(config.Config("SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Logout success", "data": nil})
}

// ConfigGoogle to set config of oauth
func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("Client"),
		ClientSecret: config.Config("Secret"),
		RedirectURL:  config.Config("redirect_url"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}
	return conf
}

// GetUserInfo fetches user's email, name from google
func GetUserInfo(token string) (model.GoogleResponse, error) {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		panic(err)
	}
	ptoken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {ptoken}},
	}
	req, err := http.DefaultClient.Do(res)
	if err != nil {
		return model.GoogleResponse{}, err
	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return model.GoogleResponse{}, err
	}
	var data model.GoogleResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	return data, nil
}

// AuthGoogle to authenticate user with google
func Auth(c *fiber.Ctx) error {
	path := ConfigGoogle()
	url := path.AuthCodeURL("state")
	return c.Redirect(url)
}

func Callback(c *fiber.Ctx) error {
	token, err := ConfigGoogle().Exchange(c.Context(), c.FormValue("code"))
	if err != nil {
		panic(err)
	}
	userInfo, err := GetUserInfo(token.AccessToken)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Unable to get user info"})
	}

	var user model.User
	if err := database.DB.Where("email = ?", userInfo.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User not found, create a new user
			user = model.User{
				Email: userInfo.Email,
				Names: userInfo.Name,
				// Picture: userInfo.Picture,
			}
			if err := database.DB.Create(&user).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{"error": "Unable to create user"})
			}
		} else {
			return c.Status(500).JSON(fiber.Map{"error": "Database error"})
		}
	} else {
		// User exists, update user info if needed
		user.Names = userInfo.Name
		// user.Picture = userInfo.Picture
		if err := database.DB.Save(&user).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Unable to update user"})
		}
	}

	return c.Status(200).JSON(fiber.Map{"email": userInfo.Email, "full_name": userInfo.Name, "login": true})
}
