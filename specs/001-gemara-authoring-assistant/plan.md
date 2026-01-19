# Implementation Plan: Gemara Artifact Authoring Assistant

**Branch**: `001-gemara-authoring-assistant` | **Date**: 2025-01-27 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-gemara-authoring-assistant/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

MCP Server for authoring and validating Gemara Definition artifacts (Layers 1-3). Supports three journeys: (1) Auto-Documentation - generate Layer 2 artifacts from technical evidence via LLM-driven pipeline, (2) Inheritance Discovery - find and import existing Layer 2 catalogs, (3) Framework Pivot - analyze gaps between existing controls and new regulatory requirements. Uses CUE schema validation, file-based storage, and MCP protocol for stateless operation.

## Technical Context

**Language/Version**: Go 1.24.0  
**Primary Dependencies**: 
- `github.com/mark3labs/mcp-go` v0.43.2 (MCP server framework)
- `cuelang.org/go` v0.15.1 (CUE schema validation)
- `github.com/gemaraproj/go-gemara` (Gemara schema types and definitions)
- `github.com/spf13/cobra` v1.10.1 (CLI framework)
- `github.com/stretchr/testify` v1.11.1 (testing)
- `go.opentelemetry.io/otel` (observability - primary metrics API)
- `go.opentelemetry.io/otel/exporters/prometheus` (optional Prometheus export)

**Storage**: File-based storage (`internal/storage/file_storage.go`) for Gemara artifacts (Layer 1, Layer 2, Layer 3). No database required - stateless operation per NFR-001.

**Testing**: `github.com/stretchr/testify` for unit, integration, and contract tests. Test structure: `tests/unit/`, `tests/integration/`, `tests/contract/`.

**Target Platform**: Linux server (MCP server via stdio or HTTP transport). Supports local development and remote cloud deployment per NFR-002.

**Project Type**: Single project (CLI/server application). Structure: `cmd/gemara-mcp-server/` (CLI), `mcp/` (MCP server), `tools/` (MCP tool handlers), `internal/` (core logic).

**Performance Goals**: 
- Layer 2 artifact generation completes in under 10 minutes (SC-001)
- Framework pivot gap analysis reports in under 15 minutes (SC-005)
- 90% deterministic outcomes for artifact generation (NFR-003, SC-008)

**Constraints**: 
- Stateless operation - no data persistence between requests (NFR-001)
- 90% deterministic outcomes when processing same input (NFR-003)
- All remote communications encrypted (NFR-004)
- User authentication required (NFR-005)
- Logical isolation of request-scoped data between sessions (NFR-006)

**Scale/Scope**: 
- MCP server handling multiple concurrent requests
- Support for querying and storing Layer 1, Layer 2, Layer 3 artifacts
- Processing technical evidence from Git repositories via external MCP servers
- Generating and validating Gemara artifacts conforming to official schemas

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### I. Dependency Management ✅
- **Status**: PASS
- **Check**: All dependencies use latest stable versions, pinned in go.mod
- **Rationale**: Dependencies are well-maintained, production-ready versions

### II. Code Style Standards ✅
- **Status**: PASS
- **Check**: Go code follows style guidelines (lowercase files, package naming, error handling, gofmt/goimports)
- **Rationale**: Existing codebase demonstrates adherence to Go conventions

### III. Centralized Constants (NON-NEGOTIABLE) ✅
- **Status**: PASS
- **Check**: Constants centralized in `internal/consts/consts.go`
- **Rationale**: Single source of truth for configuration values

### IV. Required Questions Before Implementation ✅
- **Status**: PASS
- **Check**: User personas and problems clearly defined in spec (compliance engineers, security teams)
- **Rationale**: Specification includes user stories with clear personas and problems

### V. Testing Requirements ✅
- **Status**: PASS
- **Check**: Test structure exists (`tests/unit/`, `tests/integration/`, `tests/contract/`), testify used
- **Rationale**: Testing infrastructure in place, TDD approach followed

### VI. PR Workflow Standards ✅
- **Status**: PASS
- **Check**: Standard PR workflow applies, Conventional Commits format
- **Rationale**: Repository follows standard Git workflow

### VII. Design Documentation ✅
- **Status**: PASS
- **Check**: Design decisions documented in `specs/001-gemara-authoring-assistant/` directory
- **Rationale**: Specification and planning artifacts provide design context

### VIII. Incremental Improvement ✅
- **Status**: PASS
- **Check**: Feature broken into three independent user stories (P1, P2, P3)
- **Rationale**: Incremental delivery approach enables independent testing and value delivery

**Overall Status**: ✅ **PASS** - All constitution gates pass. Proceeding to Phase 0 research.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
cmd/gemara-mcp-server/
├── main.go              # CLI entry point
└── root/
    ├── cmd.go           # Root command
    ├── serve.go         # Serve command (MCP server)
    └── version.go       # Version command

mcp/
├── server.go            # MCP server implementation
└── server_test.go       # Server tests

tools/
├── authoring/           # MCP tool handlers for authoring journeys
│   ├── parse_evidence.go
│   ├── generate_artifact.go
│   ├── generate_policy.go
│   ├── query_info.go
│   ├── query_threats.go
│   ├── validate_artifact.go
│   ├── search_inheritance.go
│   ├── import_controls.go
│   ├── analyze_pivot.go
│   ├── gap_analysis.go
│   ├── prioritize.go
│   ├── recommendations.go
│   ├── confidence.go
│   ├── compare_controls.go
│   ├── error_handler.go
│   ├── logging.go
│   ├── register_tools.go
│   ├── register_prompts.go
│   ├── utils.go
│   └── test_helpers.go
├── info/                # Info tool handlers
│   ├── gemara_info.go
│   ├── resources.go
│   ├── tools.go
│   └── validation.go
├── prompts/             # LLM prompt templates
│   ├── create-layer1.md
│   ├── create-layer2.md
│   ├── create-layer3.md
│   ├── gemara-context.md
│   ├── quick-start.md
│   └── prompts.go
└── tools.go            # Tool registration

internal/
├── consts/
│   └── consts.go        # Centralized constants
├── errors.go            # Error handling infrastructure
├── metrics/             # Observability
│   ├── metrics.go
│   └── prometheus.go
├── parsing/             # Technical evidence parsing
│   ├── interface.go
│   ├── types.go
│   ├── file_based.go
│   ├── dependency_parser.go
│   ├── regulatory_parser.go
│   └── parsers_test.go
├── storage/             # Gemara artifact storage
│   ├── interface.go
│   ├── types.go
│   ├── file_storage.go
│   ├── ranking.go
│   └── storage_test.go
└── validation/          # CUE schema validation
    ├── cue_validator.go
    ├── schemas.go
    ├── version_manager.go
    ├── layer1_validator.go
    ├── layer2_validator.go
    ├── layer3_validator.go
    ├── schemas/
    │   └── README.md
    └── validator_test.go

storage/                # Top-level storage interface
├── interface.go
└── storage.go

tests/
├── contract/            # Contract tests for MCP tools
│   ├── parse_technical_evidence_test.go
│   ├── generate_layer2_artifact_test.go
│   ├── generate_layer3_policy_test.go
│   ├── query_gemara_info_test.go
│   ├── query_threat_library_test.go
│   ├── validate_gemara_artifact_test.go
│   ├── search_inheritance_opportunities_test.go
│   ├── import_inherited_controls_test.go
│   ├── analyze_framework_pivot_test.go
│   └── test_helpers.go
├── integration/         # Integration tests
│   ├── auto_documentation_test.go
│   ├── inheritance_discovery_test.go
│   ├── framework_pivot_test.go
│   ├── layer1_layer3_validation_test.go
│   └── layer3_policy_generation_test.go
└── unit/                # Unit tests
    └── validation_test.go
```

**Structure Decision**: Single project structure with clear separation of concerns:
- `cmd/` - CLI application entry points
- `mcp/` - MCP server protocol implementation
- `tools/` - MCP tool handlers organized by domain (authoring, info)
- `internal/` - Core business logic (parsing, storage, validation)
- `tests/` - Test suites organized by test type (contract, integration, unit)

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No violations - all constitution gates pass.

---

## Phase 0: Outline & Research

**Status**: ✅ **COMPLETE**

All research tasks completed. Key decisions documented in `research.md`:

1. **MCP Server Architecture**: Use `github.com/mark3labs/mcp-go` framework with tool-based architecture for stateless LLM-server interaction
2. **CUE Schema Validation**: Use `cuelang.org/go` with Gemara schema definitions for deterministic validation
3. **Config Parsing Interface**: Extensible parser interface with file-based implementation - parsing ONLY for obscure formats (Dockerfile, Kubernetes), common formats (YAML, JSON, Markdown, text) passed directly to LLM
4. **Gemara Schema Import**: CRITICAL - All Gemara schema types MUST be imported from `github.com/gemaraproj/go-gemara`. Schema types MUST NOT be redefined locally. If import fails, build MUST fail.
5. **Storage Design**: File-based storage for Gemara artifacts (Layer 1, 2, 3) with query interface
6. **Observability**: Use OpenTelemetry metrics API (`go.opentelemetry.io/otel/metric`) as primary instrumentation. Prometheus export (`go.opentelemetry.io/otel/exporters/prometheus`) is optional for compatibility.

**Clarifications Resolved**:
- Parsing strategy: Only parse obscure formats, pass common formats to LLM
- Schema import policy: Strict requirement to import from go-gemara package
- Storage approach: File-based, no database required
- Observability: OTEL metrics primary, Prometheus optional

**Output**: `research.md` complete with all decisions.

---

## Phase 1: Design & Contracts

**Status**: ✅ **COMPLETE**

### Data Model

**Output**: `data-model.md` complete with:
- Dual approach: Gemara-native schema structures + CUE schema field relationships
- Entity definitions: TechnicalEvidence, Capability, Threat, Control, Layer2Artifact, Layer1Guidance, Layer3Policy
- Supporting types: AttackPatternReference, ThreatReference, GuidelineReference, ComplianceTargetReference, ControlReference
- Data flow pipelines for all three journeys
- Validation rules and state transitions

**Key Design Decisions**:
- Use Gemara-native structures from official schemas
- Express relationships via CUE schema fields (not RDF predicates)
- Support both native schema validation and explicit relationship tracking

### API Contracts

**Output**: `contracts/mcp_tools.yaml` complete with:
- MCP tool definitions for all three journeys
- Input/output schemas for each tool
- Error handling specifications
- Tool dependencies and execution order

**Tools Defined**:
- Auto-Documentation: `parse_technical_evidence`, `query_threat_library`, `generate_layer2_artifact`, `validate_gemara_artifact`
- Inheritance Discovery: `search_inheritance_opportunities`, `import_inherited_controls`
- Framework Pivot: `analyze_framework_pivot`, `generate_layer3_policy`
- Info Tools: `query_gemara_info`

### Quick Start Guide

**Output**: `quickstart.md` complete with:
- Prerequisites and setup instructions
- Step-by-step examples for all three journeys
- Tool usage patterns and expected outputs
- Integration scenarios with external MCP servers

### Agent Context Update

**Status**: ⏳ **PENDING**

Agent context update script will be run after Phase 1 completion to add new technology from current plan to agent-specific context files.

**Command**: `.specify/scripts/bash/update-agent-context.sh cursor-agent`

---

## Phase 2: Implementation Planning

**Status**: ⏳ **READY FOR TASKS**

Phase 2 will be handled by `/speckit.tasks` command, which will break down the plan into concrete implementation tasks.

**Prerequisites Met**:
- ✅ Research complete (Phase 0)
- ✅ Design complete (Phase 1)
- ✅ Contracts defined
- ✅ Data model documented
- ✅ Quick start guide available

**Next Steps**:
1. Run `/speckit.tasks` to generate task breakdown
2. Begin implementation following task list
3. Execute tests and validation
4. Complete polish and documentation

---

## Post-Phase 1 Constitution Re-Check

**Status**: ✅ **PASS**

All constitution gates re-evaluated after Phase 1 design:

- **Dependency Management**: ✅ All dependencies use latest stable versions, pinned in go.mod
- **Code Style Standards**: ✅ Go conventions followed, constants centralized
- **Centralized Constants**: ✅ `internal/consts/consts.go` used for all constants
- **Required Questions**: ✅ User personas and problems clearly defined in spec
- **Testing Requirements**: ✅ Test structure defined (unit, integration, contract)
- **PR Workflow**: ✅ Standard workflow applies
- **Design Documentation**: ✅ All design artifacts complete (research.md, data-model.md, contracts/)
- **Incremental Improvement**: ✅ Three independent user stories enable incremental delivery

**Overall Status**: ✅ **PASS** - Ready for Phase 2 task breakdown.

---

## Summary

Implementation plan complete for Gemara Artifact Authoring Assistant. All research clarifications resolved, design artifacts generated. System focuses on LLM-assisted authoring with deterministic validation, supporting three user journeys: Auto-Documentation, Inheritance Discovery, and Framework Pivot. Ready for task breakdown and implementation.

**Branch**: `001-gemara-authoring-assistant`  
**Plan Path**: `/home/jpower/Documents/upstream-repos/gemara-mcp-server/specs/001-gemara-authoring-assistant/plan.md`  
**Generated Artifacts**:
- `research.md` - Phase 0 research decisions
- `data-model.md` - Phase 1 data model design
- `contracts/mcp_tools.yaml` - Phase 1 API contracts
- `quickstart.md` - Phase 1 quick start guide
- `plan.md` - This implementation plan
