package promptsets

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// MCPPromptServer interface for registering prompts with mark3labs MCP server
type MCPPromptServer interface {
	AddPrompt(prompt mcp.Prompt, handler server.PromptHandlerFunc)
}

// RegisterPromptsWithMark3LabsServer registers prompts with a mark3labs MCP server
func RegisterPromptsWithMark3LabsServer(mcpServer MCPPromptServer, group *PromptSetGroup) error {
	// Get all prompts from the group
	for _, ps := range group.ListPromptSets() {
		for _, prompt := range ps.ListPrompts() {
			// Create MCP prompt with description
			promptOpts := []mcp.PromptOption{
				mcp.WithPromptDescription(prompt.Description),
			}

			// Add arguments for each variable in the prompt
			for varName, varDesc := range prompt.Variables {
				promptOpts = append(promptOpts, mcp.WithArgument(
					varName,
					mcp.ArgumentDescription(varDesc),
					mcp.RequiredArgument(),
				))
			}

			mcpPrompt := mcp.NewPrompt(prompt.Name, promptOpts...)

			// Create handler that uses our prompt generation
			handler := func(ctx context.Context, req mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				// Extract variables from request params
				variables := make(map[string]interface{})

				// Parse arguments if provided in params
				// Arguments is map[string]string in mark3labs
				if req.Params.Arguments != nil {
					for key, value := range req.Params.Arguments {
						variables[key] = value
					}
				}

				// Generate the prompt using our system
				reqInternal := PromptRequest{
					PromptName: prompt.Name,
					Variables:  variables,
				}

				resp, err := ps.GeneratePrompt(ctx, reqInternal)
				if err != nil {
					return nil, fmt.Errorf("failed to generate prompt: %w", err)
				}

				// Convert to MCP prompt result
				return mcp.NewGetPromptResult(
					prompt.Name,
					[]mcp.PromptMessage{
						mcp.NewPromptMessage(
							mcp.RoleUser,
							mcp.NewTextContent(resp.Content),
						),
					},
				), nil
			}

			// Register the prompt
			mcpServer.AddPrompt(mcpPrompt, handler)
		}
	}

	return nil
}

// RegisterGemaraPromptsWithMark3LabsServer registers Gemara-specific prompts
func RegisterGemaraPromptsWithMark3LabsServer(mcpServer MCPPromptServer) error {
	group := GetGemaraPromptSets()
	return RegisterPromptsWithMark3LabsServer(mcpServer, group)
}

// RegisterUserFacingPromptsWithMark3LabsServer registers user-facing prompts
func RegisterUserFacingPromptsWithMark3LabsServer(mcpServer MCPPromptServer) error {
	group := GetUserFacingPromptSets()
	return RegisterPromptsWithMark3LabsServer(mcpServer, group)
}

// RegisterAllPromptsWithMark3LabsServer registers all prompts (Gemara + user-facing)
func RegisterAllPromptsWithMark3LabsServer(mcpServer MCPPromptServer) error {
	// Register Gemara prompts
	if err := RegisterGemaraPromptsWithMark3LabsServer(mcpServer); err != nil {
		return fmt.Errorf("failed to register Gemara prompts: %w", err)
	}

	// Register user-facing prompts
	if err := RegisterUserFacingPromptsWithMark3LabsServer(mcpServer); err != nil {
		return fmt.Errorf("failed to register user-facing prompts: %w", err)
	}

	return nil
}
