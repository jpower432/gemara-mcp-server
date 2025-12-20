# CURSOR.md - Development Guide

This guide provides step-by-step instructions for running, testing, and developing the Gemara MCP Server.

## Prerequisites

- **Go 1.24+** - [Install Go](https://go.dev/doc/install)
- **Make** - Usually pre-installed on Linux/macOS
- **Podman or Docker** (optional) - For containerized builds
- **Cursor IDE** (optional) - For MCP client testing

## Quick Start

### 1. Clone and Navigate

```bash
cd gemara-mcp-server
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Build the Binary

```bash
make build
# Binary will be in bin/gemara-mcp-server
```

Or build directly:

```bash
go build -o bin/gemara-mcp-server ./cmd/gemara-mcp-server
```

### 4. Run the Server

**Local Development (Stdio Transport):**
```bash
./bin/gemara-mcp-server
# or with debug logging
./bin/gemara-mcp-server --debug
```

**Remote/Sandboxed Environments (StreamableHTTP via Container):**
See [Container Development](#container-development) section below.

## Development Workflow

### Building

```bash
# Build binary
make build

# Build and install to ~/.local/bin
make install

# Build container image
make container-build
```

### Running Locally

**Stdio Transport:**
```bash
# Basic run
./bin/gemara-mcp-server

# With debug logging
./bin/gemara-mcp-server --debug

# Explicit stdio transport
./bin/gemara-mcp-server --transport stdio --debug
```

**Note:** For remote or sandboxed environments, use StreamableHTTP transport via containers (see [Container Development](#container-development) section).

### Testing

```bash
# Run all tests
make test

# Run tests with verbose output
go test -v ./...

# Run tests for specific package
go test ./tools/...

# Run tests with coverage
go test -cover ./...
```

### Code Quality

```bash
# Format code
go fmt ./...

# Run goimports (if installed)
goimports -w .

# Check for issues
go vet ./...

# Build to check for compilation errors
go build ./...
```

## Running in Cursor IDE

### Local Development (Stdio Transport)

For local development, use stdio transport:

1. **Update `.cursor/mcp.json`:**
   ```json
   {
     "mcpServers": {
       "gemara-mcp-server": {
         "command": "/absolute/path/to/bin/gemara-mcp-server",
         "args": ["--debug"]
       }
     }
   }
   ```

2. **Restart Cursor IDE** to pick up the MCP server configuration

3. **Test in Cursor:**
   - Open Cursor chat
   - The MCP server tools should be available
   - Try: "List all Layer 1 guidance documents"

### Remote/Sandboxed Environments (StreamableHTTP)

For remote or sandboxed environments, use StreamableHTTP transport via containers. See [Container Development](#container-development) section for setup instructions.

## Container Development

**Use containers for remote or sandboxed environments that require StreamableHTTP transport.**

### Build Container

```bash
# Build image
make container-build
# or
podman build -t gemara-mcp-server:latest -f Containerfile .
```

### Run Container (StreamableHTTP)

**Default (Writable Artifacts):**
```bash
# Ensure artifacts directory exists
mkdir -p artifacts
chmod 755 artifacts

# Set JWT_SECRET for OAuth HMAC provider (required for HTTP transport)
# This secret is used to validate JWT tokens (HMAC-SHA256)
# For production, use Okta/Google/Azure providers instead
export JWT_SECRET="your-32-byte-secret-key-minimum-length-required"

# Run with StreamableHTTP transport (allows storing new artifacts)
# Using --host=0.0.0.0 for container networking (Podman/Docker)
make container-run
# or
podman run --rm --userns=keep-id -p 8080:8080 \
  -v "$(pwd)/artifacts:/app/artifacts:z" \
  --user $(id -u):$(id -g) \
  -e JWT_SECRET="${JWT_SECRET}" \
  gemara-mcp-server:latest
```

**Read-Only Mode (Query Only):**
```bash
# Ensure artifacts directory exists
mkdir -p artifacts

# Set JWT_SECRET for OAuth HMAC provider
export JWT_SECRET="your-32-byte-secret-key-minimum-length-required"

# Run with read-only artifacts (cannot store new artifacts, query only)
make container-run-readonly
# or
podman run --rm --userns=keep-id -p 8080:8080 \
  -v "$(pwd)/artifacts:/app/artifacts:z,ro" \
  --user $(id -u):$(id -g) \
  -e JWT_SECRET="${JWT_SECRET}" \
  gemara-mcp-server:latest
```

**Note:** 
- **HTTP transport avoids Podman mount complexity**: STDIO transport with Podman requires complex user namespace and mount configurations that often fail. HTTP transport only needs port forwarding.
- The `--userns=keep-id` flag ensures the container user matches your host user ID, preventing permission issues when writing to the mounted artifacts directory.
- The `:z` flag (lowercase) sets a shared SELinux context, which is less restrictive than `:Z` (private context).
- Ensure the `artifacts` directory exists and is writable before running the container.
- **Security**: The server uses `--host=0.0.0.0` to bind to all interfaces for container networking. OAuth 2.1 middleware (oauth-mcp-proxy) protects access - ensure `JWT_SECRET` is set for HMAC provider. Clients must send valid JWT tokens signed with this secret.
- For local development (non-container), use STDIO transport (default) or `--host=127.0.0.1` for HTTP.

The server will be accessible at `http://localhost:8080/mcp` for StreamableHTTP connections.

**Note:** In read-only mode, tools that store artifacts (`store_layer1_yaml`, `store_layer2_yaml`, `store_layer3_yaml`) will fail with a storage error.

### Configure Cursor for StreamableHTTP

Update `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "gemara-mcp-server": {
      "url": "http://localhost:8080/mcp",
      "headers": {
        "Authorization": "Bearer <JWT_TOKEN>"
      }
    }
  }
}
```

**OAuth 2.1 Authentication:**

The server uses `oauth-mcp-proxy` library with HMAC provider for OAuth 2.1 Bearer token validation. Clients must send JWT tokens signed with HMAC-SHA256 using the `JWT_SECRET` environment variable.

**JWT Token Requirements:**
- Signed with HMAC-SHA256 using `JWT_SECRET`
- Must include `aud` (audience) claim: `"api://gemara-mcp-server"`
- Must include `sub` (subject) claim for user identification
- Optional: `email`, `preferred_username` claims

**Example JWT Generation (for testing):**
```go
import "github.com/golang-jwt/jwt/v5"

token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "sub":   "user-123",
    "email": "user@example.com",
    "preferred_username": "john.doe",
    "aud":   "api://gemara-mcp-server",
    "exp":   time.Now().Add(time.Hour).Unix(),
    "iat":   time.Now().Unix(),
})
tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
```

**Note**: If `JWT_SECRET` is not set, the server will reject requests with 401 Unauthorized (RFC 6750 compliant). For production deployments, consider using Okta, Google, or Azure OIDC providers instead of HMAC.

**Security Requirements (per oauth-mcp-proxy guidelines):**
- `JWT_SECRET` must be at least 32 bytes (256 bits) for security
- Generate securely: `openssl rand -base64 32`
- Never commit secrets to git - use environment variables only
- Use HTTPS in production (HTTP is acceptable for local development only)
- Tokens are never logged - only token hashes (SHA-256) are logged for debugging

## Security Checklist

This server follows security best practices from [oauth-mcp-proxy security guidelines](https://github.com/tuannvm/oauth-mcp-proxy/blob/main/docs/SECURITY.md).

### ✅ Implemented Security Features

- **Secrets Management**: All secrets loaded from environment variables (never hardcoded)
- **JWT Secret Validation**: Minimum 32-byte length enforced at startup
- **OAuth 2.1 Compliance**: Uses oauth-mcp-proxy library for RFC 6750 compliant authentication
- **Audience Validation**: Tokens must include `aud: "api://gemara-mcp-server"` claim
- **Token Caching**: 5-minute cache for validated tokens (handled by library)
- **Secure Logging**: Only token hashes (SHA-256) logged, never full tokens
- **Custom Logger**: Integrated with slog for production logging
- **Localhost Binding**: Defaults to `127.0.0.1` to prevent NeighborJacking attacks

### ⚠️ Production Considerations

**Before deploying to production:**

- [ ] **HTTPS Required**: Configure TLS certificates (Let's Encrypt, AWS ACM, etc.)
- [ ] **OIDC Provider**: Migrate from HMAC to Okta/Google/Azure for production
- [ ] **Rate Limiting**: Add rate limiting to OAuth endpoints (recommended)
- [ ] **Security Headers**: Verify HSTS, CSP headers are configured
- [ ] **Secret Rotation**: Plan for 90-day secret rotation schedule
- [ ] **Monitoring**: Set up alerts for authentication failures
- [ ] **Audit Logging**: Review OAuth provider audit logs regularly

**Current Prototype Limitations:**
- Uses HTTP (not HTTPS) - acceptable for local development only
- HMAC provider requires manual JWT generation - not suitable for user authentication
- No rate limiting on OAuth endpoints (library handles validation, but doesn't rate limit)

See [oauth-mcp-proxy Security Guide](https://github.com/tuannvm/oauth-mcp-proxy/blob/main/docs/SECURITY.md) for complete security best practices.

### Clean Up

```bash
make container-clean
# or
podman rmi gemara-mcp-server:latest
```

## Artifacts Directory

The server automatically looks for artifacts in an `artifacts/` directory:

```
artifacts/
├── layer1/
│   └── *.yaml
├── layer2/
│   └── *.yaml
├── layer3/
│   └── *.yaml
└── layer4/
    └── *.yaml
```

**Create the directory structure:**
```bash
mkdir -p artifacts/{layer1,layer2,layer3,layer4}
```

**Note:** The server will create these directories automatically if they don't exist.

## Common Development Tasks

### Adding a New Tool

1. **Create handler function** in appropriate file (e.g., `tools/layer1.go`):
   ```go
   func (g *GemaraAuthoringTools) handleNewTool(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
       // Implementation
   }
   ```

2. **Add tool definition** in `tools/register_tools.go`:
   ```go
   func (g *GemaraAuthoringTools) newNewTool() server.ServerTool {
       return server.ServerTool{
           Tool: mcp.NewTool(
               "new_tool_name",
               mcp.WithDescription("Tool description"),
               // Parameters...
           ),
           Handler: g.handleNewTool,
       }
   }
   ```

3. **Register tool** in `registerTools()`:
   ```go
   tools = append(tools, g.newNewTool())
   ```

4. **Rebuild and test:**
   ```bash
   make build
   ./bin/gemara-mcp-server --debug
   ```

### Adding a New Prompt

1. **Create prompt file** in `tools/prompts/` (e.g., `new-prompt.md`)

2. **Embed prompt** in `tools/prompts/prompts.go`:
   ```go
   //go:embed new-prompt.md
   var NewPrompt string
   ```

3. **Add prompt definition** in `tools/register_prompts.go`:
   ```go
   func (g *GemaraAuthoringTools) newNewPrompt() server.ServerPrompt {
       return server.ServerPrompt{
           Prompt: mcp.NewPrompt(
               "new-prompt",
               mcp.WithPromptDescription("Prompt description"),
           ),
           Handler: func(_ context.Context, _ mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
               return mcp.NewGetPromptResult(
                   "Prompt Title",
                   []mcp.PromptMessage{
                       mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompts.NewPrompt)),
                   },
               ), nil
           },
       }
   }
   ```

4. **Register prompt** in `registerPrompts()`

### Adding a New Resource

1. **Add resource handler** in `tools/resources.go`:
   ```go
   func (g *GemaraAuthoringTools) handleNewResource(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
       // Implementation
   }
   ```

2. **Add resource definition** in `tools/register_resources.go`:
   ```go
   func (g *GemaraAuthoringTools) newNewResource() server.ServerResource {
       return server.ServerResource{
           Resource: mcp.NewResource(
               "gemara://resource/path",
               "Resource Name",
               mcp.WithResourceDescription("Description"),
               mcp.WithMIMEType("text/plain"),
           ),
           Handler: g.handleNewResource,
       }
   }
   ```

3. **Register resource** in `registerResources()`

## Debugging

### Enable Debug Logging

```bash
./bin/gemara-mcp-server --debug
```

### Check Logs

Debug logs are written to `stderr`. For stdio transport, logs may be mixed with protocol messages.

For StreamableHTTP, logs appear in the terminal:

```bash
./bin/gemara-mcp-server --transport streamable --port 8080 --debug
```

### Common Issues

**Issue: "Storage not available"**
- **Solution:** Ensure the artifacts directory exists and is writable
- **Check:** `ls -la artifacts/`

**Issue: "Failed to initialize artifact storage"**
- **Solution:** Check directory permissions
- **Fix:** `chmod 755 artifacts/`

**Issue: "Port already in use" (Container)**
- **Solution:** Use a different port or stop the existing container
- **Fix:** Change port mapping: `podman run -p 8081:8080 ...`

**Issue: Podman mount/permission errors with STDIO transport**
- **Solution:** Use HTTP transport (`--transport=streamable-http`) instead of STDIO for container deployments
- **Why:** Podman's user namespace mapping (`--userns=keep-id`) and volume mounts can cause permission issues with STDIO transport
- **Fix:** HTTP transport avoids mount complexity - only requires port forwarding

**Issue: "CUE validation failed"**
- **Solution:** Check YAML syntax and schema compliance
- **Debug:** Use `validate_gemara_yaml` tool first

**Issue: MCP server not appearing in Cursor**
- **Solution:** 
  1. For stdio: Verify binary path in `.cursor/mcp.json` is absolute
  2. For StreamableHTTP: Verify server is running: `curl http://localhost:8080/mcp`
  3. Check `.cursor/mcp.json` configuration matches your setup
  4. Restart Cursor IDE
  5. Check Cursor logs for MCP connection errors

## Version Information

```bash
./bin/gemara-mcp-server version
```