package mcp

import (
	"github.com/mark3labs/mcp-go/server"
	"github.com/typemd/typemd/core"
)

// Serve creates and starts the MCP server via stdio.
func Serve(vault *core.Vault) error {
	s := server.NewMCPServer(
		"tmd",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	registerTools(s, vault)

	return server.ServeStdio(s)
}

// NewServer creates an MCP server with tools registered (for testing).
func NewServer(vault *core.Vault) *server.MCPServer {
	s := server.NewMCPServer(
		"tmd",
		"0.1.0",
		server.WithToolCapabilities(false),
	)

	registerTools(s, vault)

	return s
}
