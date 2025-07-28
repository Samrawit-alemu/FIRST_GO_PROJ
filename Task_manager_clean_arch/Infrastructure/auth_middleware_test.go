package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"os"
	"taskmanager/domain"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setupTestRouter(t *testing.T, jwtService IJWTService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	testHandler := func(c *gin.Context) {
		role, _ := c.Get("role")
		userID, _ := c.Get("user_id")

		assert.Equal(t, "user", role)
		assert.NotEmpty(t, userID)
		c.Status(http.StatusOK)
	}
	router.GET("/test", AuthMiddleware(jwtService), testHandler)
	return router
}

func TestAuthMiddleware(t *testing.T) {
	os.Setenv("JWT_SECRET", "a_secret_for_testing")
	defer os.Unsetenv("JWT_SECRET")

	jwtService := NewJWTService()
	testUser := domain.User{ID: primitive.NewObjectID(), Role: "user"}
	validToken, err := jwtService.GenerateToken(testUser)
	assert.NoError(t, err)

	router := setupTestRouter(t, jwtService)

	t.Run("Success - Valid Token", func(t *testing.T) {

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Failure - No Authorization Header", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Authorization header required")
	})

	t.Run("Failure - Invalid Token Format", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "invalid-format "+validToken)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid token format")
	})

	t.Run("Failure - Malformed or Expired Token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer malformed.jwt.token")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		assert.Contains(t, rr.Body.String(), "Invalid or expired token")
	})
}

func TestRoleAuthMiddleware(t *testing.T) {

	gin.SetMode(gin.TestMode)

	protectedHandler := func(c *gin.Context) {
		c.Status(http.StatusOK)
	}

	t.Run("Success - User has required role", func(t *testing.T) {
		router := gin.New()
		setContextMiddleware := func(c *gin.Context) {
			c.Set("role", "admin")
			c.Next()
		}
		router.GET("/protected", setContextMiddleware, RoleAuthMiddleware("admin"), protectedHandler)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Failure - User has wrong role", func(t *testing.T) {
		router := gin.New()
		setContextMiddleware := func(c *gin.Context) {
			c.Set("role", "user")
			c.Next()
		}
		router.GET("/protected", setContextMiddleware, RoleAuthMiddleware("admin"), protectedHandler)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)

		router.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusForbidden, rr.Code)
		assert.Contains(t, rr.Body.String(), "Forbidden")
	})

	t.Run("Failure - Role not set in context", func(t *testing.T) {
		router := gin.New()
		router.GET("/protected", RoleAuthMiddleware("admin"), protectedHandler)

		rr := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/protected", nil)

		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusForbidden, rr.Code)
	})
}
