package service

import (
	"context"
	"errors"
	"flagon/pkg/config"
	"flagon/pkg/model"
	"flagon/pkg/repository"
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
	VerifyJwtToken(ctx context.Context, tokenString string) (*jwt.Token, error)
}

type authService struct {
	userRepo  repository.UserRepository
	authCfg   config.Authentication
	tokenRepo repository.TokenRepository
	secretKey []byte
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) AuthService {
	cfg := config.GetConfig()
	return &authService{
		userRepo:  userRepo,
		authCfg:   cfg.Auth,
		tokenRepo: tokenRepo,
		secretKey: []byte(cfg.Auth.Secret),
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
	existingUser, err := s.userRepo.FindByUsernameOrEmail(ctx, req.Username, req.Email)
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

	return user, s.userRepo.Create(ctx, user)
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

type JwtTokenClaim struct {
	jwt.RegisteredClaims
}

func (s *authService) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}

	accessJTI := uuid.New().String()
	refreshJTI := uuid.New().String()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtTokenClaim{jwt.RegisteredClaims{
		ID:        accessJTI,
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authCfg.RefreshTokenLifetime)), // Refresh token expires in 30 days
	}})

	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtTokenClaim{jwt.RegisteredClaims{
		ID:        refreshJTI,
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authCfg.RefreshTokenLifetime)), // Refresh token expires in 30 days
	}})

	refreshTokenString, err := refreshToken.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	// Store tokens in cache
	if err := s.tokenRepo.AddJwtToken(ctx, user.ID, accessJTI, s.authCfg.AccessTokenLifetime); err != nil {
		return nil, err
	}
	if err := s.tokenRepo.AddJwtToken(ctx, user.ID, refreshJTI, s.authCfg.RefreshTokenLifetime); err != nil {
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
	token, err := s.VerifyJwtToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	// Find user
	userID, _ := token.Claims.GetSubject()
	user, err := s.userRepo.FindByID(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessJTI := uuid.New().String()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtTokenClaim{jwt.RegisteredClaims{
		ID:        accessJTI,
		Subject:   user.ID.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.authCfg.RefreshTokenLifetime)), // Refresh token expires in 30 days
	}})

	accessTokenString, err := accessToken.SignedString(s.secretKey)
	if err != nil {
		return nil, err
	}

	// Store new access token in cache
	if err := s.tokenRepo.AddJwtToken(ctx, user.ID, accessJTI, s.authCfg.AccessTokenLifetime); err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessTokenString,
		RefreshToken: req.RefreshToken, // Return the same refresh token
		User:         user,
	}, nil
}

func (s *authService) Logout(ctx context.Context, userID uuid.UUID, jti string) error {
	return s.tokenRepo.RemoveJwtToken(ctx, userID, jti)
}

func (s *authService) VerifyJwtToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JwtTokenClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if !token.Valid {
		return nil, errors.New("refresh token is invalid")
	}

	// Get claims
	claims, ok := token.Claims.(*JwtTokenClaim)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	// Check if refresh token is valid in repository
	isValid, err := s.tokenRepo.IsJwtValid(ctx, userUUID, claims.ID)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, errors.New("refresh token has been revoked")
	}

	return token, nil
}
