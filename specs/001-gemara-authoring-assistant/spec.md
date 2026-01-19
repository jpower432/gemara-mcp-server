# Feature Specification: Gemara Artifact Authoring Assistant

**Feature Branch**: `001-gemara-authoring-assistant`  
**Created**: 2025-01-27  
**Status**: Draft

## Scope

This MCP Server is dedicated to **Definition authoring** for Gemara artifacts. The Gemara framework consists of two main sections:
- **Definition Section** (Layers 1-3): Layer 1 (GuidanceDocument), Layer 2 (Catalog), Layer 3 (Policy) - **IN SCOPE**
- **Measurement Section** (Layers 4-6): Layer 4 (EvaluationLog) and higher layers - **OUT OF SCOPE**

The system supports authoring and validation of Layers 1-3 only. Measurement and evaluation capabilities (Layer 4+) are explicitly excluded from this feature.

**Authoring vs Query Scope**:
- **Layer 1 (GuidanceDocument)**: Reference-only (query/storage), not authored by this system
- **Layer 2 (Catalog)**: Full authoring support for Control Catalogs. Threat Catalog (Layer 2) is reference-only (query/storage), not authored by this system
- **Layer 3 (Policy)**: Full authoring support for Policy documents. Risk Catalog (Layer 3) is reference-only (query/storage), not authored by this system

## Understanding Gemara Layers

### Key Concepts

**Definition Section (Layers 1-3)** - **IN SCOPE for this MCP server**:

- **Layer 1 (Guidance)**: Guidelines establish compliance targets from external standards - providing the foundation for all other layers. These are reference-only artifacts that define "what" compliance means.
- **Layer 2 (Threats & Controls)**: 
  - **Threats** are identified by attack patterns (like from MITRE ATT&CK), mapping security concerns to standardized threat models
  - **Controls** serve dual purposes:
    - Mitigate threats - connecting security controls to threat models
    - Satisfy guidelines - connecting security controls to compliance requirements
  - This creates a bridge between security threats and compliance requirements, enabling threat-informed compliance
- **Layer 3 (Policy & Risk)**: 
  - **Policies** mandate controls based on organizational risk appetite, translating technical controls into business decisions
  - **Risks** are managed in the Policy determining what is accpeted and what is mitigated, determining which controls need to be implemented.
  - This layer connects technical controls to business context and organizational governance

**Measurement Section (Layers 4-6)** - **OUT OF SCOPE for this MCP server**:

- **Layer 4 (Evaluation)**: Activities in the Evaluation layer provide inspection of code, configurations, and deployments.
- **Layer 5 (Enforcement)**: Activities in the Enforcement layer provide prevention or remediation.
- **Layer 6 (Audit)**: Activities in the Audit layer provide a review of organizational policy and conformance.

**Note**: This MCP server focuses on **Definition authoring** (Layers 1-3). The Measurement layers (4-6) are out of scope but are included in this table for complete context of the Gemara framework.

### Gemara's Application to SDLC

Gemara integrates with the Software Development Lifecycle (SDLC) by connecting Definition artifacts (Policies, Risks, Controls) with Measurement activities (Evaluation, Enforcement, Audit) across the development pipeline. The following diagram illustrates how Gemara layers interact with SDLC stages:

```mermaid
graph TD
    subgraph SDLC [Software Development Lifecycle]
        Git((Git))
        CC[Code & Configuration]
        DP[Deployment Pipeline]
        LE[Live Environment]

        Git --> CC
        CC -->| | DP
        DP -->| | LE
    end

    subgraph Definition
        POL[Policies]
        RSK[Risks]
        CTR[Controls]

        POL -- Manage --> RSK
        POL -- Mandate --> CTR
        CTR -. Mitigate .-> RSK
    end

    subgraph Measurement
        IE[Intent Evaluation<br/>SAST, IAC Scans]
        BE[Behavioral Evaluation<br/>DAST, IAST]
        PE[Preventive Enforcement]
        AUD[Audit]
        
        subgraph CM [Continuous Monitoring]
            RE[Remediative Enforcement]
        end
    end

    %% Definition to SDLC
    POL -.->|Provides<br/>Design Requirements| SDLC

    %% Measurement to SDLC
    IE -.->|Reads| CC
    PE -.->|Gates| DP
    BE -.->|Reads Status &<br/>Simulates Actions| LE
    RE -.->|Patches via PR| CC
    RE -.->|Patches in place| LE

    %% Internal Logic
    POL -.->|Configure| IE
    POL -.->|Configure| BE
    IE -->|Provide Evidence| AUD
    IE -->|Provide Findings| PE
    BE -->|Telemetry| RE
    BE --> RE
    AUD -.->|Validates| POL

    %% Styling
    style SDLC fill:none,stroke:#333,stroke-width:2px
    style Definition fill:none,stroke:#333,stroke-dasharray: 5 5
    style Measurement fill:none,stroke:#333,stroke-dasharray: 5 5
    style POL stroke:#326ce5,stroke-width:2px
    style CTR stroke:#326ce5,
```

**Key Integration Points**:

1. **Definition → SDLC**: Policies provide design requirements that guide development practices, ensuring code and configuration align with organizational security and compliance standards.

2. **Intent Evaluation (Layer 4)**: 
   - Reads Code & Configuration from Git repositories
   - Uses SAST (Static Application Security Testing) and IAC (Infrastructure as Code) scans
   - Configured by Policies to check for compliance with Controls
   - Provides evidence for Audit and findings for Preventive Enforcement

3. **Preventive Enforcement (Layer 5)**:
   - Gates the Deployment Pipeline, blocking non-compliant code from reaching production
   - Receives findings from Intent Evaluation
   - Ensures only compliant artifacts proceed through the pipeline

4. **Behavioral Evaluation (Layer 5)**:
   - Reads status and simulates actions in the Live Environment
   - Uses DAST (Dynamic Application Security Testing) and IAST (Interactive Application Security Testing)
   - Configured by Policies to test Controls in runtime
   - Provides telemetry to Remediative Enforcement

5. **Remediative Enforcement (Layer 5 - Continuous Monitoring)**:
   - Patches via Pull Requests back to Code & Configuration (Git)
   - Patches in place in Live Environment for immediate remediation
   - Receives telemetry from Behavioral Evaluation
   - Operates continuously to maintain compliance

6. **Audit (Layer 6)**:
   - Receives evidence from Intent Evaluation
   - Validates Policies, ensuring they remain effective and aligned with organizational needs
   - Creates immutable records of compliance activities

**Relevance to This MCP Server**:

This MCP server focuses on **Definition authoring** (Layers 1-3), which provides the foundation for the SDLC integration shown above. The artifacts created by this server (Layer 2 Controls, Layer 3 Policies) configure and inform the Measurement activities (Layers 4-6) that integrate with the SDLC. While the Measurement layers are out of scope for this server, understanding how Definition artifacts connect to SDLC helps users create more effective and actionable artifacts.

## Clarifications

### Session 2025-01-27

- Q: Should the system support Layer 4 (EvaluationLog) authoring, or is Layer 4 out of scope? → A: Layer 4+ (Layers 4-6, Measurement section) are out of scope. This MCP server is dedicated to Definition authoring (Layers 1-3) only.
- Q: When should the system parse technical evidence vs pass it directly to the LLM? → A: Parse only obscure formats LLMs aren't trained on; pass common formats (YAML, JSON, Markdown, text) directly to LLM. Obscure formats are defined as formats where LLMs lack sufficient training data or where structured parsing significantly improves accuracy (e.g., Dockerfile syntax, Kubernetes manifest structures, specialized configuration languages).
- Q: What is the minimum required structure when generating a Layer 2 artifact? → A: Generate complete Catalog with title, Controls, and Families (Categories optional)
- Q: What import mechanism should be used for inherited controls? → A: Use Layer 2 field "imported-controls": [...#MultiMapping] @go(ImportedControls). Also support imported-threats and imported-capabilities as needed.
- Q: Should the system support authoring Layer 1 (GuidanceDocument) and Layer 3 (Policy) artifacts, or are they only referenced/queried? → A: Layer 3 (Policy) authoring is supported - can be generated through scope definition with Layer 1 and Layer 2 applicability queried to provide context. Layer 1 (GuidanceDocument) is reference-only (query/storage).
- Q: Which layers should validation support? → A: Validate all three layers (Layer 1 GuidanceDocument, Layer 2 Catalog, Layer 3 Policy)
- Q: What is the decision-making process for selecting capabilities, threats, controls, and mapping to Layer 1 guidelines in Auto-Documentation? → A: LLM-driven with server-provided context. LLM makes decisions using server-provided context (parsed capabilities, threat library queries, Layer 1 guidance). Server validates structure and provides deterministic validation.
- Q: What is the exact data flow sequence for Auto-Documentation (Layer 2 creation)? → A: Sequential pipeline: (1) Capability Definition - LLM uses technical evidence to define capabilities via FR-002 Config Parsing, (2) Threat Mapping - Server provides threat library to map capabilities to threats via FR-006 Info Tooling (using Gemara-native structure AND "identifiedBy" Attack Pattern CUE schema field), (3) Control Selection - LLM proposes controls to mitigate threats (using Gemara-native structure AND "mitigates" Threat CUE schema field), (4) Audit Gap Analysis - Server checks proposed controls against Layer 1 audit minimums using Control "satisfies" Guideline CUE schema field relationships (FR-001, FR-006), (5) Verification - CUE Schema Validator ensures artifact structure (FR-004) conforms to Gemara-native schema with CUE schema field relationships
- Q: What happens if LLM-proposed controls are missing required regulatory requirements during Audit Gap Analysis? → A: Flag gaps in final artifact. Server identifies missing requirements by checking which required compliance targets (from Layer 1 guidelines) lack satisfying controls via "satisfies" CUE schema field relationships. Gaps are included as recommendations in the artifact. LLM can refine controls in subsequent iteration if needed.
- Q: What happens when threat library query returns no threats for a capability, or multiple threats match with different Layer 1 references? → A: Return all matches with confidence scores. Threat library query returns all matching threats with confidence scores and Layer 1 references. LLM selects most appropriate threats. If no matches, return empty result with note; LLM can still propose controls based on capability analysis.
- Q: How are Families determined in Layer 2 Catalog generation? → A: LLM creates Families from controls. LLM groups related controls into Families based on control relationships and security domains during control selection phase.
- Q: What format should control IDs follow? → A: Control IDs MUST use format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042") matching regex pattern `^[A-Z0-9]+-[0-9]{3,}$` where identifier is uppercase alphanumeric and numbering is 3+ digits. Do NOT include family in control ID. Control IDs are immutable once defined until withdrawn, even if controls are reclassified into different families.
- Q: How should the system handle Gemara's Threat Catalog (Layer 2) and Risk Catalog (Layer 3) concepts? → A: Threat Catalog (Layer 2) stores threats queried via threat library; Risk Catalog (Layer 3) stores organization risks similar to a risk register - system supports querying both
- Q: Should the system support authoring Threat Catalogs (Layer 2) and Risk Catalogs (Layer 3), or only query them? → A: Query only (reference/storage) - Threat Catalogs and Risk Catalogs are not authored by this system
- Q: How should the query_threat_library tool access threats from Threat Catalog (Layer 2)? → A: query_threat_library queries Threat Catalog (Layer 2) via query_gemara_info with query_type="threat_catalog"
- Q: Should the Framework Pivot journey (User Story 3) query Risk Catalogs (Layer 3) when analyzing gaps? → A: Risk Catalogs link to threats with context of impact vs probability. Framework pivot only introduces noncompliance risk (not new threats), so Risk Catalogs are not directly used in Framework Pivot analysis
- Q: What terminology should be used for capabilities discovered from technical evidence? → A: Use "Capability" only, remove "SecurityFeature" terminology, align with Gemara lexicon capability definitions
- Q: How should threat mapping align with Gemara schemas? → A: Use Gemara-native Threat structure from schemas,
- Q: How should control mapping align with Gemara schemas? → A: Use Gemara-native Control structure from schemas.
- Q: How should compliance requirement mapping align with Gemara schemas? → A: Use Gemara-native mapping structure.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Auto-Documentation (Priority: P1)

A compliance engineer needs to document security controls for audit purposes. They have raw technical evidence (configuration files, logs, documentation) but lack structured Layer 2 claims that meet Gemara standards. The system follows a sequential pipeline: (1) LLM defines capabilities from technical evidence, (2) Server maps capabilities to threats via threat library, (3) LLM proposes controls to mitigate threats, (4) Server performs audit gap analysis against Layer 1 requirements, (5) Server validates artifact structure. The system generates a validated Layer 2 Gemara artifact with proper structure, relationships, and metadata.

**Why this priority**: Auto-documentation addresses the most time-consuming aspect of compliance work - converting unstructured evidence into structured, auditable artifacts. This directly reduces manual effort and enables faster audit readiness.

**Independent Test**: Can be fully tested by providing sample technical evidence (e.g., security configuration files, policy documents) and verifying that a valid Layer 2 Gemara artifact is generated with correct structure, required fields populated, and validation passing. The generated artifact can be independently reviewed and used for audit purposes.

**Acceptance Scenarios**:

1. **Given** a user provides raw technical evidence (configuration files, documentation, logs) accessible via Git repositories through external MCP servers, **When** they request Layer 2 artifact generation, **Then** the system produces a validated Layer 2 Gemara artifact with all required fields, proper relationships, and passes Gemara schema validation
2. **Given** technical evidence contains multiple security controls, **When** the system processes the evidence, **Then** it generates separate Layer 2 control entries for each distinct control identified
3. **Given** technical evidence references existing Layer 1 guidance documents that are accessible, **When** the system generates Layer 2 claims, **Then** it properly links the generated controls to the referenced Layer 1 guidance documents
4. **Given** generated Layer 2 artifact contains validation errors, **When** the user reviews the artifact, **Then** the system provides clear error messages indicating what needs correction
5. **Given** a user requests validation of a Gemara artifact, **When** they specify a Gemara specification version, **Then** the system validates the artifact against that version and provides a detailed validation report

---

### User Story 2 - Inheritance Discovery (Priority: P2)

A compliance engineer is documenting controls for a new system but suspects some controls may already exist in imported catalogs or inherited from parent systems. The system analyzes the current documentation context and suggests existing Layer 2 catalogs that contain relevant controls, reducing redundant documentation work.

**Why this priority**: Inheritance discovery prevents duplicate work and ensures consistency across related systems. It helps users leverage existing control documentation rather than recreating it, significantly reducing documentation effort.

**Independent Test**: Can be fully tested by providing a partial Layer 2 catalog or control context and verifying that the system suggests relevant existing Layer 2 catalogs with matching controls. The suggestions can be independently evaluated for relevance and accuracy without implementing other journeys.

**Acceptance Scenarios**:

1. **Given** a user is authoring Layer 2 controls for a specific technology or domain with access to existing Layer 2 catalogs and dependency information (SBOMs, architecture diagrams, CALM artifacts), **When** they request inheritance suggestions, **Then** the system identifies and suggests existing Layer 2 catalogs that contain applicable controls
2. **Given** multiple existing catalogs contain potentially relevant controls, **When** the system provides suggestions, **Then** it prioritizes suggestions by relevance and indicates which specific controls match
3. **Given** a user selects a suggested catalog to import, **When** they confirm the import, **Then** the system incorporates the inherited controls through the Gemara `imported-controls` field (`[...#MultiMapping]` structure) into their current catalog with proper attribution
4. **Given** no relevant existing catalogs are found, **When** the system searches for inheritance opportunities, **Then** it clearly indicates that no matches were found
5. **Given** dependency information (SBOMs, architecture diagrams, CALM artifacts) is available, **When** the system analyzes inheritance opportunities, **Then** it uses this information to identify relationships and suggest relevant catalogs

---

### User Story 3 - Framework Pivot (Priority: P3)

A compliance team needs to assess their current technology stack against a new regulatory framework that lacks structured documentation. The system analyzes the current Layer 2 controls and compares them against the unstructured regulatory requirements, producing a prioritized gap analysis report. Framework pivot only introduces noncompliance risk (not new threats), so Risk Catalogs are not directly used in this analysis.

**Why this priority**: Framework pivot enables organizations to adapt to new compliance requirements efficiently. While less frequent than documentation tasks, it's critical when organizations must demonstrate compliance with new standards or regulations.

**Independent Test**: Can be fully tested by providing existing Layer 2 controls and unstructured regulatory requirements, then verifying that the system produces a prioritized report identifying gaps, covered requirements, and minimum requirements that remain uncovered. The report can be independently reviewed for accuracy and completeness.

**Acceptance Scenarios**:

1. **Given** a user provides existing Layer 2 controls and unstructured regulatory requirements with access to stored Layer 2 catalogs, **When** they request a framework pivot analysis, **Then** the system produces a prioritized report identifying covered requirements, gaps, and uncovered minimum requirements
2. **Given** multiple regulatory requirements exist, **When** the system analyzes coverage, **Then** it prioritizes uncovered requirements by criticality and indicates which existing controls partially address requirements
3. **Given** regulatory requirements are ambiguous or unclear, **When** the system processes them, **Then** it flags areas needing clarification and provides best-effort analysis with confidence indicators
4. **Given** the analysis identifies gaps, **When** the user reviews the report, **Then** it provides actionable recommendations for addressing uncovered requirements

---

### Edge Cases

- What happens when technical evidence is incomplete or contradictory?
- How does the system handle technical evidence in multiple formats or languages?
- What happens when no existing Layer 2 catalogs match the inheritance search criteria?
- How does the system handle regulatory requirements that conflict with existing controls?
- What happens when generated artifacts exceed size or complexity limits?
- How does the system handle cases where technical evidence doesn't map to any known control patterns?
- What happens when regulatory requirements reference standards not available in the system?
- What happens when external MCP servers for Git repositories are unavailable or unreachable?
- How does the system handle cases where required Layer 1 documents are missing or inaccessible?
- What happens when dependency information (SBOMs, architecture diagrams, CALM artifacts) is incomplete or outdated?
- What happens when authentication fails or user credentials are invalid?
- How does the system handle cases where remote communication encryption cannot be established?
- What happens when the system cannot achieve 90% deterministic outcomes for artifact generation?
- How does the system handle cases where performance metrics cannot be exported?

## Dependencies

### External Dependencies

- **Git Repository Access**: Journey 1 (Auto-Documentation) requires access to Git repositories containing technical evidence. This dependency is fulfilled through external MCP servers (e.g., GitHub MCP Server, GitLab MCP Server) that provide Git repository access capabilities.

### Data Dependencies

- **Layer 1 Guidance Documents**: Journey 1 (Auto-Documentation) requires existing Layer 1 documents to provide compliance requirement minimums. These documents must be accessible to the system for reference during Layer 2 artifact generation.

- **Layer 2 Catalogs**: Journey 2 (Inheritance Discovery) and Journey 3 (Framework Pivot) require access to existing Layer 2 Catalogs for searching, comparison, and import operations.

- **Dependency Information**: Journey 2 (Inheritance Discovery) requires dependency information such as Software Bill of Materials (SBOMs), architecture diagrams, or FINOS CALM artifacts to identify inherited security controls and relationships between systems.

## Assumptions

- External MCP servers for Git repositories (GitHub MCP Server, GitLab MCP Server) are available and accessible when needed for Journey 1
- Layer 1 guidance documents exist and are accessible in a format the system can process
- Layer 2 catalogs are stored in a queryable format and contain sufficient metadata for searching and matching
- Dependency information (SBOMs, architecture diagrams, CALM artifacts) is available and current enough to provide meaningful inheritance suggestions
- Users have appropriate permissions to access Git repositories, Layer 1 documents, and Layer 2 catalogs
- Technical evidence in configuration files and code follows common patterns that can be parsed and analyzed
- Gemara schema versions are backward compatible or users can specify the appropriate version for validation
- System operates in environments that support both local development and remote cloud deployment
- Users authenticate through secure authentication mechanisms before accessing the system
- Performance monitoring infrastructure is available to collect exported metrics

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept raw technical evidence in common formats (YAML, JSON, text, Markdown, Dockerfile, Kubernetes manifests). For common formats (YAML, JSON, Markdown, text), evidence is passed directly to the LLM. For obscure formats (Dockerfile, Kubernetes manifests, or other formats LLMs aren't trained on), the system MUST parse and extract relevant security control information before providing to the LLM. **Obscure formats** are defined as formats where LLMs lack sufficient training data or where structured parsing significantly improves accuracy (e.g., Dockerfile syntax, Kubernetes manifest structures, specialized configuration languages).
- **FR-002**: System MUST parse technical evidence from configuration files and code to identify capabilities and control implementations, BUT ONLY for obscure formats that LLMs are not trained on. For common formats that LLMs can process directly (YAML, JSON, Markdown, plain text), the system MUST pass evidence directly to the LLM without parsing. **Obscure formats** are formats where LLMs lack sufficient training data or where structured parsing significantly improves accuracy (e.g., Dockerfile syntax, Kubernetes manifest structures, specialized configuration languages). Parsing MUST be implemented via extensible interface supporting multiple parser types (YAML, JSON, text parsers as minimum, with extensibility for additional formats). Parsed capabilities MUST be provided to the LLM as context for decision-making when parsing is performed. See FR-002a for complete Auto-Documentation pipeline context requirements.
- **FR-002a**: System MUST provide server-generated context to the LLM for decision-making in Auto-Documentation journey. Context MUST include: parsed capabilities from technical evidence (when parsing is performed for obscure formats), threat library query results (threats mapped to capabilities with Layer 1 references), and available Layer 1 guidance documents. For common formats (YAML, JSON, Markdown, text), LLM processes evidence directly; for obscure formats, parsed capabilities are provided as context. LLM uses this context to make decisions about capability selection, threat mapping, control selection, and Layer 1 guideline mapping. The Auto-Documentation data flow MUST follow this sequential pipeline: (1) Capability Definition - LLM uses technical evidence directly (common formats) or parsed capabilities (obscure formats via FR-002 Config Parsing) to define system capabilities, (2) Threat Mapping - Server provides threat library (via FR-006 Info Tooling) to map capabilities to threats, (3) Control Selection - LLM proposes controls to mitigate identified threats, (4) Audit Gap Analysis - Server checks proposed controls against Layer 1 audit minimums (FR-001, FR-006) to identify missing regulatory requirements. If gaps are found, server MUST flag them in the final artifact as gaps/recommendations (not block artifact generation). LLM can refine controls in subsequent iteration if needed, (5) Verification - CUE Schema Validator (FR-004) ensures the artifact correctly documents Capability-Threat-Control relationships.
- **FR-003**: System MUST generate Layer 2 Gemara artifacts that conform to Gemara schema validation requirements. Generated artifacts MUST be complete Catalog structures containing: required `title` field, `Controls` array, and `Families` array. `Categories` array is optional but recommended. Generated Catalogs MUST be valid Gemara Layer 2 schema instances. Artifact generation is LLM-driven using server-provided context; server validates structural correctness without modifying LLM-generated content (NFR-009). Families MUST be created by the LLM during the Control Selection phase (step 3 of Auto-Documentation pipeline) by grouping related controls based on control relationships and security domains. Control IDs MUST follow format `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042") matching regex pattern `^[A-Z0-9]+-[0-9]{3,}$` where identifier is uppercase alphanumeric and numbering is 3+ digits. Control IDs MUST NOT include family in the ID. Control IDs are immutable once defined until withdrawn, even if controls are reclassified into different families.
- **FR-004**: System MUST validate any produced Gemara artifacts against Gemara schema specifications and provide detailed validation reports. Validation MUST support all three Definition layers: Layer 1 (GuidanceDocument), Layer 2 (Catalog including Threat Catalog), and Layer 3 (Policy and Risk Catalog). Validation reports MUST include: validation status (valid/invalid), errors array (each with path, message, severity: error/warning), warnings array, and schema version used
- **FR-005**: System MUST allow users to specify the Gemara specification version for validation purposes
- **FR-006**: System MUST store and query Gemara information (Layer 1 guidance, Layer 2 catalogs including Threat Catalogs, Layer 3 policies and Risk Catalogs) for use across all journeys. Layer 1 (GuidanceDocument) artifacts are reference-only (query/storage), not authored by this system. Threat Catalog (Layer 2) stores threats that can be queried via threat library - Threat Catalogs are reference-only (query/storage), not authored by this system. Risk Catalog (Layer 3) stores organization risks similar to a risk register - Risk Catalogs are reference-only (query/storage), not authored by this system. The query_gemara_info tool MUST support query_type="threat_catalog" to query Threat Catalog (Layer 2) artifacts and query_type="risk_catalog" to query Risk Catalog (Layer 3) artifacts. When querying threat library (via query_threat_library tool), system MUST query Threat Catalog (Layer 2) via query_gemara_info with query_type="threat_catalog" and return all matching threats with confidence scores and Layer 1 references. If no threats match a capability, system MUST return empty result with explanatory note; LLM can still propose controls based on capability analysis. If multiple threats match, system MUST return all matches (not just first match) to enable LLM selection.
- **FR-006a**: System MUST support generation of Layer 3 (Policy) artifacts through scope definition. System MUST also support querying Layer 3 Risk Catalogs (which store organization risks similar to a risk register, linking to threats with impact vs probability context). Risk Catalogs are not directly used in Framework Pivot analysis since framework pivot only introduces noncompliance risk (not new threats). When generating Layer 3 policies, system MUST query Layer 1 guidance and Layer 2 catalog applicability to provide context for policy scope, imports, and adherence definitions
- **FR-007**: System MUST link generated Layer 2 controls to relevant Layer 1 guidance documents when references are identified. LLM selects Layer 1 mappings based on threat library query results (which include Layer 1 references) and available Layer 1 guidance context provided by the server. Server validates that referenced Layer 1 guidance exists and is accessible.
- **FR-008**: System MUST provide validation feedback when generated artifacts contain errors or missing required fields. Error messages MUST follow the same structure as FR-004 validation reports (path, message, severity) to ensure consistency
- **FR-009**: System MUST search existing Layer 2 catalogs to identify controls that match the current authoring context
- **FR-010**: System MUST suggest existing Layer 2 catalogs ranked by relevance to the current documentation context. Relevance ranking MUST prioritize: exact matches (technology domain + control types) > partial matches (technology domain OR control types) > related matches (similar domains or overlapping controls). Each suggestion MUST include relevance score and indicate which specific controls match
- **FR-011**: System MUST enable users to import suggested inherited controls into their current catalog using the Layer 2 `imported-controls` field structure (`imported-controls`: `[...#MultiMapping]` with `@go(ImportedControls)` tag). System MUST also support `imported-threats` and `imported-capabilities` fields for importing related entities. Imported controls MUST include proper attribution: source catalog ID (via `reference-id` in MultiMapping), import timestamp, original control ID, and source catalog metadata (name, version, description)
- **FR-012**: System MUST accept unstructured regulatory requirements (text documents, PDFs, web content) as input for framework pivot analysis. PDF parsing MUST extract text content for analysis; structured PDF metadata MAY be used if available. Web content MUST be extracted as text (HTML stripped) for processing
- **FR-013**: System MUST compare existing Layer 2 controls against unstructured regulatory requirements to identify coverage and gaps
- **FR-014**: System MUST produce prioritized reports indicating covered requirements, gaps, and uncovered minimum requirements. Reports MUST include: covered requirements list, gaps array (each with requirement, priority, reason), partial coverage array (requirement, covering controls, coverage percentage, missing aspects), recommendations array (type, description, priority, action items), and overall confidence level (float64, 0.0-1.0)
- **FR-015**: System MUST handle cases where regulatory requirements are ambiguous or incomplete. When ambiguity is detected, system MUST flag areas needing clarification, provide best-effort analysis with confidence indicators, and clearly indicate which aspects are uncertain
- **FR-016**: System MUST provide confidence indicators when analysis results have uncertainty. Confidence indicators MUST be expressed as float64 values between 0.0 (no confidence) and 1.0 (complete confidence), with clear documentation of what factors contribute to confidence calculation. Confidence factors include: (1) input quality (completeness and clarity of technical evidence or regulatory requirements), (2) threat match confidence (relevance score from threat library queries), (3) control coverage percentage (proportion of requirements covered by existing controls), (4) regulatory requirement clarity (ambiguity level of unstructured requirements), and (5) schema validation status (structural correctness of generated artifacts)

## Non-Functional Requirements

### Performance & Scalability

- **NFR-001**: System MUST operate without persisting data between requests (stateless operation)
- **NFR-002**: System MUST support both local development environments and remote cloud deployment modes
- **NFR-003**: System MUST ensure consistent artifact generation with at least 90% deterministic outcomes when processing the same input

### Security & Privacy

- **NFR-004**: All remote communications MUST be encrypted to prevent interception of mapping information
- **NFR-005**: System MUST authenticate users securely before processing requests
- **NFR-006**: System MUST ensure logical isolation of request-scoped data between different user sessions

### Reliability & Observability

- **NFR-007**: System MUST export performance metrics to enable monitoring and analysis
- **NFR-008**: System MUST track domain-specific metrics including artifact generation success rates and validation failure rates

### Data Integrity

- **NFR-009**: System MUST validate structural validity of artifacts without modifying LLM-generated content (input-output purity)

### Key Entities *(include if feature involves data)*

- **Technical Evidence**: Raw documentation, configuration files, logs, or other unstructured sources containing security control information. Attributes include source type, content, format, and metadata.
- **Layer 2 Artifact**: Generated Gemara-compliant control catalog containing structured control definitions, relationships, and metadata. Attributes include controls, validation status, and references to Layer 1 guidance.
- **Layer 2 Catalog**: Existing stored catalog of Layer 2 controls that can be searched and imported. Attributes include catalog identifier, controls contained, technology domain, and applicability scope. Includes Threat Catalog (Layer 2) which stores threats queried via threat library.
- **Layer 3 Policy**: Generated Gemara-compliant policy document containing scope definitions, imports from Layer 1 and Layer 2, and adherence definitions. Generated through scope definition with Layer 1 and Layer 2 applicability queried for context.
- **Layer 3 Risk Catalog**: Gemara-compliant risk catalog (Layer 3) that stores organization risks similar to a risk register. Risk Catalogs link to threats with context of impact vs probability. Risk Catalogs are reference-only (query/storage), not authored by this system. Risk Catalogs are not directly used in Framework Pivot analysis since framework pivot only introduces noncompliance risk (not new threats).
- **Regulatory Requirements**: Unstructured compliance requirements from standards, frameworks, or regulations. Attributes include source, content, and context.
- **Gap Analysis Report**: Prioritized analysis comparing existing controls against regulatory requirements. Attributes include covered requirements, gaps, priorities, and recommendations.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can generate a validated Layer 2 Gemara artifact from technical evidence in under 10 minutes
- **SC-002**: Generated Layer 2 artifacts pass Gemara schema validation on first attempt in 90% of cases
- **SC-003**: System identifies relevant existing Layer 2 catalogs for inheritance in 80% of searches
- **SC-004**: Users reduce redundant control documentation by 50% when using inheritance discovery
- **SC-005**: System produces framework pivot gap analysis reports in under 15 minutes
- **SC-006**: Gap analysis reports accurately identify 85% of uncovered requirements when compared to manual expert review
- **SC-007**: Users successfully complete at least one of the three journeys (auto-documentation, inheritance discovery, or framework pivot) on first attempt without training
- **SC-008**: System generates consistent artifacts (at least 90% deterministic) when processing the same technical evidence multiple times
- **SC-009**: System processes requests without requiring data persistence between requests (100% stateless operation)
- **SC-010**: Performance metrics are available for monitoring artifact generation success rates and validation failures
