package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nei7/bank/internal/db"
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

	server.router = router
	return server
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}
