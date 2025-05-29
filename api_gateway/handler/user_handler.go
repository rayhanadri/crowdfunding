package handler

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"

	"crowdfund/model"
	"crowdfund/repository"
)

type UserHandler interface {
	// AllUsers(c echo.Context) error
	GetUserByID(c echo.Context) error
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
	LoginUser(c echo.Context) error
	RefreshToken(c echo.Context) error
}

type userHandler struct {
	userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) UserHandler {
	return &userHandler{userRepo: userRepo}
}

// GetUser godoc
// @Summary Get Current User Details
// @Description Current User Details
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Success 200 {object} model.Response
// @Router /users/me [get]
func (h *userHandler) GetUserByID(c echo.Context) error {
	// fmt.Println("GetUserByID called")
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)

	idInt := userIdInt

	// fmt.Println("User ID from context:", idInt)
	// fmt.Println("User ID from JWT:", userID)

	user, err := h.userRepo.GetUserByID(idInt)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve user",
		})
	}

	user.Password = "" // Clear the password before sending the response
	return c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "Success",
		Data:    user,
	})
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserRegister true "User object"
// @Success 201 {object} model.Response
// @Router /users/register [post]
func (h *userHandler) CreateUser(c echo.Context) error {

	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	err := user.Validate()
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	user, err = h.userRepo.CreateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to create user, maybe email already exists",
		})
	}

	user.Password = "" // Clear the password before sending the response
	return c.JSON(http.StatusCreated, model.Response{
		Status:  http.StatusCreated,
		Message: "Success",
		Data:    user,
	})
}

// LoginUser godoc
// @Summary Login user
// @Description Authenticate and return a token for the user
// @Tags users
// @Accept json
// @Produce json
// @Param user body model.UserLogin true "User object"
// @Success 200 {object} model.Response
// @Router /users/login [post]
func (h *userHandler) LoginUser(c echo.Context) error {
	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// fmt.Println("User login request:", user)

	user, err := h.userRepo.LoginUser(user)
	// fmt.Println("User after login:", user)

	if err != nil {
		// fmt.Println("Password mismatch:", err)
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Invalid email or password",
		})
	}

	accessToken, refreshToken, err := GenerateTokens(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate tokens" + err.Error(),
		})
	}

	tokenString := accessToken
	// Parse the token to extract claims
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_KEY")), nil
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid token",
		})
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid token claims",
		})
	}

	// set values in context
	c.Set("id", claims.UserID)
	c.Set("email", claims.Email)
	c.Set("exp", claims.Exp)

	// return the user object
	return c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "Login successful",
		Data: map[string]interface{}{
			"user":         user,
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		},
	})
}

// RefreshToken godoc
// @Summary Refresh user access token
// @Description Refresh the access token using the refresh token
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Success 200 {object} model.Response
// @Router /users/refresh-token [post] // Updated the router path to use POST method
func (h *userHandler) RefreshToken(c echo.Context) error {
	// Get the refresh token from the request header
	refreshToken := c.Request().Header.Get("Authorization")
	if refreshToken == "" {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "Missing refresh token",
		})
	}

	// Parse the refresh token to extract claims
	refreshToken = strings.TrimPrefix(refreshToken, "Bearer ")
	token, err := jwt.ParseWithClaims(refreshToken, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_KEY")), nil
	})
	if err != nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid refresh token",
		})
	}

	claims, ok := token.Claims.(*model.Claims)
	if !ok || !token.Valid {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "Invalid token claims",
		})
	}

	userID := claims.UserID

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to retrieve user",
		})
	}

	newAccessToken, newRefreshToken, err := GenerateTokens(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate new tokens",
		})
	}

	// remove password from user object before sending response
	user.Password = "" // Clear the password before sending the response

	return c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "Tokens refreshed successfully",
		Data: map[string]interface{}{
			"accessToken":  newAccessToken,
			"refreshToken": newRefreshToken,
			"user":         user,
		},
	})
}

// Create Refresh Token
func GenerateTokens(user *model.User) (string, string, error) {
	accessClaims := model.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Exp:    float64(time.Now().Add(time.Hour * 24).Unix()), // Token expires in 24 hours
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(os.Getenv("JWT_ACCESS_KEY")))
	if err != nil {
		log.Println("Error generating access token:", err.Error())
		return "", "", err
	}

	// Refresh Token (7 days)
	refreshClaims := model.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Exp:    float64(time.Now().Add(time.Hour * 24).Unix()),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(os.Getenv("JWT_REFRESH_KEY")))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// UpdateUser godoc
// @Summary Update user details
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>"
// @Param user body model.UserRegister true "User object"
// @Success 200 {object} model.Response
// @Router /users/me [put] // Updated the router path to include user ID
func (h *userHandler) UpdateUser(c echo.Context) error {
	userID := c.Get("user_id")
	if userID == nil {
		return c.JSON(http.StatusUnauthorized, model.Response{
			Status:  http.StatusUnauthorized,
			Message: "User not authenticated",
		})
	}

	userIdFloat, ok := userID.(float64)
	if !ok {
		return c.JSON(http.StatusInternalServerError, model.Response{
			Status:  http.StatusInternalServerError,
			Message: "Invalid user ID",
		})
	}
	userIdInt := int(userIdFloat)

	idInt := userIdInt

	user := new(model.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, model.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	}
	user.ID = idInt
	user.UpdatedAt = time.Now()
	updatedUser, err := h.userRepo.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	updatedUser.Password = "" // Clear the password before sending the response
	return c.JSON(http.StatusOK, model.Response{
		Status:  http.StatusOK,
		Message: "User updated successfully",
		Data: map[string]interface{}{
			"user": updatedUser,
		},
	})
}
