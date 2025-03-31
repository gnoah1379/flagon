package service

import (
	"context"
	"errors"
	"flagon/model"
	"flagon/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, req *RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error)
	RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*LoginResponse, error)
	Logout(ctx context.Context, userID uuid.UUID, jti string) error
}

type authService struct {
	userRepo  repository.UserRepository
	jwtKey    []byte
	tokenRepo repository.TokenRepository
}

func NewAuthService(userRepo repository.UserRepository, jwtKey []byte, tokenRepo repository.TokenRepository) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtKey:    jwtKey,
		tokenRepo: tokenRepo,
	}
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Email     string `json:"email" binding:"required,email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}

func (s *authService) Register(ctx context.Context, req *RegisterRequest) (*model.User, error) {
	// Check if username or email already exists
	existingUser, err := s.userRepo.FindByUsernameOrEmail(req.Username, req.Email)
	if err == nil {
		if existingUser.Username == req.Username {
			return nil, errors.New("username already exists")
		}
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &model.User{
		ID:        uuid.New(),
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		AvatarURL: req.AvatarURL,
	}

	return user, s.userRepo.Create(user)
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *model.User
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Generate JTI for access token
	accessJTI := uuid.New().String()
	// Generate JTI for refresh token
	refreshJTI := uuid.New().String()

	// Generate access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"jti": accessJTI,
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Access token expires in 1 hour
	})

	accessTokenString, err := accessToken.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"jti": refreshJTI,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // Refresh token expires in 7 days
	})

	refreshTokenString, err := refreshToken.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	// Store tokens in cache
	if err := s.tokenRepo.AddUserToken(ctx, user.ID, accessJTI, time.Hour*1); err != nil {
		return nil, err
	}
	if err := s.tokenRepo.AddUserToken(ctx, user.ID, refreshJTI, time.Hour*24*7); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		User:         user,
	}, nil
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (s *authService) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*LoginResponse, error) {
	// Parse refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.jwtKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if !token.Valid {
		return nil, errors.New("refresh token is invalid")
	}

	// Get claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Get user ID and JTI from claims
	userID, ok := claims["sub"]
	if !ok {
		return nil, errors.New("invalid user ID in token")
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("invalid JTI in token")
	}

	// Check if refresh token is valid in repository
	isValid, err := s.tokenRepo.IsUserTokenValid(ctx, userUUID, jti)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, errors.New("refresh token has been revoked")
	}

	// Find user
	user, err := s.userRepo.FindByID(userUUID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessJTI := uuid.New().String()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"jti": accessJTI,
		"exp": time.Now().Add(time.Hour * 1).Unix(), // Access token expires in 1 hour
	})

	accessTokenString, err := accessToken.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	// Store new access token in cache
	if err := s.tokenRepo.AddUserToken(ctx, user.ID, accessJTI, time.Hour*1); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: req.RefreshToken, // Return the same refresh token
		User:         user,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID, jti string) error {
	return s.tokenRepo.RemoveUserToken(ctx, userID, jti)
}
