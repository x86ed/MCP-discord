## Why

Before deploying to production, we need robust quality gates to ensure code reliability. Currently, test coverage requirements allow production deployment with as low as 90% coverage, and infrastructure validation is not automated in CI. This change establishes stricter production-readiness criteria.

## What Changes

- Raise minimum test coverage requirement from 90% to 95%
- Standardize Go version across all workflows and build artifacts to use current stable (1.26)
- Add Terraform validation step to GitHub Actions workflow to catch infrastructure config errors early
- Update coverage threshold messaging to reflect 95% as minimum for production

## Capabilities

### New Capabilities
None - all changes enhance existing capabilities.

### Modified Capabilities
- `tiered-coverage-thresholds`: Raise minimum threshold from 90% to 95% for production readiness
- `cicd-pipeline`: Add Terraform validation step and standardize Go version to 1.26 across all workflow jobs

## Impact

- **Tests**: Coverage threshold enforcement becomes stricter - existing code must reach 95% coverage before merging
- **CI/CD**: GitHub Actions workflow will run Terraform validation on every PR and push
- **Build**: Dockerfile and GitHub Actions workflows will use consistent Go 1.26 version
- **Developer Experience**: Clearer quality expectations before production deployment
