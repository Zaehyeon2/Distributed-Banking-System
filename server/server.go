package server

import (
	"net/http"
	"time"

	"dbs/server/banking_handler"
	"dbs/server/raft_handler"

	"github.com/hashicorp/raft"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type srv struct {
	listenAddress string
	raft          *raft.Raft
	echo          *echo.Echo
}

// Start start the server
func (s srv) Start() error {
	return s.echo.StartServer(&http.Server{
		Addr:         s.listenAddress,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

// New return new server
func New(listenAddr string, r *raft.Raft) *srv {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Pre(middleware.RemoveTrailingSlash())
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	// Raft server
	raftHandler := raft_handler.New(r)
	e.POST("/raft", raftHandler.JoinRaftHandler)
	e.DELETE("/raft/:node_id", raftHandler.RemoveRaftHandler)
	e.GET("/raft/leaderstats", raftHandler.StatsLeaderHandler)
	e.GET("/raft/nodesstats", raftHandler.StatsNodesHandler)

	// Banking handler
	bankingHandler := banking_handler.New(r)
	e.POST("/bank", bankingHandler.Deposit)
	e.PUT("/bank", bankingHandler.Transfer)
	e.GET("/bank/:account", bankingHandler.Get)

	return &srv{
		listenAddress: listenAddr,
		echo:          e,
		raft:          r,
	}
}
