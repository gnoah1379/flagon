package v1

import (
	"flagon/internal/repository"
	"flagon/internal/service"
	"github.com/gin-gonic/gin"
)

type API interface {
	Register(router gin.IRouter)
}

type api struct {
	authAPI AuthAPI
}

func NewAPI(authService service.AuthService, tokenRepo repository.TokenRepository, jwtKey []byte) *api {
	return &api{
		authAPI: NewAuthAPI(authService, tokenRepo, jwtKey),
	}
}

func (a *api) Register(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		a.authAPI.Register(v1)
		protected := v1.Group("/", a.authAPI.AuthRequired())
		{
			_ = protected
		}
	}
}
