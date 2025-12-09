package tools

import (
	"context"
	"fmt"

	"github.com/complytime/gemara-mcp-server/storage"
	"github.com/complytime/gemara-mcp-server/tools/prompts"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/ossf/gemara/layer1"
	"github.com/ossf/gemara/layer2"
	"github.com/ossf/gemara/layer3"
)

// ToolKit represents a collection of MCP functionality.
type ToolKit interface {
	Name() string
	Description() string
	Register(mcpServer *server.MCPServer)
}

// GemaraAuthoringTools provides tools for creating and validating Gemara artifacts
type GemaraAuthoringTools struct {
	tools     []server.ServerTool
	prompts   []server.ServerPrompt
	resources []server.ServerResource
	// Disk-based storage with index
	storage *storage.ArtifactStorage
	// In-memory cache for quick access (populated from storage index)
	layer1Guidance map[string]*layer1.GuidanceDocument
	layer2Catalogs map[string]*layer2.Catalog
	layer3Policies map[string]*layer3.PolicyDocument
	// CUE schema cache
	schemaCache map[int]string // layer -> schema content
}

func NewGemaraAuthoringTools() (*GemaraAuthoringTools, error) {
	g := &GemaraAuthoringTools{
		layer1Guidance: make(map[string]*layer1.GuidanceDocument),
		layer2Catalogs: make(map[string]*layer2.Catalog),
		layer3Policies: make(map[string]*layer3.PolicyDocument),
		schemaCache:    make(map[int]string),
	}

	// Initialize storage
	artifactsDir := g.getArtifactsDir()
	var err error
	g.storage, err = storage.NewArtifactStorage(artifactsDir)
	if err != nil {
		// If storage initialization fails, log but continue (fallback to in-memory only)
		return g, fmt.Errorf("failed to initialize artifact storage: %w", err)
	}

	g.initTools()
	g.initPrompts()
	g.initResources()
	g.LoadArtifactsDir()

	return g, nil
}

func (g *GemaraAuthoringTools) Name() string {
	return "gemara-authoring"
}

func (g *GemaraAuthoringTools) Description() string {
	return "A set of tools related to authoring Gemara artifacts in YAML for Layers 1-4 of the Gemara model."
}

func (g *GemaraAuthoringTools) Register(s *server.MCPServer) {
	s.AddTools(g.tools...)
	s.AddPrompts(g.prompts...)
	s.AddResources(g.resources...)
}

func (g *GemaraAuthoringTools) initTools() {
	g.tools = []server.ServerTool{
		// Layer 1 Tools
		{
			Tool: mcp.NewTool(
				"list_layer1_guidance",
				mcp.WithDescription("List all available Layer 1 Guidance documents. Returns a summary of all stored guidance documents with their IDs, titles, descriptions, and metadata."),
			),
			Handler: g.handleListLayer1Guidance,
		},
		{
			Tool: mcp.NewTool(
				"get_layer1_guidance",
				mcp.WithDescription("Get detailed information about a specific Layer 1 Guidance document by its ID. Returns the full guidance document in YAML or JSON format."),
				mcp.WithString("guidance_id", mcp.Description("The unique identifier of the Layer 1 Guidance document to retrieve."), mcp.Required()),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleGetLayer1Guidance,
		},
		{
			Tool: mcp.NewTool(
				"search_layer1_guidance",
				mcp.WithDescription("Search Layer 1 Guidance documents by name, description, or author. Can optionally filter by applicability scope (boundaries, technologies, providers)."),
				mcp.WithString("search_term", mcp.Description("Search term to match against title, description, or author. Required unless scoping filters are provided.")),
				mcp.WithArray("boundaries", mcp.Description("Optional array of boundary/jurisdiction filters to apply.")),
				mcp.WithArray("technologies", mcp.Description("Optional array of technology domain filters to apply.")),
				mcp.WithArray("providers", mcp.Description("Optional array of provider/industry sector filters to apply.")),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleSearchLayer1Guidance,
		},
		{
			Tool: mcp.NewTool(
				"store_layer1_yaml",
				mcp.WithDescription("Store a Layer 1 Guidance document from raw YAML content. This preserves all YAML content without data loss. The YAML is validated with CUE before storing."),
				mcp.WithString("yaml_content", mcp.Description("Raw YAML content containing the complete Layer-1 GuidanceDocument structure. Must include metadata.id and will be validated against the Layer 1 CUE schema."), mcp.Required()),
			),
			Handler: g.handleStoreLayer1YAML,
		},
		// Layer 2 Tools
		{
			Tool: mcp.NewTool(
				"list_layer2_controls",
				mcp.WithDescription("List all available Layer 2 Controls with optional filtering by technology or Layer 1 reference. Returns controls grouped by catalog."),
				mcp.WithString("technology", mcp.Description("Optional technology filter to limit results.")),
				mcp.WithString("layer1_reference", mcp.Description("Optional Layer 1 guidance ID to filter controls that reference it.")),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleListLayer2Controls,
		},
		{
			Tool: mcp.NewTool(
				"get_layer2_control",
				mcp.WithDescription("Get detailed information about a specific Layer 2 Control by its ID. Returns the full control definition in YAML or JSON format."),
				mcp.WithString("control_id", mcp.Description("The unique identifier of the Layer 2 Control to retrieve."), mcp.Required()),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleGetLayer2Control,
		},
		{
			Tool: mcp.NewTool(
				"search_layer2_controls",
				mcp.WithDescription("Search Layer 2 Controls by name, objective, or ID. Can also filter by Layer 1 guidance reference, technology, or applicability scope."),
				mcp.WithString("search_term", mcp.Description("Search term to match against title, objective, or control ID. Required unless other filters are provided.")),
				mcp.WithString("technology", mcp.Description("Optional technology filter.")),
				mcp.WithString("layer1_reference", mcp.Description("Optional Layer 1 guidance ID to filter controls that reference it.")),
				mcp.WithArray("boundaries", mcp.Description("Optional array of boundary/jurisdiction filters to apply.")),
				mcp.WithArray("technologies", mcp.Description("Optional array of technology domain filters to apply.")),
				mcp.WithArray("providers", mcp.Description("Optional array of provider/industry sector filters to apply.")),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleSearchLayer2Controls,
		},
		{
			Tool: mcp.NewTool(
				"store_layer2_yaml",
				mcp.WithDescription("Store a Layer 2 Control Catalog from raw YAML content. This preserves all YAML content without data loss. The YAML is validated with CUE before storing."),
				mcp.WithString("yaml_content", mcp.Description("Raw YAML content containing the complete Layer-2 Catalog structure. Must include metadata.id and will be validated against the Layer 2 CUE schema."), mcp.Required()),
			),
			Handler: g.handleStoreLayer2YAML,
		},
		{
			Tool: mcp.NewTool(
				"get_layer2_guideline_mappings",
				mcp.WithDescription("Retrieve all Layer 1 guideline mappings for a Layer 2 control. Shows which Layer 1 guidance documents the control references and the specific guideline entries."),
				mcp.WithString("control_id", mcp.Description("The unique identifier of the Layer 2 Control to get mappings for."), mcp.Required()),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
				mcp.WithString("include_guidance_details", mcp.Description("Whether to include full Layer 1 guidance document details. Set to 'true' or '1' to enable.")),
			),
			Handler: g.handleGetLayer2GuidelineMappings,
		},
		// Layer 3 Tools
		{
			Tool: mcp.NewTool(
				"list_layer3_policies",
				mcp.WithDescription("List all available Layer 3 Policy documents. Returns a summary of all stored policy documents with their IDs, titles, objectives, and metadata."),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleListLayer3Policies,
		},
		{
			Tool: mcp.NewTool(
				"get_layer3_policy",
				mcp.WithDescription("Get detailed information about a specific Layer 3 Policy document by its ID. Returns the full policy document in YAML or JSON format."),
				mcp.WithString("policy_id", mcp.Description("The unique identifier of the Layer 3 Policy document to retrieve."), mcp.Required()),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleGetLayer3Policy,
		},
		{
			Tool: mcp.NewTool(
				"search_layer3_policies",
				mcp.WithDescription("Search Layer 3 Policy documents by title, objective, or other metadata."),
				mcp.WithString("search_term", mcp.Description("Search term to match against title, objective, or policy ID."), mcp.Required()),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleSearchLayer3Policies,
		},
		{
			Tool: mcp.NewTool(
				"store_layer3_yaml",
				mcp.WithDescription("Store a Layer 3 Policy document from raw YAML content. This preserves all YAML content without data loss. The YAML is validated with CUE before storing."),
				mcp.WithString("yaml_content", mcp.Description("Raw YAML content containing the complete Layer-3 PolicyDocument structure. Must include metadata.id and will be validated against the Layer 3 CUE schema."), mcp.Required()),
			),
			Handler: g.handleStoreLayer3YAML,
		},
		// Validation and Utility Tools
		{
			Tool: mcp.NewTool(
				"validate_gemara_yaml",
				mcp.WithDescription("Validate YAML content against a Gemara layer schema using CUE. Returns a detailed validation report with any errors found."),
				mcp.WithString("yaml_content", mcp.Description("Raw YAML content to validate."), mcp.Required()),
				mcp.WithNumber("layer", mcp.Description("Layer number (1-4) to validate against."), mcp.Required()),
			),
			Handler: g.handleValidateGemaraYAML,
		},
		{
			Tool: mcp.NewTool(
				"find_applicable_artifacts",
				mcp.WithDescription("Find Layer 1 and Layer 2 artifacts applicable to a given policy scope. Filters artifacts by boundaries, technologies, and providers."),
				mcp.WithArray("boundaries", mcp.Description("Optional array of boundary/jurisdiction filters.")),
				mcp.WithArray("technologies", mcp.Description("Optional array of technology domain filters.")),
				mcp.WithArray("providers", mcp.Description("Optional array of provider/industry sector filters.")),
				mcp.WithString("output_format", mcp.Description("Output format: 'yaml' (default) or 'json'.")),
			),
			Handler: g.handleFindApplicableArtifacts,
		},
		{
			Tool: mcp.NewTool(
				"get_layer_schema_info",
				mcp.WithDescription("Get information about a Gemara layer schema, including key fields, purpose, and schema location. Useful for understanding what fields are required when creating artifacts."),
				mcp.WithNumber("layer", mcp.Description("Layer number (1-4) to get schema information for."), mcp.Required()),
			),
			Handler: g.handleGetLayerSchemaInfo,
		},
	}
}

// initPrompts initializes MCP prompts for creation tasks
func (g *GemaraAuthoringTools) initPrompts() {
	// Prompts are embedded at compile time via prompts.go
	g.prompts = []server.ServerPrompt{
		{
			Prompt: mcp.NewPrompt(
				"create-layer1-guidance",
				mcp.WithPromptDescription("Guide for creating Layer 1 Guidance documents. Provides YAML structure, examples, and best practices."),
			),
			Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				return mcp.NewGetPromptResult(
					"Creating Layer 1 Guidance Documents",
					[]mcp.PromptMessage{
						mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompts.CreateLayer1Prompt)),
					},
				), nil
			},
		},
		{
			Prompt: mcp.NewPrompt(
				"create-layer2-controls",
				mcp.WithPromptDescription("Guide for creating Layer 2 Control Catalogs. Provides YAML structure, examples, and best practices."),
			),
			Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				return mcp.NewGetPromptResult(
					"Creating Layer 2 Control Catalogs",
					[]mcp.PromptMessage{
						mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompts.CreateLayer2Prompt)),
					},
				), nil
			},
		},
		{
			Prompt: mcp.NewPrompt(
				"create-layer3-policies",
				mcp.WithPromptDescription("Guide for creating Layer 3 Policy documents. Provides YAML structure, examples, and best practices."),
			),
			Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				return mcp.NewGetPromptResult(
					"Creating Layer 3 Policy Documents",
					[]mcp.PromptMessage{
						mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompts.CreateLayer3Prompt)),
					},
				), nil
			},
		},
		{
			Prompt: mcp.NewPrompt(
				"gemara-quick-start",
				mcp.WithPromptDescription("Quick start guide for creating your first Gemara artifacts. Provides step-by-step instructions and common workflows."),
			),
			Handler: func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
				return mcp.NewGetPromptResult(
					"Gemara Quick Start Guide",
					[]mcp.PromptMessage{
						mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompts.QuickStartPrompt)),
					},
				), nil
			},
		},
	}
}

// initResources initializes MCP resources for CUE schemas
func (g *GemaraAuthoringTools) initResources() {
	g.resources = []server.ServerResource{
		// Common schemas
		{
			Resource: mcp.NewResource(
				"gemara://schema/common/base",
				"Base Schema",
				mcp.WithResourceDescription("Common base CUE schema used by all Gemara layers"),
				mcp.WithMIMEType("text/x-cue"),
			),
			Handler: g.handleBaseSchemaResource,
		},
		{
			Resource: mcp.NewResource(
				"gemara://schema/common/metadata",
				"Metadata Schema",
				mcp.WithResourceDescription("Common metadata CUE schema used by all Gemara layers"),
				mcp.WithMIMEType("text/x-cue"),
			),
			Handler: g.handleMetadataSchemaResource,
		},
		{
			Resource: mcp.NewResource(
				"gemara://schema/common/mapping",
				"Mapping Schema",
				mcp.WithResourceDescription("Common mapping CUE schema used by all Gemara layers"),
				mcp.WithMIMEType("text/x-cue"),
			),
			Handler: g.handleMappingSchemaResource,
		},
		// Layer-specific schemas
		{
			Resource: mcp.NewResource(
				"gemara://schema/layer/1",
				"Layer 1 Schema",
				mcp.WithResourceDescription("CUE schema for Layer 1 Guidance documents"),
				mcp.WithMIMEType("text/x-cue"),
			),
			Handler: g.handleLayer1SchemaResource,
		},
		{
			Resource: mcp.NewResource(
				"gemara://schema/layer/2",
				"Layer 2 Schema",
				mcp.WithResourceDescription("CUE schema for Layer 2 Control Catalogs"),
				mcp.WithMIMEType("text/x-cue"),
			),
			Handler: g.handleLayer2SchemaResource,
		},
		{
			Resource: mcp.NewResource(
				"gemara://schema/layer/3",
				"Layer 3 Schema",
				mcp.WithResourceDescription("CUE schema for Layer 3 Policy documents"),
				mcp.WithMIMEType("text/x-cue"),
			),
			Handler: g.handleLayer3SchemaResource,
		},
		{
			Resource: mcp.NewResource(
				"gemara://schema/layer/4",
				"Layer 4 Schema",
				mcp.WithResourceDescription("CUE schema for Layer 4 Evaluation documents"),
				mcp.WithMIMEType("text/x-cue"),
			),
			Handler: g.handleLayer4SchemaResource,
		},
	}
}
