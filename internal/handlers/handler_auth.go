package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
	"wallet-app/internal/auth"
	"wallet-app/internal/database"
	"wallet-app/internal/utils"
)

type AuthHandler struct {
	DB *database.Queries
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// @Summary Register a new user
// @Description Register a new user with username, email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration details"
// @Success 201 {object} models.AuthResponseSwag "User registered successfully"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 409 {object} map[string]string "User already exists"
// @Router /v1/register [post]

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	_, err = h.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	})

	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			utils.RespondWithError(w, http.StatusBadRequest, "Username or email already exists")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{
		"message": "User registered successfully",
	})
}

// @Example {object} {"username": "johndoe", "password": "secretpass123"}
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login credentials"
// @Success 200 {object} models.AuthResponseSwag "Authentication successful"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} utils.ErrorResponse "Invalid credentials"
// @Router /v1/login [post]

// Login takes email and password
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := h.DB.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid username or password")
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}
