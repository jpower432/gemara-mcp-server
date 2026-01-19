# Tasks: Gemara Artifact Authoring Assistant

**Input**: Design documents from `/specs/001-gemara-authoring-assistant/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are included as contract tests and integration tests per the specification requirements.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Paths follow Go project structure at repository root:
- `cmd/gemara-mcp-server/` - Server entry point
- `mcp/` - MCP server implementation
- `tools/authoring/` - Authoring tools
- `tools/info/` - Info tools
- `internal/` - Internal packages
- `tests/` - Test files

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create internal/parsing/ directory structure
- [X] T002 Create internal/validation/ directory structure
- [X] T003 Create internal/storage/ directory structure
- [X] T004 Create internal/metrics/ directory structure
- [X] T005 Create tests/contract/ directory structure
- [X] T006 Create tests/integration/ directory structure
- [X] T007 [P] Add CUE schema files for Layer 1, Layer 2, Layer 3 validation in internal/validation/schemas/ - schemas loaded from github.com/gemaraproj/go-gemara package. CRITICAL: Schema types MUST NOT be redefined locally. All Gemara schema types MUST be imported from github.com/gemaraproj/go-gemara. If import fails, build MUST fail with clear error message.
- [X] T008 [P] Update go.mod with required dependencies (cuelang.org/go, github.com/goccy/go-yaml, github.com/mark3labs/mcp-go, github.com/gemaraproj/go-gemara, github.com/spf13/cobra, github.com/stretchr/testify, go.opentelemetry.io/otel, go.opentelemetry.io/otel/exporters/prometheus). CRITICAL: Verify github.com/gemaraproj/go-gemara is available and importable. If import fails, stop build with error. Do NOT use github.com/ossf/gemara or any other package. All Gemara schema types MUST be imported from github.com/gemaraproj/go-gemara.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T009 Add constants for all magic strings/numbers to internal/consts/consts.go per constitution (parser names, format types, error codes, metric names, control ID format pattern regex: `^[A-Z0-9]+-[0-9]{3,}$`)
- [X] T010 [P] Create internal/errors.go with error handling infrastructure and error types
- [X] T011 [P] Create internal/parsing/interface.go with ConfigParser interface definition
- [X] T012 [P] Create internal/parsing/types.go with parsing-related types (TechnicalEvidence, Capability)
- [X] T013 [P] Create internal/storage/interface.go with GemaraStorage interface for querying Layer 1, Layer 2 (including Threat Catalog), Layer 3 (including Risk Catalog) artifacts
- [X] T014 [P] Create internal/storage/types.go with storage-related types. CRITICAL: Import all Gemara schema types (Layer1Guidance, Layer2Catalog, ThreatCatalog, Layer3Policy, RiskCatalog, etc.) from github.com/gemaraproj/go-gemara. DO NOT redefine these types locally. Only add CUE schema field relationship types (AttackPatternReference, ThreatReference, GuidelineReference, ComplianceTargetReference, ControlReference) as wrappers/extensions if needed. If github.com/gemaraproj/go-gemara import fails, build MUST fail.
- [X] T015 [P] Create internal/validation/cue_validator.go with base CUE validation infrastructure. CRITICAL: Import all CUE schema definitions from github.com/gemaraproj/go-gemara. Do NOT define schemas locally. If import fails, build MUST fail.
- [X] T016 [P] Create internal/validation/version_manager.go for Gemara spec version management. CRITICAL: Use version information from github.com/gemaraproj/go-gemara package. Do NOT hardcode versions locally.
- [X] T017 [P] Create internal/validation/schemas.go for loading CUE schemas for all three Definition layers (Layer 1, Layer 2 including Threat Catalog, Layer 3 including Risk Catalog). CRITICAL: Import all schema definitions from github.com/gemaraproj/go-gemara. Do NOT define schemas locally. If import fails, build MUST fail.
- [X] T018 [P] Create internal/validation/layer1_validator.go for Layer 1 GuidanceDocument validation. CRITICAL: Use GuidanceDocument type from github.com/gemaraproj/go-gemara. Do NOT redefine locally.
- [X] T019 [P] Create internal/validation/layer2_validator.go for Layer 2 Catalog validation (including Threat Catalog). CRITICAL: Use Catalog and ThreatCatalog types from github.com/gemaraproj/go-gemara. Do NOT redefine locally.
- [X] T020 [P] Create internal/validation/layer3_validator.go for Layer 3 Policy validation (including Risk Catalog). CRITICAL: Use Policy and RiskCatalog types from github.com/gemaraproj/go-gemara. Do NOT redefine locally.
- [X] T021 Create internal/storage/file_storage.go with file-based storage implementation (implements GemaraStorage interface, supports querying Threat Catalog and Risk Catalog)
- [X] T022 Create internal/metrics/metrics.go with OpenTelemetry metrics instrumentation using go.opentelemetry.io/otel/metric API. Define domain metrics: gemara_mapping_success_rate (gauge), gemara_schema_validation_failures_total (counter). Use OTEL metrics API as primary instrumentation mechanism.
- [X] T023 Create tools/authoring/types.go with shared types for authoring tools (Gap, Recommendation, ThreatMapping, MultiMapping)
- [X] T024 Create tools/authoring/error_handler.go with error handling utilities for MCP tools
- [X] T025 Create tools/authoring/logging.go with logging utilities for MCP tools
- [X] T026 Create tools/authoring/utils.go with helper functions (extractStringArray, etc.)
- [X] T026a [P] Implement TLS 1.3 encryption for Streamable HTTP transport in mcp/server.go (NFR-004) - ensure all remote communications are encrypted
- [X] T026b [P] Implement OAuth2/OIDC authentication with PKCE for remote mode in mcp/server.go (NFR-005) - secure user authentication before processing requests
- [X] T026c [P] Implement stdio transport for local IDE integration in mcp/server.go (NFR-002) - supports local development mode via stdin/stdout
- [X] T026d [P] Implement HTTP transport abstraction in mcp/server.go (NFR-002) - supports cloud-native deployment with Streamable HTTP
- [X] T026e [P] Implement session isolation via MCP-Session-Id header handling in mcp/server.go (NFR-006) - ensure logical isolation of request-scoped data between user sessions
- [X] T026f [P] Add build-time validation in internal/validation/schemas.go to verify github.com/gemaraproj/go-gemara imports successfully and provides required schema types (Layer1Guidance, Layer2Catalog, ThreatCatalog, Layer3Policy, RiskCatalog). If import fails or types are missing, build MUST fail with clear error message. Do NOT define fallback local types.

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Auto-Documentation (Priority: P1) 🎯 MVP

**Goal**: Generate validated Layer 2 Gemara artifacts from raw technical evidence via LLM-driven 5-phase pipeline: (1) Capability Definition, (2) Threat Mapping, (3) Control Selection, (4) Audit Gap Analysis, (5) Verification.

**Independent Test**: Can be fully tested by providing sample technical evidence (e.g., security configuration files, policy documents) and verifying that a valid Layer 2 Gemara artifact is generated with correct structure, required fields populated, and validation passing. The generated artifact can be independently reviewed and used for audit purposes.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [X] T027 [P] [US1] Contract test for parse_technical_evidence tool in tests/contract/parse_technical_evidence_test.go
- [X] T028 [P] [US1] Contract test for query_gemara_info tool (Layer 1, Layer 2, Layer 3, threat_catalog, risk_catalog) in tests/contract/query_gemara_info_test.go
- [X] T029 [P] [US1] Contract test for query_threat_library tool in tests/contract/query_threat_library_test.go
- [X] T030 [P] [US1] Contract test for validate_gemara_artifact tool (Layer 2) in tests/contract/validate_gemara_artifact_test.go
- [X] T031 [P] [US1] Contract test for generate_layer2_artifact tool in tests/contract/generate_layer2_artifact_test.go
- [X] T032 [P] [US1] Integration test for Auto-Documentation journey in tests/integration/auto_documentation_test.go

### Implementation for User Story 1

- [X] T033 [P] [US1] Create internal/parsing/file_based.go with YAML parser implementation (implements ConfigParser interface, only for obscure YAML formats)
- [X] T034 [P] [US1] Create internal/parsing/file_based.go with JSON parser implementation (extends ConfigParser, only for obscure JSON formats)
- [X] T035 [P] [US1] Create internal/parsing/file_based.go with text parser implementation (extends ConfigParser, only for obscure text formats)
- [X] T036 [US1] Implement parse_technical_evidence MCP tool handler in tools/authoring/parse_evidence.go (depends on T011, T012, T033-T035) - passes common formats directly to LLM, parses obscure formats
- [X] T037 [US1] Implement query_gemara_info MCP tool handler in tools/authoring/query_info.go (depends on T013, T014, T021) - supports querying Layer 1, Layer 2, Layer 3 artifacts, threat_catalog, risk_catalog query types
- [X] T038 [US1] Implement query_threat_library MCP tool handler in tools/authoring/query_threats.go (depends on T013, T014, T037) - queries Threat Catalog via query_gemara_info with query_type="threat_catalog"
- [X] T039 [US1] Implement validate_gemara_artifact MCP tool handler for Layer 2 in tools/authoring/validate_artifact.go (depends on T015-T020)
- [X] T040 [US1] Implement generate_layer2_artifact MCP tool handler in tools/authoring/generate_artifact.go (depends on T036, T037, T038, T039) - orchestrates 5-phase pipeline: (1) Capability Definition (common formats → LLM, obscure formats → parse), (2) Threat Mapping (via query_threat_library using Gemara-native structure AND "identifiedBy" Attack Pattern CUE schema field), (3) Control Selection (LLM creates Families using Gemara-native structure AND "mitigates" Threat CUE schema field), (4) Audit Gap Analysis (using Control "satisfies" Guideline CUE schema field relationships), (5) Verification
- [X] T041 [US1] Add control ID format validation in tools/authoring/generate_artifact.go - validate format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042"), ensure no family in ID, enforce immutability
- [X] T042 [US1] Implement audit gap analysis logic in tools/authoring/generate_artifact.go - check proposed controls against Layer 1 audit minimums using Control "satisfies" Guideline CUE schema field relationships, flag gaps as recommendations (depends on T037, T038)
- [X] T043 [US1] Register parse_technical_evidence, query_gemara_info, query_threat_library, validate_gemara_artifact, generate_layer2_artifact tools in tools/authoring/register_tools.go
- [X] T044 [US1] Add error handling and logging to all US1 MCP tool handlers in tools/authoring/
- [X] T045 [US1] Add OpenTelemetry metrics tracking for artifact generation success rate in tools/authoring/generate_artifact.go using go.opentelemetry.io/otel/metric API

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Inheritance Discovery (Priority: P2)

**Goal**: Enable users to discover and import existing Layer 2 catalogs containing relevant controls, reducing redundant documentation work through inheritance suggestions and import capabilities.

**Independent Test**: Can be fully tested by providing a partial Layer 2 catalog or control context and verifying that the system suggests relevant existing Layer 2 catalogs with matching controls. The suggestions can be independently evaluated for relevance and accuracy without implementing other journeys.

### Tests for User Story 2

- [X] T046 [P] [US2] Contract test for search_inheritance_opportunities tool in tests/contract/search_inheritance_opportunities_test.go
- [X] T047 [P] [US2] Contract test for import_inherited_controls tool in tests/contract/import_inherited_controls_test.go
- [X] T048 [P] [US2] Integration test for Inheritance Discovery journey in tests/integration/inheritance_discovery_test.go

### Implementation for User Story 2

- [X] T049 [P] [US2] Create internal/storage/ranking.go with relevance ranking algorithm (exact matches > partial matches > related matches)
- [X] T050 [P] [US2] Create internal/parsing/dependency_parser.go for parsing SBOMs, architecture diagrams, CALM artifacts
- [X] T051 [US2] Implement search_inheritance_opportunities MCP tool handler in tools/authoring/search_inheritance.go (depends on T013, T014, T037, T049)
- [X] T052 [US2] Create tools/authoring/compare_controls.go with control comparison logic for matching controls across catalogs
- [X] T053 [US2] Implement import_inherited_controls MCP tool handler in tools/authoring/import_controls.go (depends on T051, T052) - uses `imported-controls` field with MultiMapping structure (`[...#MultiMapping]` with `@go(ImportedControls)` tag)
- [X] T054 [US2] Register search_inheritance_opportunities and import_inherited_controls tools in tools/authoring/register_tools.go
- [X] T055 [US2] Add support for imported-threats and imported-capabilities fields in tools/authoring/import_controls.go
- [X] T056 [US2] Add error handling and logging to all US2 MCP tool handlers
- [X] T057 [US2] Add OpenTelemetry metrics tracking for inheritance discovery operations using go.opentelemetry.io/otel/metric API

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Framework Pivot (Priority: P3)

**Goal**: Analyze existing Layer 2 controls against unstructured regulatory requirements to produce prioritized gap analysis reports identifying covered requirements, gaps, and uncovered minimum requirements.

**Independent Test**: Can be fully tested by providing existing Layer 2 controls and unstructured regulatory requirements, then verifying that the system produces a prioritized report identifying gaps, covered requirements, and minimum requirements that remain uncovered. The report can be independently reviewed for accuracy and completeness.

### Tests for User Story 3

- [X] T058 [P] [US3] Contract test for analyze_framework_pivot tool in tests/contract/analyze_framework_pivot_test.go
- [X] T059 [P] [US3] Contract test for generate_layer3_policy tool in tests/contract/generate_layer3_policy_test.go
- [X] T060 [P] [US3] Contract test for validate_gemara_artifact tool (Layer 3) in tests/contract/validate_layer3_artifact_test.go
- [X] T061 [P] [US3] Integration test for Framework Pivot journey in tests/integration/framework_pivot_test.go
- [X] T062 [P] [US3] Integration test for Layer 3 Policy generation in tests/integration/layer3_policy_generation_test.go

### Implementation for User Story 3

- [X] T063 [P] [US3] Create internal/parsing/regulatory_parser.go for parsing unstructured regulatory requirements (text, PDF, web content)
- [X] T064 [US3] Create tools/authoring/gap_analysis.go with gap analysis engine for comparing controls against requirements
- [X] T065 [US3] Create tools/authoring/prioritize.go with prioritization logic for uncovered requirements
- [X] T066 [US3] Create tools/authoring/recommendations.go with recommendation generator for addressing gaps
- [X] T067 [US3] Create tools/authoring/confidence.go with confidence indicator calculation (float64 0.0-1.0)
- [X] T067a [US3] Document confidence calculation factors in tools/authoring/confidence.go per FR-016 - factors: input quality, threat match confidence, control coverage percentage, regulatory requirement clarity, schema validation status
- [X] T068 [US3] Implement analyze_framework_pivot MCP tool handler in tools/authoring/analyze_pivot.go (depends on T037, T063, T064, T065, T066, T067) - note: Risk Catalogs not directly used
- [X] T069 [US3] Implement generate_layer3_policy MCP tool handler in tools/authoring/generate_policy.go (depends on T013, T014, T037) - generates Policy via scope definition with Layer 1 and Layer 2 applicability queried for context, uses Guideline "establishes" Compliance Target CUE schema field relationships
- [X] T070 [US3] Extend validate_gemara_artifact MCP tool handler for Layer 3 in tools/authoring/validate_artifact.go (extends T039)
- [X] T071 [US3] Register analyze_framework_pivot and generate_layer3_policy tools in tools/authoring/register_tools.go
- [X] T072 [US3] Add error handling and logging to all US3 MCP tool handlers
- [X] T073 [US3] Add OpenTelemetry metrics tracking for framework pivot operations using go.opentelemetry.io/otel/metric API

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T074 [P] Add unit tests for parsing components in internal/parsing/parsers_test.go
- [X] T075 [P] Add unit tests for validation components in internal/validation/validator_test.go
- [X] T076 [P] Add unit tests for storage components in internal/storage/storage_test.go
- [X] T077 [P] Add unit tests for gap analysis in tools/authoring/gap_analysis_test.go
- [X] T078 [P] Add OpenTelemetry metrics export for gemara_mapping_success_rate and gemara_schema_validation_failures_total in internal/metrics/metrics.go using go.opentelemetry.io/otel/metric API. Optional: Add Prometheus exporter (go.opentelemetry.io/otel/exporters/prometheus) if Prometheus compatibility is required.
- [X] T079 [P] Update documentation in README.md with usage examples for all three journeys
- [X] T080 [P] Validate quickstart.md examples work with implemented tools (examples match tool signatures and expected outputs)
- [X] T081 Add comprehensive error messages following FR-004 validation report structure (path, message, severity) across all tools
- [X] T082 Ensure all tools support Gemara spec version specification per FR-005
- [X] T083 Verify all constants are centralized in internal/consts/consts.go per constitution (including control ID format pattern)
- [X] T084 [P] Add control ID format validation tests in tests/unit/validation_test.go - test format `<identifier>-<numbering>` matching regex `^[A-Z0-9]+-[0-9]{3,}$` (e.g., "AC-001", "SEC-042"), test immutability, test family-independence
- [X] T085 [P] Add integration tests for Layer 1 and Layer 3 validation in tests/integration/
- [X] T086 Ensure 90% deterministic outcomes for artifact generation (NFR-003) via CUE validation and context tools (documented in docs/DETERMINISM.md)
- [X] T087 Performance optimization: ensure Layer 2 artifact generation completes in under 10 minutes (SC-001) (implementation designed to meet requirement, documented in docs/DETERMINISM.md)
- [X] T088 Performance optimization: ensure framework pivot gap analysis reports complete in under 15 minutes (SC-005)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Final Phase)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - May integrate with US1 but should be independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - May integrate with US1/US2 but should be independently testable

### Within Each User Story

- Tests (if included) MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Contract test for parse_technical_evidence tool in tests/contract/parse_technical_evidence_test.go"
Task: "Contract test for query_gemara_info tool in tests/contract/query_gemara_info_test.go"
Task: "Contract test for query_threat_library tool in tests/contract/query_threat_library_test.go"
Task: "Contract test for validate_gemara_artifact tool in tests/contract/validate_gemara_artifact_test.go"
Task: "Contract test for generate_layer2_artifact tool in tests/contract/generate_layer2_artifact_test.go"
Task: "Integration test for Auto-Documentation journey in tests/integration/auto_documentation_test.go"

# Launch all parsers for User Story 1 together:
Task: "Create internal/parsing/file_based.go with YAML parser implementation"
Task: "Create internal/parsing/file_based.go with JSON parser implementation"
Task: "Create internal/parsing/file_based.go with text parser implementation"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1
   - Developer B: User Story 2
   - Developer C: User Story 3
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- CRITICAL: All Gemara schema types MUST be imported from github.com/gemaraproj/go-gemara. Do NOT redefine locally. If import fails, build MUST fail.
- Control IDs MUST follow format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042") matching regex `^[A-Z0-9]+-[0-9]{3,}$`
- Use OpenTelemetry metrics API (go.opentelemetry.io/otel/metric) as primary instrumentation mechanism
- Parsing is ONLY for obscure formats (Dockerfile, Kubernetes) - common formats (YAML, JSON, Markdown, text) passed directly to LLM

---

## Task Summary

- **Total Tasks**: 88
- **Setup Tasks**: 8 (Phase 1)
- **Foundational Tasks**: 22 (Phase 2) - includes transport abstraction, session isolation, and schema import validation
- **User Story 1 Tasks**: 19 (Phase 3) - MVP (includes 5-phase pipeline orchestration, control ID validation, audit gap analysis)
- **User Story 2 Tasks**: 12 (Phase 4)
- **User Story 3 Tasks**: 16 (Phase 5)
- **Polish Tasks**: 15 (Phase 6) (includes control ID format validation tests and performance optimization)

### Parallel Opportunities Identified

- **Phase 1**: 2 parallel tasks
- **Phase 2**: 18 parallel tasks (within foundational phase)
- **Phase 3**: 6 parallel test tasks, 3 parallel parser tasks
- **Phase 4**: 3 parallel test tasks, 2 parallel implementation tasks
- **Phase 5**: 5 parallel test tasks, 1 parallel implementation task
- **Phase 6**: 8 parallel tasks

### Independent Test Criteria

- **User Story 1**: Provide sample technical evidence → verify valid Layer 2 Gemara artifact generated with correct structure and validation passing
- **User Story 2**: Provide partial Layer 2 catalog → verify system suggests relevant existing Layer 2 catalogs with matching controls
- **User Story 3**: Provide existing Layer 2 controls and unstructured regulatory requirements → verify prioritized report identifying gaps, covered requirements, and uncovered minimums

### Suggested MVP Scope

**MVP = Phase 1 + Phase 2 + Phase 3 (User Story 1)**

This delivers Auto-Documentation capability, enabling users to generate validated Layer 2 Gemara artifacts from technical evidence. This addresses the most time-consuming aspect of compliance work and provides immediate value.
