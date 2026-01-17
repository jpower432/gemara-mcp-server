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
- `internal/` - Internal packages
- `tests/` - Test files

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [ ] T001 Create internal/parsing/ directory structure
- [ ] T002 Create internal/validation/ directory structure
- [ ] T003 Create internal/storage/ directory structure
- [ ] T004 Create internal/metrics/ directory structure
- [ ] T005 Create tests/contract/ directory structure
- [ ] T006 Create tests/integration/ directory structure
- [ ] T007 [P] Add CUE schema files for Layer 1, Layer 2, Layer 3 validation in internal/validation/schemas/
- [ ] T008 [P] Update go.mod with required dependencies (cuelang.org/go, github.com/goccy/go-yaml, github.com/mark3labs/mcp-go, github.com/gemaraproj/go-gemara, github.com/spf13/cobra, github.com/stretchr/testify)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T009 Add constants for all magic strings/numbers to internal/consts/consts.go per constitution (parser names, format types, error codes, metric names, control ID format pattern regex: `<identifier>-<numbering>`)
- [ ] T010 [P] Create internal/errors.go with error handling infrastructure and error types
- [ ] T011 [P] Create internal/parsing/interface.go with ConfigParser interface definition
- [ ] T012 [P] Create internal/parsing/types.go with parsing-related types (TechnicalEvidence, SecurityFeature)
- [ ] T013 [P] Create internal/storage/interface.go with GemaraStorage interface for querying Layer 1, Layer 2, Layer 3 artifacts
- [ ] T014 [P] Create internal/storage/types.go with storage-related types
- [ ] T015 [P] Create internal/validation/cue_validator.go with base CUE validation infrastructure
- [ ] T016 [P] Create internal/validation/version_manager.go for Gemara spec version management
- [ ] T017 [P] Create internal/validation/schemas.go for loading CUE schemas for all three Definition layers
- [ ] T018 [P] Create internal/validation/layer1_validator.go for Layer 1 GuidanceDocument validation
- [ ] T019 [P] Create internal/validation/layer2_validator.go for Layer 2 Catalog validation
- [ ] T020 [P] Create internal/validation/layer3_validator.go for Layer 3 Policy validation
- [ ] T021 Create internal/storage/file_storage.go with file-based storage implementation (implements GemaraStorage interface)
- [ ] T022 Create internal/metrics/metrics.go with domain metrics definitions (gemara_mapping_success_rate, gemara_schema_validation_failures_total)
- [ ] T023 Create tools/authoring/types.go with shared types for authoring tools
- [ ] T024 Create tools/authoring/error_handler.go with error handling utilities for MCP tools
- [ ] T025 Create tools/authoring/logging.go with logging utilities for MCP tools

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Auto-Documentation (Priority: P1) 🎯 MVP

**Goal**: Generate validated Layer 2 Catalog artifacts from raw technical evidence with proper structure (title, Controls, Families), relationships, and metadata. System follows 5-phase sequential pipeline: (1) LLM defines capabilities from technical evidence, (2) Server maps capabilities to threats via threat library, (3) LLM proposes controls to mitigate threats (creates Families), (4) Server performs audit gap analysis against Layer 1 requirements, (5) Server validates artifact structure. Control IDs use format `<identifier>-<numbering>` and are immutable.

**Independent Test**: Can be fully tested by providing sample technical evidence (e.g., security configuration files, policy documents) and verifying that a valid Layer 2 Gemara artifact is generated with correct structure, required fields populated, and validation passing. The generated artifact can be independently reviewed and used for audit purposes.

### Tests for User Story 1

> **NOTE: Write these tests FIRST, ensure they FAIL before implementation**

- [ ] T026 [P] [US1] Contract test for parse_technical_evidence tool in tests/contract/test_parse_technical_evidence.go
- [ ] T027 [P] [US1] Contract test for query_gemara_info tool in tests/contract/test_query_gemara_info.go
- [ ] T028 [P] [US1] Contract test for query_threat_library tool in tests/contract/test_query_threat_library.go
- [ ] T029 [P] [US1] Contract test for validate_gemara_artifact tool (Layer 2) in tests/contract/test_validate_gemara_artifact.go
- [ ] T030 [P] [US1] Contract test for generate_layer2_artifact tool in tests/contract/test_generate_layer2_artifact.go
- [ ] T031 [P] [US1] Integration test for Auto-Documentation journey in tests/integration/test_auto_documentation.go

### Implementation for User Story 1

- [ ] T032 [P] [US1] Create internal/parsing/file_based.go with YAML parser implementation (implements ConfigParser interface)
- [ ] T033 [P] [US1] Create internal/parsing/file_based.go with JSON parser implementation (extends ConfigParser)
- [ ] T034 [P] [US1] Create internal/parsing/file_based.go with text parser implementation (extends ConfigParser)
- [ ] T035 [US1] Implement parse_technical_evidence MCP tool handler in tools/authoring/parse_evidence.go (depends on T011, T012, T032-T034)
- [ ] T036 [US1] Implement query_gemara_info MCP tool handler in tools/info/gemara_info.go (depends on T013, T014, T021) - supports querying Layer 1, Layer 2, Layer 3 artifacts
- [ ] T037 [US1] Implement query_threat_library MCP tool handler in tools/authoring/query_threats.go (depends on T013, T014, T036)
- [ ] T038 [US1] Implement validate_gemara_artifact MCP tool handler for Layer 2 in tools/authoring/validate_artifact.go (depends on T015-T020)
- [ ] T039 [US1] Implement generate_layer2_artifact MCP tool handler in tools/authoring/generate_artifact.go (depends on T035, T036, T037, T038) - orchestrates 5-phase pipeline: (1) Capability Definition, (2) Threat Mapping, (3) Control Selection, (4) Audit Gap Analysis, (5) Verification
- [ ] T040 [US1] Add control ID format validation in tools/authoring/generate_artifact.go - validate format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042"), ensure no family in ID, enforce immutability
- [ ] T041 [US1] Implement audit gap analysis logic in tools/authoring/generate_artifact.go - check proposed controls against Layer 1 audit minimums, flag gaps as recommendations (depends on T036, T037)
- [ ] T042 [US1] Register parse_technical_evidence, query_gemara_info, query_threat_library, validate_gemara_artifact, generate_layer2_artifact tools in tools/authoring/register_tools.go
- [ ] T043 [US1] Add error handling and logging to all US1 MCP tool handlers in tools/authoring/
- [ ] T044 [US1] Add metrics tracking for artifact generation success rate in tools/authoring/generate_artifact.go

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently. Users can provide technical evidence and receive validated Layer 2 Catalog artifacts.

---

## Phase 4: User Story 2 - Inheritance Discovery (Priority: P2)

**Goal**: Identify inherited security controls to reduce redundant documentation. System analyzes current documentation context and suggests existing Layer 2 catalogs containing relevant controls, enabling import via Gemara `imported-controls` field with proper attribution.

**Independent Test**: Can be fully tested by providing a partial Layer 2 catalog or control context and verifying that the system suggests relevant existing Layer 2 catalogs with matching controls. The suggestions can be independently evaluated for relevance and accuracy without implementing other journeys.

### Tests for User Story 2

- [ ] T045 [P] [US2] Contract test for search_inheritance_opportunities tool in tests/contract/test_search_inheritance.go
- [ ] T046 [P] [US2] Contract test for import_inherited_controls tool in tests/contract/test_import_controls.go
- [ ] T047 [P] [US2] Integration test for Inheritance Discovery journey in tests/integration/test_inheritance_discovery.go

### Implementation for User Story 2

- [ ] T048 [P] [US2] Create internal/storage/ranking.go with relevance ranking algorithm (exact matches > partial matches > related matches)
- [ ] T049 [P] [US2] Create internal/parsing/dependency_parser.go for parsing SBOMs, architecture diagrams, CALM artifacts
- [ ] T050 [US2] Implement search_inheritance_opportunities MCP tool handler in tools/authoring/search_inheritance.go (depends on T013, T014, T036, T048)
- [ ] T051 [US2] Create tools/authoring/compare_controls.go with control comparison logic for matching controls across catalogs
- [ ] T052 [US2] Implement import_inherited_controls MCP tool handler in tools/authoring/import_controls.go (depends on T050, T051) - uses `imported-controls` field with MultiMapping structure
- [ ] T053 [US2] Register search_inheritance_opportunities and import_inherited_controls tools in tools/authoring/register_tools.go
- [ ] T054 [US2] Add support for imported-threats and imported-capabilities fields in tools/authoring/import_controls.go
- [ ] T055 [US2] Add error handling and logging to all US2 MCP tool handlers
- [ ] T056 [US2] Add metrics tracking for inheritance discovery operations

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently. Users can search for inheritance opportunities and import controls with proper attribution.

---

## Phase 5: User Story 3 - Framework Pivot (Priority: P3)

**Goal**: Assess current technology against a new, unstructured regulatory standard. System analyzes existing Layer 2 controls and compares them against unstructured regulatory requirements, producing a prioritized gap analysis report with confidence indicators.

**Independent Test**: Can be fully tested by providing existing Layer 2 controls and unstructured regulatory requirements, then verifying that the system produces a prioritized report identifying gaps, covered requirements, and minimum requirements that remain uncovered. The report can be independently reviewed for accuracy and completeness.

### Tests for User Story 3

- [ ] T057 [P] [US3] Contract test for analyze_framework_pivot tool in tests/contract/test_framework_pivot.go
- [ ] T058 [P] [US3] Contract test for generate_layer3_policy tool in tests/contract/test_generate_layer3_policy.go
- [ ] T059 [P] [US3] Contract test for validate_gemara_artifact tool (Layer 3) in tests/contract/test_validate_gemara_artifact.go
- [ ] T060 [P] [US3] Integration test for Framework Pivot journey in tests/integration/test_framework_pivot.go
- [ ] T061 [P] [US3] Integration test for Layer 3 Policy generation in tests/integration/test_layer3_policy_generation.go

### Implementation for User Story 3

- [ ] T062 [P] [US3] Create internal/parsing/regulatory_parser.go for parsing unstructured regulatory requirements (text, PDF, web content)
- [ ] T063 [US3] Create tools/authoring/gap_analysis.go with gap analysis engine for comparing controls against requirements
- [ ] T064 [US3] Create tools/authoring/prioritize.go with prioritization logic for uncovered requirements
- [ ] T065 [US3] Create tools/authoring/recommendations.go with recommendation generator for addressing gaps
- [ ] T066 [US3] Create tools/authoring/confidence.go with confidence indicator calculation (float64 0.0-1.0)
- [ ] T067 [US3] Implement analyze_framework_pivot MCP tool handler in tools/authoring/analyze_pivot.go (depends on T036, T062, T063, T064, T065, T066)
- [ ] T068 [US3] Implement generate_layer3_policy MCP tool handler in tools/authoring/generate_policy.go (depends on T013, T014, T036, FR-006a) - generates Policy via scope definition with Layer 1 and Layer 2 context
- [ ] T069 [US3] Extend validate_gemara_artifact MCP tool handler for Layer 3 in tools/authoring/validate_artifact.go (extends T038)
- [ ] T070 [US3] Register analyze_framework_pivot and generate_layer3_policy tools in tools/authoring/register_tools.go
- [ ] T071 [US3] Add error handling and logging to all US3 MCP tool handlers
- [ ] T072 [US3] Add metrics tracking for framework pivot operations

**Checkpoint**: All user stories should now be independently functional. Users can perform framework pivot analysis and generate Layer 3 Policy artifacts.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T073 [P] Add unit tests for parsing components in internal/parsing/parsers_test.go
- [ ] T074 [P] Add unit tests for validation components in internal/validation/validator_test.go
- [ ] T075 [P] Add unit tests for storage components in internal/storage/storage_test.go
- [ ] T076 [P] Add unit tests for gap analysis in tools/authoring/gap_analysis_test.go
- [ ] T077 [P] Add observability metrics export for gemara_mapping_success_rate and gemara_schema_validation_failures_total in internal/metrics/prometheus.go
- [ ] T078 [P] Update documentation in README.md with usage examples for all three journeys
- [ ] T079 [P] Validate quickstart.md examples work with implemented tools
- [ ] T080 Add comprehensive error messages following FR-004 validation report structure (path, message, severity) across all tools
- [ ] T081 Ensure all tools support Gemara spec version specification per FR-005
- [ ] T082 Verify all constants are centralized in internal/consts/consts.go per constitution (including control ID format pattern)
- [ ] T083 Add control ID format validation tests in tests/unit/validation_test.go - test format `<identifier>-<numbering>`, test immutability, test family-independence
- [ ] T084 Add integration tests for Layer 1 and Layer 3 validation in tests/integration/
- [ ] T085 Ensure 90% deterministic outcomes for artifact generation (NFR-003) via CUE validation and context tools
- [ ] T086 Performance optimization: ensure Layer 2 artifact generation completes in under 10 minutes (SC-001)
- [ ] T087 Performance optimization: ensure framework pivot analysis completes in under 15 minutes (SC-005)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3+)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Polish (Phase 6)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories. MVP scope.
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Uses storage and validation from US1 but independently testable
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Uses storage, validation, and parsing from US1/US2 but independently testable

### Within Each User Story

- Tests (contract + integration) MUST be written and FAIL before implementation
- Parsing/Storage/Validation infrastructure before MCP tool handlers
- MCP tool handlers before tool registration
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Parsing implementations within US1 marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all contract tests for User Story 1 together:
Task: "Contract test for parse_technical_evidence tool in tests/contract/test_parse_technical_evidence.go"
Task: "Contract test for query_gemara_info tool in tests/contract/test_query_gemara_info.go"
Task: "Contract test for query_threat_library tool in tests/contract/test_query_threat_library.go"
Task: "Contract test for validate_gemara_artifact tool (Layer 2) in tests/contract/test_validate_gemara_artifact.go"
Task: "Contract test for generate_layer2_artifact tool in tests/contract/test_generate_layer2_artifact.go"

# Launch parser implementations together:
Task: "Create internal/parsing/file_based.go with YAML parser implementation"
Task: "Create internal/parsing/file_based.go with JSON parser implementation"
Task: "Create internal/parsing/file_based.go with text parser implementation"

# Launch all parser implementations together:
Task: "Create internal/parsing/file_based.go with YAML parser implementation"
Task: "Create internal/parsing/file_based.go with JSON parser implementation"
Task: "Create internal/parsing/file_based.go with text parser implementation"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Auto-Documentation)
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
   - Developer A: User Story 1 (Auto-Documentation)
   - Developer B: User Story 2 (Inheritance Discovery)
   - Developer C: User Story 3 (Framework Pivot)
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
- All MCP tools must follow stateless operation (NFR-001)
- Validation must support all three Definition layers (Layer 1, Layer 2, Layer 3)
- Import mechanism uses Gemara `imported-controls` field with MultiMapping structure
- Layer 3 Policy generation uses scope definition with Layer 1 and Layer 2 context
- Control IDs must use format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042"), must NOT include family, and are immutable once defined
- Auto-Documentation follows 5-phase sequential pipeline: (1) Capability Definition, (2) Threat Mapping, (3) Control Selection, (4) Audit Gap Analysis, (5) Verification

---

## Task Summary

- **Total Tasks**: 87
- **Setup Tasks**: 8 (Phase 1)
- **Foundational Tasks**: 17 (Phase 2)
- **User Story 1 Tasks**: 19 (Phase 3) - MVP (includes 5-phase pipeline orchestration, control ID validation, audit gap analysis)
- **User Story 2 Tasks**: 12 (Phase 4)
- **User Story 3 Tasks**: 16 (Phase 5)
- **Polish Tasks**: 15 (Phase 6) (includes control ID format validation tests)

### Parallel Opportunities Identified

- **Phase 1**: 2 parallel tasks (T007, T008)
- **Phase 2**: 15 parallel tasks (T010-T020, T022-T025)
- **Phase 3**: 10 parallel tasks (T026-T030, T032-T034)
- **Phase 4**: 3 parallel tasks (T045-T047, T048-T049)
- **Phase 5**: 5 parallel tasks (T057-T061, T062)
- **Phase 6**: 7 parallel tasks (T073-T079)

### Independent Test Criteria

- **User Story 1**: Provide technical evidence → receive validated Layer 2 Catalog artifact
- **User Story 2**: Provide partial catalog → receive inheritance suggestions → import controls
- **User Story 3**: Provide Layer 2 controls + regulatory requirements → receive gap analysis report + Layer 3 Policy

### Suggested MVP Scope

**MVP = Phase 1 + Phase 2 + Phase 3 (User Story 1 only)**

This delivers the core Auto-Documentation capability, enabling users to generate validated Layer 2 Catalog artifacts from technical evidence. This is the highest-priority user need and provides immediate value.
