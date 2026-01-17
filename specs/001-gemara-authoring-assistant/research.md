# Research: Gemara Artifact Authoring Assistant

**Date**: 2025-01-27  
**Feature**: 001-gemara-authoring-assistant  
**Status**: Complete

## Research Tasks

### 1. MCP Server Architecture for Authoring Workflows

**Task**: Research MCP server patterns for LLM-assisted authoring workflows with deterministic validation

**Decision**: Use existing `github.com/mark3labs/mcp-go` framework with tool-based architecture. Tools provide structured interfaces for LLM interaction. Server acts as stateless processor receiving raw text, providing context tools, and validating outputs.

**Rationale**: 
- MCP protocol designed for LLM-server interaction
- Tool-based pattern enables structured data exchange
- Stateless design aligns with NFR-001 requirement
- Existing codebase already uses mcp-go

**Alternatives Considered**:
- Custom protocol: Rejected - MCP provides standard interface for IDE integration
- Stateful server: Rejected - Violates NFR-001 statelessness requirement

### 2. CUE Schema Validation for Gemara Artifacts

**Task**: Research CUE validation patterns for ensuring 90% deterministic outcomes

**Decision**: Use `cuelang.org/go` CUE engine with Gemara schema definitions. Implement version-aware validation allowing users to specify Gemara spec version. CUE provides structural validation ensuring artifact correctness.

**Rationale**:
- CUE designed for configuration validation and data validation
- Strong typing and constraint system enables deterministic validation
- Gemara framework likely provides CUE schemas
- Version specification addresses FR-005 requirement

**Alternatives Considered**:
- JSON Schema validation: Rejected - Less powerful constraint system, lower determinism
- Custom validation: Rejected - CUE provides proven validation framework

**Implementation Notes**:
- Load Gemara CUE schemas for Layer 1, 2, 3 artifacts
- Support version selection for schema loading
- Provide detailed validation reports with error locations

### 3. Config Parsing Interface Design

**Task**: Research interface patterns for extensible config parsing

**Decision**: Define Go interface for config parsers with file-based default implementation. Interface allows pluggable parsers for different evidence types (config files, code, documentation).

**Rationale**:
- Interface pattern enables extensibility (FR-003 requirement)
- File-based default provides lightweight implementation
- Allows future parsers for SBOMs, CALM artifacts, etc.
- Stateless operation - parsers process input without state

**Alternatives Considered**:
- Monolithic parser: Rejected - Violates extensibility requirement
- External parsing service: Rejected - Adds complexity, violates statelessness

**Interface Design**:
```go
type ConfigParser interface {
    ParseEvidence(source string, content []byte) ([]SecurityFeature, error)
    SupportedFormats() []string
}
```

### 4. Gemara Info Storage Interface

**Task**: Research storage patterns for queryable Gemara information (Layer 1/2/3 artifacts)

**Decision**: Define storage interface with file-based implementation. Storage provides query capabilities for Layer 1 guidance, Layer 2 catalogs, Layer 3 policies. Stateless operation - storage queries are request-scoped.

**Rationale**:
- Interface enables different storage backends (file, in-memory, future: database)
- File-based default aligns with lightweight deployment
- Query interface supports inheritance discovery and framework pivot journeys
- Stateless - no persistence between requests

**Alternatives Considered**:
- Database storage: Deferred - File-based sufficient for MVP, can add later
- In-memory only: Rejected - Need to load existing artifacts for queries

**Interface Design**:
```go
type GemaraStorage interface {
    QueryLayer1(guidanceID string) (*Layer1Guidance, error)
    QueryLayer2(catalogID string) (*Layer2Catalog, error)
    SearchLayer2(query SearchQuery) ([]*Layer2Catalog, error)
    LoadArtifact(path string) (Artifact, error)
}
```

### 5. Auto-Documentation Journey Data Flow

**Task**: Research capability-threat-control mapping patterns for automated artifact generation

**Decision**: Implement workflow: Parse evidence → Extract capabilities → Map to threats → Propose controls → Validate against audit minimums → Generate artifact.

**Rationale**:
- User-provided example flow matches compliance engineering practices
- Threat library provides mapping from capabilities to threats
- NIST 800-53 audit minimums provide regulatory baseline
- CUE validation ensures structural correctness

**Data Flow**:
1. **Capability Definition** (FR-003): Config parser extracts security features/capabilities
2. **Threat Mapping** (FR-006): Query threat library from Gemara info storage
3. **Control Selection**: LLM proposes controls based on threats
4. **Audit Gap Analysis** (FR-004, FR-006): Validate against Layer 1 guidance (NIST 800-53 minimums)
5. **Verification** (FR-004): CUE validation ensures correct structure

**Alternatives Considered**:
- Manual control selection: Rejected - Defeats automation purpose
- Pre-defined control templates: Considered - May add as enhancement

### 6. Inheritance Discovery Pattern

**Task**: Research patterns for identifying inherited controls from dependency information

**Decision**: Use dependency information (SBOMs, architecture diagrams, CALM artifacts) to identify relationships. Search Layer 2 catalogs by technology domain, control type, or dependency relationships. Rank suggestions by relevance.

**Rationale**:
- Dependency information provides relationship context
- Search interface enables catalog discovery
- Relevance ranking improves user experience
- Import mechanism enables control reuse

**Implementation Approach**:
- Parse dependency information to extract relationships
- Build search index from Layer 2 catalogs (technology, domain, controls)
- Match current authoring context against catalog metadata
- Rank by relevance (exact match > partial match > related)

### 7. Framework Pivot Analysis Pattern

**Task**: Research patterns for comparing existing controls against unstructured regulatory requirements

**Decision**: Use LLM to extract structured requirements from unstructured text. Compare against existing Layer 2 controls. Generate gap analysis report with priorities.

**Rationale**:
- LLM excels at extracting structure from unstructured text
- Comparison against existing controls identifies coverage
- Prioritization enables focused compliance efforts
- Report format supports audit readiness

**Analysis Flow**:
1. Parse unstructured regulatory requirements (PDF, text, web)
2. Extract structured requirements using LLM
3. Compare against existing Layer 2 controls
4. Identify gaps and partial coverage
5. Prioritize by criticality
6. Generate actionable report

### 8. Dual Transport Support (Stdio + HTTP)

**Task**: Research MCP transport patterns for local and remote deployment

**Decision**: Support both stdio (local IDE) and HTTP (cloud) transports. Use mcp-go transport abstraction. HTTP transport requires TLS encryption and OAuth2 authentication.

**Rationale**:
- Stdio enables local IDE integration (Cursor/VS Code)
- HTTP enables cloud-native deployment
- TLS encryption required for remote (NFR-004)
- OAuth2 provides secure authentication (NFR-005)

**Implementation Notes**:
- mcp-go supports transport abstraction
- Stdio: Direct stdin/stdout communication
- HTTP: Streamable HTTP with TLS 1.3
- Session isolation via MCP-Session-Id header

### 9. Deterministic Outcome Strategies

**Task**: Research approaches to achieve 90% deterministic artifact generation

**Decision**: Combine structured context tools (threat library, control templates) with CUE validation. Context tools provide consistent reference data. CUE validation ensures structural correctness. LLM uses structured context to generate consistent artifacts.

**Rationale**:
- Structured context reduces LLM variability
- CUE validation catches structural inconsistencies
- Version-aware schemas ensure consistency
- 90% target accounts for LLM non-determinism in content generation

**Strategies**:
- Provide comprehensive threat library as context
- Use control templates where applicable
- Validate structure strictly with CUE
- Accept semantic variation (content) while ensuring structural consistency

### 10. Observability and Metrics

**Task**: Research metrics patterns for MCP server observability

**Decision**: Export Prometheus-compatible metrics. Track domain-specific metrics: `gemara_mapping_success_rate`, `gemara_schema_validation_failures_total`. Integrate with central collector.

**Rationale**:
- Prometheus standard for metrics collection
- Domain metrics track feature-specific outcomes
- Central collector enables monitoring and alerting
- Aligns with NFR-007 and NFR-008 requirements

**Metrics Design**:
- Counter: `gemara_schema_validation_failures_total` (by version, layer)
- Gauge: `gemara_mapping_success_rate` (percentage)
- Histogram: Request processing time
- Counter: Tool invocation counts by type

### 11. Control ID Format and Immutability

**Task**: Research control ID format requirements and immutability constraints

**Decision**: Control IDs MUST use format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042"). Control IDs MUST NOT include family in the identifier. Control IDs are immutable once defined until withdrawn, even if controls are reclassified into different families.

**Rationale**:
- Control IDs serve as stable references across catalog versions
- Family reclassification is common but should not break existing references
- Format `<identifier>-<numbering>` provides clear structure without coupling to family
- Immutability ensures referential integrity and prevents breaking changes

**Alternatives Considered**:
- Family-prefixed IDs (e.g., "AC-AC-001"): Rejected - Breaks when controls move families
- UUID-based IDs: Rejected - Not human-readable, harder to reference
- Sequential numbering only: Rejected - Lacks identifier prefix for namespacing

---

## Summary

All research tasks completed. Key decisions:
1. MCP tool-based architecture with stateless operation
2. CUE validation for deterministic structural correctness
3. Interface-based design for parsing and storage (extensibility)
4. Structured data flows for each journey (5-phase pipeline for Auto-Documentation)
5. Dual transport support (stdio + HTTP)
6. Prometheus metrics for observability
7. Control ID format: `<identifier>-<numbering>` (immutable, family-independent)

No blocking clarifications remain. Ready for Phase 1 design.
