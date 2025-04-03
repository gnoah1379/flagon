package v1

import (
	"flagon/pkg/api/v1/response"
	"flagon/pkg/repository"
	"flagon/pkg/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthAPI interface {
	Register(router gin.IRouter)
	AuthRequired() gin.HandlerFunc
}

type authApi struct {
	authService service.AuthService
	tokenRepo   repository.TokenRepository
}

func NewAuthAPI(authService service.AuthService, tokenRepo repository.TokenRepository) AuthAPI {
	return &authApi{
		authService: authService,
		tokenRepo:   tokenRepo,
	}
}

func (api *authApi) Register(router gin.IRouter) {
	router.POST("/register", api.HandleRegister)
	router.POST("/login", api.HandleLogin)
	router.POST("/refresh-token", api.HandleRefreshToken)
	router.POST("/logout", api.HandleLogout, api.AuthRequired())
}

// HandleRegister
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param user body service.RegisterRequest true "User registration details"
// @Success 200 {object} response.SuccessResponse[model.User] "Registration successful"
// @Failure 400 {object} response.ErrorResponse[string] "Invalid request"
// @Router /register [post]
func (api *authApi) HandleRegister(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendBadRequest(c, response.ErrInvalidRequest, err.Error())
		return
	}

	user, err := api.authService.Register(c, &req)
	if err != nil {
		response.SendBadRequest(c, response.ErrInvalidRequest, err.Error())
		return
	}

	response.SendOK(c, "Registration successful", user)
}

// HandleLogin
// @Summary Login user
// @Description Authenticate user and return access token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body service.LoginRequest true "User login credentials"
// @Success 200 {object} response.SuccessResponse[service.LoginResponse] "Login successful"
// @Failure 401 {object} response.ErrorResponse[string] "Invalid credentials"
// @Router /login [post]
func (api *authApi) HandleLogin(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendBadRequest(c, response.ErrInvalidRequest, err.Error())
		return
	}

	resp, err := api.authService.Login(c.Request.Context(), &req)
	if err != nil {
		response.SendUnauthorized(c, response.ErrInvalidCredentials, err.Error())
		return
	}

	response.SendOK(c, "Login successful", resp)
}

// HandleRefreshToken
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body service.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} response.SuccessResponse[service.LoginResponse] "Token refreshed successfully"
// @Failure 401 {object} response.ErrorResponse[string] "Invalid token"
// @Router /refresh-token [post]
func (api *authApi) HandleRefreshToken(c *gin.Context) {
	var req service.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.SendBadRequest(c, response.ErrInvalidRequest, err.Error())
		return
	}

	resp, err := api.authService.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		response.SendUnauthorized(c, response.ErrInvalidToken, err.Error())
		return
	}

	response.SendOK(c, "Token refreshed successfully", resp)
}

// HandleLogout
// @Summary Logout user
// @Description Invalidate current access token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse[string]	"Logout successful"
// @Failure 401 {object} response.ErrorResponse[string] "Unauthorized"
// @Failure 500 {object} response.ErrorResponse[string] "Internal server error"
// @Router /logout [post]
func (api *authApi) HandleLogout(c *gin.Context) {
	// Get user ID and JTI from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		response.SendUnauthorized(c, response.ErrUnauthorized, "User not authenticated")
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid user ID")
		return
	}

	jti, exists := c.Get("jti")
	if !exists {
		response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid token")
		return
	}

	err = api.authService.Logout(c.Request.Context(), userUUID, jti.(string))
	if err != nil {
		response.SendInternalServerError(c, response.ErrInternalServer, "Failed to logout")
		return
	}

	response.SendOK(c, "Logout successful", nil)
}

func (api *authApi) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Check if token starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := api.authService.VerifyJwtToken(c, tokenString)

		if err != nil {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		claims, _ := token.Claims.(*service.JwtTokenClaim)

		userUUID, err := uuid.Parse(claims.Subject)
		if err != nil {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid user ID in token")
			c.Abort()
			return
		}

		// Check if token is valid in repository
		isValid, err := api.tokenRepo.IsJwtValid(c, userUUID, claims.ID)
		if err != nil {
			response.SendInternalServerError(c, response.ErrInternalServer, "Failed to validate token")
			c.Abort()
			return
		}

		if !isValid {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Token has been revoked")
			c.Abort()
			return
		}

		// Set user ID and JTI in context
		c.Set("userID", userUUID)
		c.Set("jti", claims.ID)
		c.Set("claims", claims)
		c.Next()
	}
}
