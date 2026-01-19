# Data Model: Gemara Artifact Authoring Assistant

**Date**: 2025-01-27  
**Feature**: 001-gemara-authoring-assistant  
**Status**: Phase 1 Design

## Design Approach

This data model uses a **dual approach** to align with Gemara's lexicon and schemas:

1. **Gemara-Native Schema Structure**: All entities conform to Gemara's native schema structures from the official schemas (https://github.com/gemaraproj/gemara/tree/main). This ensures structural compatibility and validation.

2. **CUE Schema Field Relationships**: All relationships are expressed through CUE schema structures (per Layer Structure Table in spec.md). These are CUE schema field names, not RDF predicates:
   - Threat "identifiedBy" Attack Pattern (Layer 2) - expressed via CUE schema field
   - Control "mitigates" Threat (Layer 2) - expressed via CUE schema field
   - Control "satisfies" Guideline (Layer 2) - expressed via CUE schema field
   - Guideline "establishes" Compliance Target (Layer 1) - expressed via CUE schema field

This dual approach ensures the codebase declaratively reflects Gemara's specification while maintaining both structural correctness (via native schemas) and semantic clarity (via CUE schema field relationships).

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
- Generates: Capability (one-to-many)
- Used in: Auto-Documentation Journey

**Validation Rules**:
- Content must not be empty
- Format must be supported by available parsers
- SourceType must be valid enum value

**State Transitions**: None (immutable input)

---

### Capability

**Description**: Extracted security capability identified from technical evidence. Aligned with Gemara lexicon capability definitions.

**Attributes**:
- `ID`: string - Unique identifier
- `Name`: string - Capability name (e.g., "Full-disk encryption via LUKS")
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

**Description**: Security threat identified from threat library. Uses Gemara-native Threat structure from schemas AND CUE schema field relationships. Relationships are expressed through CUE schema structures, not RDF triples.

**Attributes** (per Gemara-native Threat structure from schemas):
- `ID`: string - Threat identifier (from Gemara threat library)
- `Name`: string - Threat name (e.g., "Physical theft of storage media")
- `Description`: string - Threat description
- `ThreatCategory`: string - Category of threat
- `Layer1Reference`: string - Reference to Layer 1 guidance document
- `AffectedCapabilities`: []string - Capability types this threat affects

**CUE Schema Field Relationships** (expressed through CUE schema structures):
- `identifiedBy`: AttackPatternReference - CUE schema field: Threat "identifiedBy" Attack Pattern (MITRE ATT&CK). Expressed via CUE schema field, not RDF predicate. Per Gemara Layer 2 CUE schema model.
- `mitigatedBy`: []ControlReference - CUE schema field: Control "mitigates" Threat (reverse relationship). Expressed via CUE schema field, not RDF predicate. Per Gemara Layer 2 CUE schema model.

**Supporting Types**:
- `AttackPatternReference`: Contains Attack Pattern ID (MITRE ATT&CK) and relationship metadata per Gemara schema
- `ControlReference`: Contains Control ID and relationship metadata per Gemara schema

**Relationships** (additional context):
- Mapped from: Capability (many-to-many) - Capabilities map to threats via threat library queries
- Referenced in: Layer1Guidance (many-to-one)

**Validation Rules**:
- ID must reference valid threat from Gemara library
- identifiedBy.AttackPatternID must reference valid MITRE ATT&CK pattern
- Layer1Reference must reference valid Layer 1 guidance
- Structure must conform to Gemara Layer 2 Threat schema

**Source**: Loaded from Threat Catalog (Layer 2) via query_gemara_info with query_type="threat_catalog"

---

### Control

**Description**: Security control aligned with Gemara Layer 2 Control structure. Uses Gemara-native Control structure from schemas AND CUE schema field relationships. Relationships are expressed through CUE schema structures, not RDF triples.

**Attributes** (per Gemara-native Control structure from schemas):
- `ID`: string - Control identifier. Format: `<identifier>-<numbering>` (e.g., "AC-001", "SEC-042"). MUST NOT include family in ID. Control IDs are immutable once defined until withdrawn, even if controls are reclassified into different families.
- `Title`: string - Control title/name (e.g., "Mandatory TPM-backed encryption")
- `Objective`: string - Control objective/description
- `Description`: string - Detailed control description
- `Family`: string - Control family grouping
- `Status`: string - Status (proposed, validated, imported)

**CUE Schema Field Relationships** (expressed through CUE schema structures):
- `mitigates`: []ThreatReference - CUE schema field: Control "mitigates" Threat (SPDX Mitigation). Expressed via CUE schema field, not RDF predicate. Per Gemara Layer 2 CUE schema model. Threat references this control mitigates.
- `satisfies`: []GuidelineReference - CUE schema field: Control "satisfies" Guideline (Compliance Inheritance). Expressed via CUE schema field, not RDF predicate. Per Gemara Layer 2 CUE schema model. Guideline references this control satisfies for compliance inheritance.

**Supporting Types**:
- `ThreatReference`: Contains threat ID and relationship metadata per Gemara schema. Used in "mitigates" CUE schema field relationship.
- `GuidelineReference`: Contains guideline ID (Layer 1 reference), compliance target, and relationship metadata per Gemara schema. Used in "satisfies" CUE schema field relationship.

**Relationships** (additional context):
- Part of: Layer2Artifact (many-to-one)
- Inherited from: Layer2Catalog (many-to-one, optional)

**Validation Rules**:
- Title must not be empty
- Must mitigate at least one threat (mitigates array must not be empty)
- Satisfies array may reference Layer 1 guidelines for compliance inheritance
- Status must be valid enum value
- Structure must conform to Gemara Layer 2 Control schema
- CUE schema field relationships must be valid per Gemara CUE schema model

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

**Description**: Layer 1 guidance document aligned with Gemara Layer 1 GuidanceDocument structure. Uses Gemara-native GuidanceDocument structure from schemas AND CUE schema field relationships. Relationships are expressed through CUE schema structures, not RDF triples.

**Attributes** (per Gemara-native GuidanceDocument structure from schemas):
- `ID`: string - Guidance identifier
- `Title`: string - Guidance title
- `Description`: string - Guidance description
- `Framework`: string - Framework name (e.g., "NIST 800-53", "CRA")
- `Version`: string - Framework version
- `Guidelines`: []Guideline - Guidelines in this guidance document

**CUE Schema Field Relationships** (expressed through CUE schema structures):
- `establishes`: []ComplianceTargetReference - CUE schema field: Guideline "establishes" Compliance Target (CRA / NIST 800-53). Expressed via CUE schema field, not RDF predicate. Per Gemara Layer 1 CUE schema model. Compliance targets established by guidelines.
- `satisfiedBy`: []ControlReference - CUE schema field: Control "satisfies" Guideline (reverse relationship). Expressed via CUE schema field, not RDF predicate. Per Gemara Layer 2 CUE schema model. Controls that satisfy these guidelines.

**Supporting Types**:
- `Guideline`: Contains guideline ID, title, description per Gemara schema
- `ComplianceTargetReference`: Contains compliance target identifier and relationship metadata per Gemara schema. Used in "establishes" CUE schema field relationship.
- `ControlReference`: Contains Control ID and relationship metadata per Gemara schema. Used in "satisfiedBy" CUE schema field relationship.

**Relationships** (additional context):
- Referenced by: Threat (many-to-one)

**Validation Rules**:
- ID must be unique
- Framework must be valid
- Guidelines must not be empty
- Establishes array contains compliance targets per Gemara CUE schema model
- Structure must conform to Gemara Layer 1 GuidanceDocument schema
- CUE schema field relationships must be valid per Gemara CUE schema model

**Source**: Loaded from Gemara info storage (reference-only, not authored by this system)

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

### ThreatCatalog

**Description**: Layer 2 Threat Catalog storing threats that can be queried via threat library. Threat Catalogs are reference-only (query/storage), not authored by this system.

**Attributes**:
- `ID`: string - Threat Catalog identifier
- `Title`: string - Threat Catalog title
- `Description`: string - Threat Catalog description
- `Threats`: []Threat - List of threats in this catalog
- `Metadata`: map[string]string - Additional metadata
- `Source`: string - Source location (file path, URL, etc.)

**Relationships**:
- Contains: Threat (one-to-many)
- Queried by: query_threat_library tool (via query_gemara_info with query_type="threat_catalog")

**Validation Rules**:
- ID must be unique
- Must contain at least one threat
- Must conform to Gemara Layer 2 Threat Catalog schema

**Source**: Loaded from Gemara info storage (reference-only, not authored by this system)

---

### RiskCatalog

**Description**: Layer 3 Risk Catalog storing organization risks. Risk Catalogs link to threats and are mitigated by controls. Risk Catalogs are reference-only (query/storage), not authored by this system. Risk Catalogs are NOT directly used in Framework Pivot analysis since framework pivot only introduces non-technical noncompliance risk (not new threats).

**Attributes**:
- `ID`: string - Risk Catalog identifier
- `Title`: string - Risk Catalog title
- `Description`: string - Risk Catalog description
- `Risks`: []Risk - List of risks in this catalog
- `Metadata`: map[string]string - Additional metadata
- `Source`: string - Source location (file path, URL, etc.)

**Risk Attributes**:
- `ID`: string - Risk identifier
- `ThreatID`: string - Reference to threat from Threat Catalog
- `Impact`: string - Impact level (e.g., "high", "medium", "low")
- `Probability`: string - Probability level (e.g., "high", "medium", "low")
- `Description`: string - Risk description

**Relationships**:
- Contains: Risk (one-to-many)
- Links to: Threat (many-to-one via ThreatID)
- Queried by: query_gemara_info with query_type="risk_catalog"

**Validation Rules**:
- ID must be unique
- Must contain at least one risk
- Must conform to Gemara Layer 3 Risk Catalog schema
- Risk ThreatID must reference valid threat from Threat Catalog

**Source**: Loaded from Gemara info storage (reference-only, not authored by this system)

**Note**: Risk Catalogs focus on threat-based risk assessment (impact × probability). Framework Pivot introduces compliance gaps (noncompliance risk), not new threats, so Risk Catalogs are not directly used in Framework Pivot gap analysis.

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

### AttackPatternReference (per Gemara-native Threat structure)
- `AttackPatternID`: string - Reference to MITRE ATT&CK pattern ID
- `RelationshipMetadata`: map[string]interface{} - Additional relationship metadata per Gemara schema
- Used in: Threat "identifiedBy" Attack Pattern CUE schema field relationship

### ThreatReference (per Gemara-native Control structure)
- `ThreatID`: string - Reference to threat ID
- `RelationshipMetadata`: map[string]interface{} - Additional relationship metadata per Gemara schema
- Used in: Control "mitigates" Threat CUE schema field relationship

### GuidelineReference (per Gemara-native Control structure)
- `GuidelineID`: string - Reference to Layer 1 guideline ID
- `ComplianceTarget`: string - Compliance target established by guideline
- `RelationshipMetadata`: map[string]interface{} - Additional relationship metadata per Gemara schema
- Used in: Control "satisfies" Guideline CUE schema field relationship

### ComplianceTargetReference (per Gemara-native GuidanceDocument structure)
- `ComplianceTargetID`: string - Compliance target identifier
- `Description`: string - Compliance target description
- `Framework`: string - Framework source (e.g., "NIST 800-53", "CRA")
- `RelationshipMetadata`: map[string]interface{} - Additional relationship metadata per Gemara schema
- Used in: Guideline "establishes" Compliance Target CUE schema field relationship

### ControlReference (per Gemara-native CUE schema field relationships)
- `ControlID`: string - Reference to control ID
- `RelationshipMetadata`: map[string]interface{} - Additional relationship metadata per Gemara schema
- Used in: Control "satisfies" Guideline reverse relationship (satisfiedBy)

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

### MultiMapping (for imported-controls)
- `ReferenceID`: string - Reference ID pointing to source catalog/control (per Gemara `#MultiMapping` structure)
- `Entries`: []MappingEntry - Array of mapping entries, each with `reference-id` field
- Used in: `imported-controls`, `imported-threats`, `imported-capabilities` fields

---

## Data Flow Summary

### Auto-Documentation Journey
1. TechnicalEvidence → (Parse if obscure format) OR (Pass directly to LLM if common format) → Capability
2. Capabilities + ThreatCatalog (via query_gemara_info query_type="threat_catalog") → Threat (using Gemara-native structure AND "identifiedBy" Attack Pattern CUE schema field relationship)
3. Threat → LLM → Control (proposed, using Gemara-native structure AND "mitigates" Threat CUE schema field relationship)
4. Control + Layer1Guidance → Audit Gap Analysis (checks Control "satisfies" Guideline CUE schema field relationships against required compliance targets)
5. Control → CUEValidator → Layer2Artifact (validated, conforms to Gemara-native schema with CUE schema field relationships)

### Inheritance Discovery Journey
1. Partial Layer2Artifact + DependencyInfo → SearchQuery
2. SearchQuery → GemaraStorage → []Layer2Catalog (ranked)
3. Layer2Catalog → Import via Gemara `imported-controls` field (`[...#MultiMapping]`) → Layer2Artifact

### Framework Pivot Journey
1. Layer2Artifact + RegulatoryRequirement → Comparison
2. Comparison → GapAnalysisReport (noncompliance risk, not new threats)
3. GapAnalysisReport → Recommendations
4. Note: Risk Catalogs are NOT directly used since framework pivot only introduces noncompliance risk (not new threats)

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
