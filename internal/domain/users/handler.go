package users

import (
	"encoding/json"
	"net/http"

	"go-crud-api/pkg/web"

	"github.com/go-playground/validator/v10"
)

// AuthHandler handles authentication requests.
type AuthHandler struct {
	service  *Service
	validate *validator.Validate
}

// NewAuthHandler creates a new AuthHandler.
func NewAuthHandler(service *Service) *AuthHandler {
	return &AuthHandler{
		service:  service,
		validate: validator.New(),
	}
}

// RegisterRequest is the request payload for user registration.
type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginRequest is the request payload for user login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse is the response payload for user login.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Register handles user registration.
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} web.Response{data=User} "User registered successfully"
// @Failure 400 {object} web.Response{error=web.ApiError} "Bad request or validation error"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondWithError(w, "bad_request", "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		web.RespondWithError(w, "validation_error", err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		// TODO: Handle specific errors, e.g., email already exists
		web.RespondWithError(w, "internal_error", "Could not create user", http.StatusInternalServerError)
		return
	}

	web.RespondWithJSON(w, http.StatusCreated, web.Response{Data: user})
}

// @Summary Log in a user
// @Description Authenticate user with email and password, returns JWT tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "User login credentials"
// @Success 200 {object} web.Response{data=LoginResponse} "User logged in successfully"
// @Failure 400 {object} web.Response{error=web.ApiError} "Bad request or validation error"
// @Failure 401 {object} web.Response{error=web.ApiError} "Unauthorized (invalid credentials)"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		web.RespondWithError(w, "bad_request", "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		web.RespondWithError(w, "validation_error", err.Error(), http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// TODO: Handle specific errors, e.g., invalid credentials
		web.RespondWithError(w, "unauthorized", "Invalid email or password", http.StatusUnauthorized)
		return
	}

	resp := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	web.RespondWithJSON(w, http.StatusOK, web.Response{Data: resp})
}

// @Summary List all users
// @Description Get a list of all registered users (Admin only)
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} web.Response{data=[]User} "List of users"
// @Failure 401 {object} web.Response{error=web.ApiError} "Unauthorized"
// @Failure 403 {object} web.Response{error=web.ApiError} "Forbidden (not admin)"
// @Failure 500 {object} web.Response{error=web.ApiError} "Internal server error"
// @Router /v1/users [get]
func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.List(r.Context())
	if err != nil {
		web.RespondWithError(w, "internal_error", "Could not fetch users", http.StatusInternalServerError)
		return
	}

	web.RespondWithJSON(w, http.StatusOK, web.Response{Data: users})
}
