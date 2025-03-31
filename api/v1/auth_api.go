package v1

import (
	"flagon/api/response"
	"flagon/repository"
	"flagon/service"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthAPI interface {
	Register(router gin.IRouter)
	AuthRequired() gin.HandlerFunc
}

type authApi struct {
	authService service.AuthService
	tokenRepo   repository.TokenRepository
	jwtKey      []byte
}

func NewAuthAPI(authService service.AuthService, tokenRepo repository.TokenRepository, jwtKey []byte) AuthAPI {
	return &authApi{
		authService: authService,
		tokenRepo:   tokenRepo,
		jwtKey:      jwtKey,
	}
}

func (api *authApi) Register(router gin.IRouter) {
	router.POST("/signup", api.handleRegister)
	router.POST("/signin", api.handleLogin)
	router.POST("/refresh-token", api.handleRefreshToken)
	router.POST("/logout", api.handleLogout, api.AuthRequired())

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
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return api.jwtKey, nil
		})

		if err != nil {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		// Check if token is valid
		if !token.Valid {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Token is invalid")
			c.Abort()
			return
		}

		// Get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid token claims")
			c.Abort()
			return
		}

		// Get user ID and JTI from claims
		userID, ok := claims["sub"]
		if !ok {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid user ID in token")
			c.Abort()
			return
		}

		userUUID, err := uuid.Parse(userID.(string))
		if err != nil {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid user ID in token")
			c.Abort()
			return
		}

		jti, ok := claims["jti"].(string)
		if !ok {
			response.SendUnauthorized(c, response.ErrUnauthorized, "Invalid JTI in token")
			c.Abort()
			return
		}

		// Check if token is valid in repository
		isValid, err := api.tokenRepo.IsUserTokenValid(c, userUUID, jti)
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
		c.Set("userID", userID)
		c.Set("jti", jti)
		c.Next()
	}
}

func (api *authApi) handleRegister(c *gin.Context) {
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

func (api *authApi) handleLogin(c *gin.Context) {
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

func (api *authApi) handleRefreshToken(c *gin.Context) {
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

func (api *authApi) handleLogout(c *gin.Context) {
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
