# Implementation Plan: Gemara Artifact Authoring Assistant

**Branch**: `001-gemara-authoring-assistant` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-gemara-authoring-assistant/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build a high-performance Go MCP Server that enables LLMs to automate the synthesis, mapping, and validation of Gemara security artifacts with deterministic precision. The server acts as a "Compliance Orchestrator" supporting three critical user journeys: Auto-Documentation (generate Layer 2 Catalog artifacts from technical evidence via 5-phase sequential pipeline), Inheritance Discovery (identify inherited controls), and Framework Pivot (assess technology against regulatory standards). The implementation leverages CUE for schema validation of all three Definition layers (Layer 1 GuidanceDocument, Layer 2 Catalog, Layer 3 Policy), provides config parsing interfaces, and operates statelessly with dual transport support (stdio for local IDEs, HTTP for cloud deployment). Layer 3 Policy generation is supported through scope definition with Layer 1 and Layer 2 applicability queried for context. Control IDs follow format `<identifier>-<numbering>` and are immutable once defined.

## Technical Context

**Language/Version**: Go 1.24.0  
**Primary Dependencies**: 
- `cuelang.org/go` - CUE schema validation engine for all three Definition layers
- `github.com/goccy/go-yaml` - YAML parsing for Gemara artifacts
- `github.com/mark3labs/mcp-go` - MCP protocol implementation
- `github.com/gemaraproj/go-gemara` - Gemara framework integration
- `github.com/spf13/cobra` - CLI command structure
- `github.com/stretchr/testify` - Testing framework

**Storage**: Stateless operation (NFR-001) - no persistence between requests. Gemara information storage uses interface-based design with lightweight file-based default implementation (FR-006). Query operations are request-scoped. Layer 1 (GuidanceDocument) artifacts are reference-only (query/storage), not authored by this system.

**Testing**: `github.com/stretchr/testify` for unit and integration tests. Contract tests for MCP tool interfaces. Integration tests for CUE validation workflows across all three Definition layers.

**Target Platform**: Linux server (MCP server). Supports dual transport: stdio for local IDE integration (Cursor/VS Code) and HTTP for cloud-native deployment.

**Project Type**: Single project - MCP server extension to existing gemara-mcp-server codebase

**Performance Goals**: 
- Generate validated Layer 2 Catalog artifacts in under 10 minutes (SC-001)
- Generate Layer 3 Policy artifacts through scope definition in reasonable time
- 90% deterministic artifact generation (NFR-003, SC-008)
- Framework pivot gap analysis in under 15 minutes (SC-005)
- High-performance request processing with stateless operation

**Constraints**: 
- Stateless operation - no data persistence between requests (NFR-001)
- Input-output purity - validate structure without modifying LLM content (NFR-009)
- Dual transport support required (local stdio + remote HTTP)
- 90% deterministic outcomes for artifact generation
- Encrypted remote communications (NFR-004)
- Secure authentication required (NFR-005)
- Scope limited to Definition layers (Layers 1-3); Measurement layers (Layers 4-6) explicitly out of scope
- Control IDs MUST use format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042") and MUST NOT include family. Control IDs are immutable once defined until withdrawn

**Scale/Scope**: 
- MCP server handling LLM-driven artifact authoring requests
- Support for Layer 1 (GuidanceDocument - reference only), Layer 2 (Catalog - full authoring), and Layer 3 (Policy - generation via scope definition) Gemara artifacts
- Validation support for all three Definition layers
- Integration with external MCP servers (GitHub/GitLab) for technical evidence access
- Config parsing interface for extensible evidence processing

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Dependency Management ✅
- **Status**: PASS
- **Check**: Using latest stable versions from provided dependency list. Dependencies are well-maintained (CUE, go-yaml, mcp-go, go-gemara). Versions will be pinned in go.mod.
- **Rationale**: All dependencies are stable, actively maintained projects with good adoption rates.

### II. Code Style Standards ✅
- **Status**: PASS
- **Check**: Go code follows project conventions (lowercase_with_underscores file names, short package names). License headers required. `.golangci.yml` checks enforced.
- **Rationale**: Existing codebase already follows Go conventions. New code will maintain consistency.

### III. Centralized Constants (NON-NEGOTIABLE) ✅
- **Status**: PASS
- **Check**: Constants will be placed in `internal/consts/consts.go` (already exists). No magic strings/numbers inline. Control ID format pattern will be centralized.
- **Rationale**: Project already has constants structure. New constants will follow existing pattern.

### IV. Required Questions Before Implementation ✅
- **Status**: PASS
- **Check**: Design questions addressed in spec (user personas: compliance engineers, problem: automate artifact authoring). No data storage persistence (stateless operation).
- **Rationale**: Specification clearly defines user needs and system behavior. Stateless design eliminates data lifecycle concerns.

### V. Testing Requirements ✅
- **Status**: PASS
- **Check**: Tests required for all code changes. Using testify framework. Contract tests for MCP tools, integration tests for validation workflows.
- **Rationale**: Testing framework established. Test coverage will be maintained.

### VI. PR Workflow Standards ✅
- **Status**: PASS
- **Check**: Feature branch created (`001-gemara-authoring-assistant`). PRs will follow conventional commits format.
- **Rationale**: Standard workflow applies. No exceptions needed.

### VII. Design Documentation ✅
- **Status**: PASS
- **Check**: Design decisions documented in plan.md, research.md, data-model.md. Architecture decisions will be documented.
- **Rationale**: Planning phase includes design documentation. Decisions will be captured.

### VIII. Incremental Improvement ✅
- **Status**: PASS
- **Check**: Feature is incremental addition to existing MCP server. Interface-based design allows extensibility.
- **Rationale**: Extends existing functionality without breaking changes. Interface patterns enable future improvements.

**Overall Status**: ✅ **PASS** - All constitution gates pass. Proceeding to Phase 0 research.

## Project Structure

### Documentation (this feature)

```text
specs/001-gemara-authoring-assistant/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md             # Feature specification
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── mcp_tools.yaml   # MCP tool definitions
│   └── api_schema.yaml  # API contract schema
├── checklists/
│   └── requirements.md  # Specification quality checklist
├── ANALYSIS_REPORT.md   # Cross-artifact consistency analysis
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/gemara-mcp-server/
└── main.go              # Server entry point

mcp/
├── server.go            # MCP server implementation
└── server_test.go       # Server tests

tools/
├── authoring/           # Authoring tools (existing + new)
│   ├── layer1.go
│   ├── layer2.go
│   ├── layer3.go        # NEW: Layer 3 Policy generation
│   ├── register_tools.go
│   ├── parse_evidence.go        # NEW: parse_technical_evidence tool
│   ├── query_threats.go         # NEW: query_threat_library tool
│   ├── validate_artifact.go     # NEW: validate_gemara_artifact tool (all 3 layers)
│   ├── generate_artifact.go     # NEW: generate_layer2_artifact tool
│   ├── generate_policy.go        # NEW: generate_layer3_policy tool (scope-based)
│   ├── search_inheritance.go    # NEW: search_inheritance_opportunities tool
│   ├── import_controls.go       # NEW: import_inherited_controls tool (uses imported-controls field)
│   ├── analyze_pivot.go          # NEW: analyze_framework_pivot tool
│   ├── error_handler.go         # NEW: Error handling
│   ├── logging.go               # NEW: Logging
│   ├── compare_controls.go      # NEW: Control comparison logic
│   ├── gap_analysis.go          # NEW: Gap analysis engine
│   ├── prioritize.go            # NEW: Prioritization logic
│   ├── recommendations.go       # NEW: Recommendation generator
│   ├── confidence.go            # NEW: Confidence indicators
│   └── types.go                 # NEW: Shared types
├── info/                # Info tools (existing)
│   ├── gemara_info.go
│   ├── validation.go
│   └── ...
└── prompts/             # Prompt templates (existing)
    ├── create-layer2.md
    └── ...

internal/
├── consts/
│   └── consts.go       # Centralized constants (including control ID format pattern)
├── parsing/             # NEW: Config parsing interfaces
│   ├── interface.go    # Parser interface
│   ├── file_based.go    # File-based parser implementation
│   ├── dependency_parser.go  # NEW: Dependency info parser (SBOMs, CALM)
│   ├── regulatory_parser.go # NEW: Regulatory requirement parser
│   ├── types.go        # NEW: Parsing-related types
│   └── parsers_test.go
├── validation/          # NEW: Enhanced validation (all 3 layers)
│   ├── cue_validator.go # CUE-based schema validator
│   ├── layer1_validator.go # NEW: Layer 1 GuidanceDocument validator
│   ├── layer2_validator.go # NEW: Layer 2 Catalog validator
│   ├── layer3_validator.go # NEW: Layer 3 Policy validator
│   ├── version_manager.go # Gemara spec version management
│   ├── schemas.go      # CUE schema loading (all 3 layers)
│   └── validator_test.go
├── storage/             # NEW: Gemara info storage interface
│   ├── interface.go     # Storage interface
│   ├── file_storage.go  # File-based storage implementation
│   ├── ranking.go       # NEW: Relevance ranking algorithm
│   ├── types.go         # NEW: Storage-related types
│   └── storage_test.go
├── metrics/             # NEW: Observability
│   ├── metrics.go       # Metrics definitions
│   └── prometheus.go    # Prometheus integration
└── errors.go            # NEW: Error handling infrastructure

storage/                 # Existing storage (may be extended)
├── interface.go
└── storage.go

tests/
├── contract/            # NEW: Contract tests for MCP tools
│   ├── test_parse_technical_evidence.go
│   ├── test_query_threat_library.go
│   ├── test_validate_gemara_artifact.go
│   ├── test_generate_layer2_artifact.go
│   ├── test_generate_layer3_policy.go
│   ├── test_search_inheritance.go
│   ├── test_import_controls.go
│   └── test_framework_pivot.go
├── integration/         # NEW: Integration tests
│   ├── test_auto_documentation.go
│   ├── test_inheritance_discovery.go
│   ├── test_framework_pivot.go
│   └── test_layer3_policy_generation.go
└── unit/                # Unit tests for new components
    ├── parsing_test.go
    ├── validation_test.go
    ├── storage_test.go
    └── gap_analysis_test.go
```

**Structure Decision**: Extending existing single-project structure. New components organized under `internal/` following Go conventions. Authoring capabilities integrated into existing `tools/authoring/` directory. Validation expanded to support all three Definition layers. Layer 3 Policy generation tool added. Import functionality uses Gemara `imported-controls` field with MultiMapping structure. Control ID format validation and constants centralized. Test structure expanded to include contract and integration tests for all journeys and Layer 3 generation.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all constitution gates pass. Interface-based design enables extensibility without violating simplicity principles.

---

## Phase 0: Research ✅

**Status**: Complete  
**Output**: `research.md`

All research tasks completed. Key decisions documented:
- MCP tool-based architecture with stateless operation
- CUE validation for deterministic structural correctness (all three Definition layers)
- Interface-based design for parsing and storage
- Structured data flows for each journey (5-phase pipeline for Auto-Documentation)
- Dual transport support (stdio + HTTP)
- Prometheus metrics for observability
- Gemara `imported-controls` field with MultiMapping structure for control inheritance
- Layer 3 Policy generation via scope definition with Layer 1 and Layer 2 context
- Control ID format: `<identifier>-<numbering>` (immutable, family-independent)

No blocking clarifications remain.

---

## Phase 1: Design & Contracts ✅

**Status**: Complete  
**Outputs**: `data-model.md`, `contracts/mcp_tools.yaml`, `quickstart.md`

### Data Model
- 10 core entities defined (TechnicalEvidence, SecurityFeature, Threat, Control, Layer2Artifact, Layer2Catalog, Layer1Guidance, Layer3Policy, RegulatoryRequirement, GapAnalysisReport)
- Supporting types and relationships documented
- State transitions and validation rules specified
- Import mechanism uses Gemara `imported-controls` field (`[...#MultiMapping]` structure) with `@go(ImportedControls)` tag
- Layer 3 Policy entity added with scope definition and import relationships
- MultiMapping structure documented for imported-controls, imported-threats, imported-capabilities
- Control ID format constraint: `<identifier>-<numbering>`, immutable, family-independent

### Contracts
- 8 MCP tools defined with input/output schemas:
  - `parse_technical_evidence` - Config parsing (FR-002, FR-003)
  - `query_threat_library` - Threat mapping (FR-006)
  - `validate_gemara_artifact` - Schema validation for all three Definition layers (FR-004, FR-005)
  - `query_gemara_info` - Info storage queries (FR-006)
  - `search_inheritance_opportunities` - Inheritance discovery (FR-009, FR-010)
  - `import_inherited_controls` - Control import via Gemara `imported-controls` field (FR-011)
  - `analyze_framework_pivot` - Gap analysis (FR-012, FR-013, FR-014)
  - `generate_layer2_artifact` - Auto-documentation orchestration (FR-001 through FR-008)
  - `generate_layer3_policy` - Layer 3 Policy generation via scope definition (FR-006a)

### Quick Start Guide
- Step-by-step examples for all three journeys
- Common patterns and error handling
- Best practices documented
- Import workflow using Gemara `imported-controls` field with MultiMapping structure
- Layer 3 Policy generation workflow
- Control ID format examples and validation

### Agent Context
- Updated Cursor IDE context file with Go 1.24.0 and project structure

---

## Phase 2: Task Breakdown

**Status**: Complete  
**Output**: `tasks.md` (generated by `/speckit.tasks`)

Task breakdown completed with tasks organized by user story. All three journeys (Auto-Documentation, Inheritance Discovery, Framework Pivot) have complete task coverage. Layer 3 Policy generation tasks included. Control ID format validation tasks included.

---

## Post-Design Constitution Re-Check ✅

**Status**: PASS

All constitution principles remain satisfied after design phase:
- Dependency management: All dependencies are stable versions
- Code style: Go conventions maintained
- Constants: Centralized in internal/consts/consts.go (including control ID format pattern)
- Testing: Comprehensive test coverage planned for all three Definition layers
- Design documentation: Complete (plan.md, research.md, data-model.md, contracts/)
- Incremental improvement: Extends existing codebase without breaking changes

**Overall Status**: ✅ **PASS** - Ready for implementation
