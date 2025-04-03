package v1

import (
	"github.com/gin-gonic/gin"
)

type API interface {
	Register(router gin.IRouter)
}

func New(authAPI AuthAPI) API {
	return &api{
		Auth: authAPI,
	}

}

type api struct {
	Auth AuthAPI
}

func (a *api) Register(r gin.IRouter) {
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		a.Auth.Register(v1)
		protected := v1.Group("/", a.Auth.AuthRequired())
		{
			_ = protected
		}
	}
}
