# Data Model: Gemara Artifact Authoring Assistant

**Date**: 2025-01-27  
**Feature**: 001-gemara-authoring-assistant  
**Status**: Phase 1 Design

## Entities

### Technical Evidence

**Description**: Raw documentation, configuration files, logs, or other unstructured sources containing security control information.

**Attributes**:
- `SourceType`: string - Type of evidence source (file, git_repo, document, log)
- `Content`: []byte - Raw content of the evidence
- `Format`: string - Format of the content (yaml, json, text, markdown, etc.)
- `Metadata`: map[string]string - Additional metadata (path, repository, timestamp, etc.)
- `ParserID`: string - Identifier for parser that should process this evidence

**Relationships**:
- Processed by: ConfigParser (many-to-one)
- Generates: SecurityFeature (one-to-many)
- Used in: Auto-Documentation Journey

**Validation Rules**:
- Content must not be empty
- Format must be supported by available parsers
- SourceType must be valid enum value

**State Transitions**: None (immutable input)

---

### SecurityFeature

**Description**: Extracted security capability or feature identified from technical evidence.

**Attributes**:
- `ID`: string - Unique identifier
- `Name`: string - Feature name (e.g., "Full-disk encryption via LUKS")
- `Description`: string - Detailed description
- `EvidenceRef`: string - Reference to source TechnicalEvidence
- `CapabilityType`: string - Type of capability (encryption, authentication, access_control, etc.)
- `Configuration`: map[string]interface{} - Configuration details extracted from evidence

**Relationships**:
- Extracted from: TechnicalEvidence (many-to-one)
- Maps to: Threat (many-to-many via ThreatMapping)
- Mitigated by: Control (many-to-many)

**Validation Rules**:
- Name must not be empty
- EvidenceRef must reference valid TechnicalEvidence
- CapabilityType must be valid enum value

---

### Threat

**Description**: Security threat identified from threat library that maps to capabilities.

**Attributes**:
- `ID`: string - Threat identifier (from Gemara threat library)
- `Name`: string - Threat name (e.g., "Physical theft of storage media")
- `Description`: string - Threat description
- `ThreatCategory`: string - Category of threat
- `Layer1Reference`: string - Reference to Layer 1 guidance document
- `AffectedCapabilities`: []string - Capability types this threat affects

**Relationships**:
- Mapped from: SecurityFeature (many-to-many via ThreatMapping)
- Mitigated by: Control (many-to-many)
- Referenced in: Layer1Guidance (many-to-one)

**Validation Rules**:
- ID must reference valid threat from Gemara library
- Layer1Reference must reference valid Layer 1 guidance

**Source**: Loaded from Gemara info storage (threat library)

---

### Control

**Description**: Security control proposed or identified to mitigate threats.

**Attributes**:
- `ID`: string - Control identifier. Format: `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042"). MUST NOT include family in ID. Control IDs are immutable once defined until withdrawn, even if controls are reclassified into different families.
- `Name`: string - Control name (e.g., "Mandatory TPM-backed encryption")
- `Description`: string - Control description
- `ControlType`: string - Type of control (preventive, detective, corrective)
- `MitigatesThreats`: []string - Threat IDs this control mitigates
- `Layer1Reference`: string - Reference to Layer 1 guidance (e.g., NIST 800-53 AC-3)
- `AuditMinimums`: []string - Required audit minimums from Layer 1
- `Status`: string - Status (proposed, validated, imported)

**Relationships**:
- Mitigates: Threat (many-to-many)
- References: Layer1Guidance (many-to-one)
- Part of: Layer2Artifact (many-to-one)
- Inherited from: Layer2Catalog (many-to-one, optional)

**Validation Rules**:
- Name must not be empty
- Must mitigate at least one threat
- Layer1Reference must reference valid Layer 1 guidance
- Status must be valid enum value

**State Transitions**:
- `proposed` → `validated` (after CUE validation passes)
- `proposed` → `rejected` (if validation fails)
- `imported` → `validated` (inherited controls)

---

### Layer2Artifact

**Description**: Generated Gemara-compliant control catalog containing structured control definitions.

**Attributes**:
- `ID`: string - Artifact identifier
- `Metadata`: ArtifactMetadata - Catalog metadata (name, version, description, etc.)
- `Controls`: []Control - List of controls in this artifact
- `ImportedControls`: []MultiMapping - List of imported control references using Gemara `imported-controls` field (`[...#MultiMapping]` with `@go(ImportedControls)` tag) for inherited controls
- `ValidationStatus`: string - Validation status (pending, valid, invalid)
- `ValidationErrors`: []ValidationError - List of validation errors if invalid
- `GemaraVersion`: string - Gemara specification version used for validation
- `GeneratedAt`: time.Time - Timestamp of generation
- `SourceEvidence`: []string - References to TechnicalEvidence used

**Relationships**:
- Contains: Control (one-to-many)
- Imports via: MultiMapping (one-to-many) - Uses Gemara `imported-controls` field (`[...#MultiMapping]` structure)
- Validated by: CUEValidator (many-to-one)
- Generated from: TechnicalEvidence (many-to-many)
- References: Layer1Guidance (many-to-many)

**Validation Rules**:
- Must conform to Gemara Layer 2 schema (CUE validation)
- Controls must have valid Layer1References
- Metadata must include required fields (name, version)

**State Transitions**:
- `pending` → `valid` (CUE validation passes)
- `pending` → `invalid` (CUE validation fails)
- `invalid` → `valid` (after corrections)

---

### Layer2Catalog

**Description**: Existing stored catalog of Layer 2 controls that can be searched and imported.

**Attributes**:
- `ID`: string - Catalog identifier
- `Name`: string - Catalog name
- `Description`: string - Catalog description
- `TechnologyDomain`: string - Technology domain (kubernetes, cloud, etc.)
- `ApplicabilityScope`: string - Scope of applicability
- `Controls`: []Control - Controls in this catalog
- `Metadata`: map[string]string - Additional metadata
- `Source`: string - Source location (file path, URL, etc.)

**Relationships**:
- Contains: Control (one-to-many)
- Imported into: Layer2Artifact (many-to-many)
- Searched in: Inheritance Discovery Journey

**Validation Rules**:
- ID must be unique
- Must contain at least one control
- TechnologyDomain used for search matching

---

### Layer1Guidance

**Description**: Layer 1 guidance document providing compliance requirement minimums.

**Attributes**:
- `ID`: string - Guidance identifier
- `Title`: string - Guidance title
- `Description`: string - Guidance description
- `Framework`: string - Framework name (e.g., "NIST 800-53")
- `Version`: string - Framework version
- `AuditMinimums`: []AuditMinimum - Required audit minimums
- `Threats`: []Threat - Threats referenced in this guidance

**Relationships**:
- Referenced by: Control (many-to-one)
- Referenced by: Threat (many-to-one)
- Used in: Auto-Documentation Journey (audit gap analysis)

**Validation Rules**:
- ID must be unique
- Framework must be valid
- AuditMinimums must not be empty

**Source**: Loaded from Gemara info storage

---

### Layer3Policy

**Description**: Generated Gemara-compliant policy document containing scope definitions, imports from Layer 1 and Layer 2, and adherence definitions.

**Attributes**:
- `ID`: string - Policy identifier
- `Title`: string - Policy title
- `Metadata`: PolicyMetadata - Policy metadata (id, description, author)
- `Contacts`: Contacts - Responsible and accountable contacts
- `Scope`: Scope - Scope definition (in: {})
- `Imports`: Imports - Imports from Layer 1 (GuidanceImport) and Layer 2 (CatalogImport)
- `Adherence`: Adherence - Adherence definitions

**Relationships**:
- Imports: Layer1Guidance (many-to-many via GuidanceImport)
- Imports: Layer2Catalog (many-to-many via CatalogImport)
- Generated from: Scope definition with Layer 1 and Layer 2 applicability queried for context

**Validation Rules**:
- Must conform to Gemara Layer 3 schema (CUE validation)
- Title must not be empty
- Scope must be defined
- Imports must reference valid Layer 1 and Layer 2 artifacts

**State Transitions**:
- `pending` → `valid` (CUE validation passes)
- `pending` → `invalid` (CUE validation fails)
- `invalid` → `valid` (after corrections)

**Generation Process**:
- Generated through scope definition
- Layer 1 guidance applicability queried to provide context
- Layer 2 catalog applicability queried to provide context

---

### RegulatoryRequirement

**Description**: Unstructured compliance requirement from standards, frameworks, or regulations.

**Attributes**:
- `ID`: string - Requirement identifier
- `Source`: string - Source document/location
- `Content`: string - Raw requirement text
- `ExtractedStructure`: map[string]interface{} - Structured extraction (if available)
- `Priority`: string - Priority level (critical, high, medium, low)
- `Status`: string - Status (covered, gap, partial)

**Relationships**:
- Compared against: Layer2Control (many-to-many)
- Part of: GapAnalysisReport (many-to-one)

**Validation Rules**:
- Content must not be empty
- Priority must be valid enum value
- Status must be valid enum value

---

### GapAnalysisReport

**Description**: Prioritized analysis comparing existing controls against regulatory requirements.

**Attributes**:
- `ID`: string - Report identifier
- `GeneratedAt`: time.Time - Generation timestamp
- `Framework`: string - Regulatory framework analyzed
- `CoveredRequirements`: []RegulatoryRequirement - Requirements with coverage
- `Gaps`: []Gap - Uncovered requirements
- `PartialCoverage`: []PartialCoverage - Requirements with partial coverage
- `Recommendations`: []Recommendation - Actionable recommendations
- `Confidence`: float64 - Confidence level (0.0-1.0)

**Relationships**:
- Analyzes: RegulatoryRequirement (one-to-many)
- Compares: Layer2Control (many-to-many)
- Generated in: Framework Pivot Journey

**Validation Rules**:
- Must contain at least one requirement
- Confidence must be between 0.0 and 1.0
- Gaps must be prioritized

---

## Supporting Types

### ArtifactMetadata
- `Name`: string
- `Version`: string
- `Description`: string
- `Author`: string
- `CreatedAt`: time.Time

### ValidationError
- `Path`: string - JSON path to error location
- `Message`: string - Error message
- `Severity`: string - Error severity (error, warning)

### AuditMinimum
- `ID`: string - Minimum identifier (e.g., "AC-3")
- `Description`: string - Minimum requirement description
- `Required`: bool - Whether this is required

### Gap
- `Requirement`: RegulatoryRequirement
- `Priority`: string
- `Reason`: string - Why this is a gap

### PartialCoverage
- `Requirement`: RegulatoryRequirement
- `CoveringControls`: []Control
- `CoveragePercentage`: float64
- `MissingAspects`: []string

### Recommendation
- `Type`: string - Recommendation type
- `Description`: string
- `Priority`: string
- `ActionItems`: []string

### ThreatMapping
- `SecurityFeatureID`: string
- `ThreatID`: string
- `Confidence`: float64
- `Rationale`: string

### MultiMapping (for imported-controls)
- `ReferenceID`: string - Reference ID pointing to source catalog/control (per Gemara `#MultiMapping` structure)
- `Entries`: []MappingEntry - Array of mapping entries, each with `reference-id` field
- Used in: `imported-controls`, `imported-threats`, `imported-capabilities` fields

---

## Data Flow Summary

### Auto-Documentation Journey
1. TechnicalEvidence → ConfigParser → SecurityFeature
2. SecurityFeature + ThreatLibrary → ThreatMapping → Threat
3. Threat → LLM → Control (proposed)
4. Control + Layer1Guidance → AuditMinimum validation
5. Control → CUEValidator → Layer2Artifact (validated)

### Inheritance Discovery Journey
1. Partial Layer2Artifact + DependencyInfo → SearchQuery
2. SearchQuery → GemaraStorage → []Layer2Catalog (ranked)
3. Layer2Catalog → Import via Gemara `imported-controls` field (`[...#MultiMapping]`) → Layer2Artifact

### Framework Pivot Journey
1. Layer2Artifact + RegulatoryRequirement → Comparison
2. Comparison → GapAnalysisReport
3. GapAnalysisReport → Recommendations

---

## Storage Considerations

**Stateless Operation** (NFR-001):
- All entities are request-scoped
- No persistence between requests
- Storage interface provides query capabilities but doesn't maintain state
- File-based storage loads artifacts on-demand per request

**Interface-Based Storage**:
- `GemaraStorage` interface enables different implementations
- File-based default for MVP
- Future: Database, object storage, etc.
