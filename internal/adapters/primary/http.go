package primary

import (
	"github.com/gin-gonic/gin"
	"github.com/roxxers/surfe-techtest/internal/core/services"
	"github.com/roxxers/surfe-techtest/internal/ports"
)

type HTTPServer struct {
	// Usually have defined values for what IP to bind on and port number here
	db      ports.Database
	service *services.Service
}

func NewHTTPServer(db ports.Database, service *services.Service) *HTTPServer {
	return &HTTPServer{db, service}
}

func (s *HTTPServer) Serve(addr string) {
	router := gin.New()
	controller := NewController(s.service)

	// Add any required auth middleware here and logging middleware

	// Using GET to sim a CRUD like REST API
	router.GET("/api/v1/user/:id", controller.FetchUser)
	router.GET("/api/v1/user/:id/actioncount", controller.GetUserActionCount)
	router.POST("/api/v1/actions/probablity", controller.CalculateNextActionProbablity)
	router.GET("/api/v1/users/referalindex", controller.CalculateAllUserReferalIndexes)
	router.Run(addr)
}
