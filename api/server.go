package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/reinhardbuyabo/simplebank/db/sqlc"
)

// Server servers HTTP requests for our banking service
type Server struct {
	store  *db.Store   // allow us to interact with database when processing api request
	router *gin.Engine // allow us to send each API request to the correct handler for processing
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	// add routes to the router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

// start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
