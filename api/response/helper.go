package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendError sends an error response
func SendError(c *gin.Context, code int, message string, details any) {
	c.JSON(code, NewErrorResponse(code, message, details))
}

// SendSuccess sends a success response
func SendSuccess(c *gin.Context, code int, message string, data any) {
	c.JSON(code, NewSuccessResponse(code, message, data))
}

// SendBadRequest sends a 400 Bad Request response
func SendBadRequest(c *gin.Context, message string, details any) {
	SendError(c, http.StatusBadRequest, message, details)
}

// SendUnauthorized sends a 401 Unauthorized response
func SendUnauthorized(c *gin.Context, message string, details any) {
	SendError(c, http.StatusUnauthorized, message, details)
}

// SendForbidden sends a 403 Forbidden response
func SendForbidden(c *gin.Context, message string, details any) {
	SendError(c, http.StatusForbidden, message, details)
}

// SendNotFound sends a 404 Not Found response
func SendNotFound(c *gin.Context, message string, details any) {
	SendError(c, http.StatusNotFound, message, details)
}

// SendInternalServerError sends a 500 Internal Server Error response
func SendInternalServerError(c *gin.Context, message string, details any) {
	SendError(c, http.StatusInternalServerError, message, details)
}

// SendOK sends a 200 OK response
func SendOK(c *gin.Context, message string, data any) {
	SendSuccess(c, http.StatusOK, message, data)
}

// SendCreated sends a 201 Created response
func SendCreated(c *gin.Context, message string, data any) {
	SendSuccess(c, http.StatusCreated, message, data)
}
