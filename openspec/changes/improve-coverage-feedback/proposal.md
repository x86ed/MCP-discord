## Why

The GitHub Actions coverage check currently shows a failure message when coverage is between 90-95%, even though the workflow passes. This creates confusion as developers see a red "❌" message despite meeting the 90% threshold. Better messaging would differentiate between "good enough" (90-95%) and "excellent" (95%+) coverage.

## What Changes

- Update coverage threshold enforcement in `.github/workflows/test.yml` to use tiered messaging
- 90-94.99% coverage: pass with warning message (⚠️ yellow indicator)
- 95%+ coverage: pass with success message (✅ green indicator)
- Below 90%: fail with error message (❌ red indicator, exit 1)
- Remove duplicate threshold check logic in current workflow

## Capabilities

### New Capabilities
- `tiered-coverage-thresholds`: Implement multi-level coverage thresholds with appropriate messaging for each tier (fail < 90%, warn 90-95%, success >= 95%)

### Modified Capabilities
<!-- No existing capabilities being modified -->

## Impact

- Affected code: `.github/workflows/test.yml` (coverage job, enforce threshold step)
- User experience: Clearer feedback on coverage quality - developers will know if they're "good" vs "excellent"
- No breaking changes - workflow still enforces 90% minimum
- No impact on test execution or coverage calculation
