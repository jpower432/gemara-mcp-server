# Quick Start: Gemara Artifact Authoring Assistant

**Feature**: 001-gemara-authoring-assistant  
**Date**: 2025-01-27

## Overview

The Gemara Artifact Authoring Assistant enables LLMs to automate the creation, validation, and management of Gemara security artifacts through MCP tools. This guide demonstrates the three critical user journeys.

## Prerequisites

- Gemara MCP Server running (local stdio or remote HTTP)
- Access to external MCP servers (GitHub/GitLab) for technical evidence
- Layer 1 guidance documents available in Gemara info storage
- Layer 2 catalogs available for inheritance discovery

## Journey 1: Auto-Documentation

**Goal**: Generate a validated Layer 2 Gemara artifact from raw technical evidence.

### Step 1: Parse Technical Evidence

```yaml
Tool: parse_technical_evidence
Input:
  source: "kubernetes/deployment.yaml"
  content: |
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: secure-app
    spec:
      template:
        spec:
          securityContext:
            runAsNonRoot: true
            runAsUser: 1000
  format: "yaml"
```

**Expected Output**: List of capabilities extracted (e.g., "Non-root user execution", "Security context enforcement")

### Step 2: Map Capabilities to Threats

```yaml
Tool: query_threat_library
Input:
  capability_types: ["access_control", "runtime_security"]
```

**Expected Output**: Threats that map to the capabilities (e.g., "Privilege escalation", "Container breakout")

### Step 3: Generate and Validate Artifact

```yaml
Tool: generate_layer2_artifact
Input:
  technical_evidence:
    - source: "kubernetes/deployment.yaml"
      content: "..."
      format: "yaml"
  artifact_metadata:
    name: "Secure App Controls"
    version: "1.0.0"
    description: "Security controls for secure-app deployment"
  gemara_version: "1.0.0"
```

**Expected Output**: Validated Layer 2 artifact with controls, validation status, and audit gap analysis.

### Validation Check

```yaml
Tool: validate_gemara_artifact
Input:
  artifact: { ... generated artifact ... }
  layer: 2
  gemara_version: "1.0.0"
```

**Expected Output**: Validation report with any errors or warnings.

---

## Journey 2: Inheritance Discovery

**Goal**: Identify existing Layer 2 catalogs that contain controls you can inherit.

### Step 1: Search for Inheritance Opportunities

```yaml
Tool: search_inheritance_opportunities
Input:
  current_context:
    technology_domain: "kubernetes"
    control_types: ["access_control", "network_policy"]
    partial_controls:
      - name: "Pod Security Policy"
        type: "access_control"
  dependency_info:
    sboms:
      - "sbom.json"
    calm_artifacts:
      - "architecture.calm.yaml"
  max_results: 5
```

**Expected Output**: Ranked list of suggested Layer 2 catalogs with relevance scores and matching controls.

### Step 2: Import Inherited Controls

```yaml
Tool: import_inherited_controls
Input:
  catalog_id: "kubernetes-base-controls-v1.2"
  control_ids: ["ctrl-001", "ctrl-002"]  # Optional: specific controls
  target_artifact_id: "my-artifact-001"
```

**Expected Output**: Confirmation of imported controls with attribution.

---

## Journey 3: Framework Pivot

**Goal**: Assess your current technology against a new regulatory framework.

### Step 1: Analyze Framework Coverage

```yaml
Tool: analyze_framework_pivot
Input:
  layer2_controls:
    - id: "ctrl-001"
      name: "Encryption at Rest"
      layer1_reference: "NIST-800-53:SC-28"
    - id: "ctrl-002"
      name: "Access Control"
      layer1_reference: "NIST-800-53:AC-3"
  regulatory_requirements:
    source: "New Regulation v2.0"
    content: |
      Section 4.1: All systems MUST implement encryption for data at rest.
      Section 4.2: Access controls MUST be enforced at all system boundaries.
      Section 4.3: Systems MUST implement key rotation policies.
    format: "text"
  framework_name: "New Regulation v2.0"
  include_recommendations: true
```

**Expected Output**: Gap analysis report with:
- Covered requirements (Section 4.1, 4.2)
- Gaps (Section 4.3 - key rotation)
- Partial coverage (if any)
- Prioritized recommendations

---

## Common Patterns

### Querying Gemara Information

```yaml
Tool: query_gemara_info
Input:
  query_type: "layer1"
  layer1_id: "nist-800-53-v5"
```

```yaml
Tool: query_gemara_info
Input:
  query_type: "search_layer2"
  search_query:
    technology_domain: "kubernetes"
    keywords: ["encryption", "access control"]
```

### Validating Artifacts

```yaml
Tool: validate_gemara_artifact
Input:
  artifact: { ... your artifact ... }
  layer: 2
  gemara_version: "1.0.0"  # Optional: use latest if not specified
```

---

## Error Handling

### Validation Errors

If validation fails, the `validate_gemara_artifact` tool returns detailed errors:

```yaml
Output:
  valid: false
  errors:
    - path: "/controls/0/layer1_reference"
      message: "Invalid Layer 1 reference format"
      severity: "error"
```

Fix errors and re-validate until `valid: true`.

### Parser Errors

If evidence parsing fails:

```yaml
Output:
  features: []
  parser_used: "yaml_parser"
  errors:
    - "Unable to parse YAML: invalid syntax at line 5"
```

Check evidence format and content, then retry.

### No Inheritance Matches

If no inheritance opportunities found:

```yaml
Output:
  suggestions: []
  total_found: 0
```

This is normal - proceed with manual control authoring.

---

## Best Practices

1. **Always Validate**: Use `validate_gemara_artifact` before considering an artifact complete
2. **Check Audit Gaps**: Review audit gap analysis from auto-documentation to ensure compliance
3. **Leverage Inheritance**: Search for inheritance opportunities before authoring new controls
4. **Version Awareness**: Specify Gemara version for validation to ensure consistency
5. **Iterative Refinement**: Use validation errors to iteratively improve artifacts

---

## Next Steps

- Review generated artifacts for accuracy
- Integrate with your compliance workflow
- Extend config parsers for your specific evidence types
- Customize threat library mappings for your domain
