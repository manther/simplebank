package api	

import (
	db "github.com/manther/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)
type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	// router.POST("/transfers", server.createTransfer)
	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// TODO: add graceful shutdown
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}