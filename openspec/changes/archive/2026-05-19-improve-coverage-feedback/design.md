## Context

The current CI workflow in `.github/workflows/test.yml` has a single coverage threshold (90%) with binary pass/fail logic. When coverage is between 90-95%, the workflow passes but shows a confusing failure message. This change improves messaging to provide clearer feedback about coverage quality.

Current implementation uses bash `bc` command to compare coverage against threshold and exits with 1 if below 90%.

## Goals / Non-Goals

**Goals:**
- Provide tiered coverage feedback (fail < 90%, warn 90-95%, success >= 95%)
- Maintain 90% minimum threshold enforcement
- Use GitHub Actions annotations for visual clarity
- Clean up duplicate threshold check logic

**Non-Goals:**
- Changing the 90% minimum threshold
- Modifying coverage calculation or test execution
- Adding coverage reporting to external services
- Changing which packages are included in coverage

## Decisions

### Decision 1: Use GitHub Actions Warning Annotations
Use `::warning::` annotations for 90-95% range instead of custom emoji-only output.

**Rationale:** GitHub Actions natively supports `::error::`, `::warning::`, and `::notice::` annotations that appear in the workflow UI. This provides better visibility than console output alone.

**Alternatives considered:**
- Custom emoji-only messages: Less discoverable, no workflow UI integration
- Comment on PR: Overkill for simple threshold messaging
- Status check contexts: Would require separate checks, more complex

### Decision 2: Three-Tier Threshold Logic
Implement three distinct ranges: < 90% (fail), 90-94.99% (warn), >= 95% (success).

**Rationale:** Provides clear differentiation between "meets minimum" and "exceeds expectations" while maintaining the existing 90% requirement.

**Alternatives considered:**
- Two tiers (pass/fail only): Doesn't address the problem
- More granular tiers (e.g., 85%, 90%, 95%, 98%): Over-engineered for current needs

### Decision 3: Single Pass Through Coverage Check
Consolidate the duplicate threshold checks into one logical block.

**Rationale:** Current workflow has duplicate logic (lines appear twice). Single block is clearer and prevents drift.

## Risks / Trade-offs

**Risk:** Warning annotations might be ignored by developers → **Mitigation:** Clear messaging in warning text, encourage aiming for 95%

**Trade-off:** Using `exit 0` for 90-95% range means workflow passes (as intended) but some teams might want stricter enforcement → **Mitigation:** Document that 95% is aspirational, 90% is required

**Risk:** bc command might not be available in future GitHub Actions runners → **Mitigation:** bc is part of standard ubuntu-latest; if removed, switch to awk comparison
