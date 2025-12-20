package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/complytime/gemara-mcp-server/tools/authoring"
	"github.com/complytime/gemara-mcp-server/tools/info"
	"github.com/complytime/gemara-mcp-server/tools/prompts"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	oauth "github.com/tuannvm/oauth-mcp-proxy"
	"github.com/tuannvm/oauth-mcp-proxy/mark3labs"
)

type ServerConfig struct {
	// Version of the server
	Version string
	// Transport mode selection (stdio/streamable)
	Transport string

	// Host for StreamableHTTP transport
	Host string
	// Port for StreamableHTTP transport
	Port int
	// Logger for HTTP server logging
	Logger *slog.Logger
}

// Server represents the MCP server
type Server struct {
	mcpServer *server.MCPServer
	config    *ServerConfig
}

// NewServer creates a new MCP server
func NewServer(cfg *ServerConfig) (*Server, error) {
	slog.Debug("Creating new MCP server", "version", cfg.Version)

	s := &Server{
		config: cfg,
	}

	mcpServer := server.NewMCPServer(
		"gemara-mcp-server",
		cfg.Version,
		server.WithLogging(),
	)

	s.mcpServer = mcpServer
	slog.Debug("MCP server instance created")

	// Register Gemara Info Tools (validation, schemas, resources)
	slog.Debug("Initializing Gemara info tools")
	infoTools, err := info.NewGemaraInfoTools()
	if err != nil {
		slog.Error("Failed to create info tools", "error", err)
		return nil, err
	}
	infoTools.Register(mcpServer)
	slog.Debug("Gemara info tools registered successfully")

	// Register Gemara Authoring Tools
	slog.Debug("Initializing Gemara authoring tools")
	authoringTools, err := authoring.NewGemaraAuthoringTools()
	if err != nil {
		slog.Error("Failed to create authoring tools", "error", err)
		return nil, err
	}
	authoringTools.Register(mcpServer)
	slog.Debug("Gemara authoring tools registered successfully")

	return s, nil
}

// Start starts the MCP server
func (s *Server) Start() error {
	switch s.config.Transport {
	case "stdio":
		return s.ServeStdio()
	case "streamable-http":
		return s.ServeStreamableHTTP()
	default:
		return fmt.Errorf("unsupported transport mode: %s", s.config.Transport)
	}
}

// ServeStdio serves the MCP server via stdio transport
func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.mcpServer)
}

func (s *Server) ServeStreamableHTTP() error {
	host := s.config.Host
	if host == "" {
		host = "0.0.0.0"
	}

	mux := http.NewServeMux()

	// Minimal OAuth setup for prototype
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is required")
	}

	serverURL := fmt.Sprintf("http://%s:%d", host, s.config.Port)
	oauthServer, oauthOption, err := mark3labs.WithOAuth(mux, &oauth.Config{
		Provider:  "hmac",
		Audience:  "api://gemara-mcp-server",
		JWTSecret: []byte(jwtSecret),
		ServerURL: serverURL,
	})
	if err != nil {
		return fmt.Errorf("failed to configure OAuth: %w", err)
	}

	// Create MCP server with OAuth middleware
	mcpServerWithOAuth := server.NewMCPServer(
		"gemara-mcp-server",
		s.config.Version,
		server.WithLogging(),
		oauthOption,
	)

	// Register tools on OAuth-protected server
	infoTools, err := info.NewGemaraInfoTools()
	if err != nil {
		return fmt.Errorf("failed to create info tools: %w", err)
	}
	infoTools.Register(mcpServerWithOAuth)

	authoringTools, err := authoring.NewGemaraAuthoringTools()
	if err != nil {
		return fmt.Errorf("failed to create authoring tools: %w", err)
	}
	authoringTools.Register(mcpServerWithOAuth)

	// Create StreamableHTTP server
	streamableOpts := []server.StreamableHTTPOption{
		server.WithHTTPContextFunc(oauth.CreateHTTPContextFunc()),
	}
	if s.config.Logger != nil {
		adapter := &slogLoggerAdapter{logger: s.config.Logger}
		streamableOpts = append(streamableOpts, server.WithLogger(adapter))
	}

	streamableServer := server.NewStreamableHTTPServer(mcpServerWithOAuth, streamableOpts...)
	mux.HandleFunc("/mcp", oauthServer.WrapMCPEndpoint(streamableServer))

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, s.config.Port),
		Handler: mux,
	}

	return srv.ListenAndServe()
}

// handleGemaraContextResource provides the Gemara context as a resource
func (s *Server) handleGemaraContextResource(_ context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return []mcp.ResourceContents{
		&mcp.TextResourceContents{
			URI:      request.Params.URI,
			MIMEType: "text/markdown",
			Text:     prompts.GemaraContext,
		},
	}, nil
}

// slogLoggerAdapter adapts slog.Logger to the util.Logger interface
type slogLoggerAdapter struct {
	logger *slog.Logger
}

// Infof implements util.Logger interface
func (a *slogLoggerAdapter) Infof(format string, v ...any) {
	a.logger.Info(fmt.Sprintf(format, v...))
}

// Errorf implements util.Logger interface
func (a *slogLoggerAdapter) Errorf(format string, v ...any) {
	a.logger.Error(fmt.Sprintf(format, v...))
}

