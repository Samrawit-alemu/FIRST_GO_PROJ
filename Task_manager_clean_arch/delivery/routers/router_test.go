package routers

import (
	"net/http"
	"net/http/httptest"
	"taskmanager/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRouter_AdminRouteIsProtected(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockJwtService := new(mocks.IJWTService)
	mockUserController := new(mocks.IUserController)
	mockTaskController := new(mocks.ITaskController)

	router := SetupRouter(mockUserController, mockTaskController, mockJwtService)

	req, _ := http.NewRequest(http.MethodPut, "/admin/promote/123", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}
