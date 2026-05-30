## Context

The project currently has:
- Test coverage requirements with 90% minimum and 95% target thresholds
- Multiple Go versions in use: matrix testing on 1.21-1.23, but go.mod declares 1.25.3 and Dockerfile uses 1.25-alpine
- Terraform infrastructure in `/infra/` directory with no automated validation in CI
- GitHub Actions workflows for testing (`test.yml`) and deployment (`deploy.yml`)

This design addresses three production-readiness gaps: enforcing stricter coverage requirements, standardizing toolchain versions, and validating infrastructure configuration automatically.

## Goals / Non-Goals

**Goals:**
- Raise minimum test coverage from 90% to 95% by modifying existing threshold checks
- Standardize all Go version references to 1.26 (current stable release with security support)
- Add Terraform validation to GitHub Actions to catch infrastructure errors before deployment
- Maintain backward compatibility with existing test and deployment workflows

**Non-Goals:**
- Changing test scope or coverage calculation methodology
- Migrating to newer Terraform versions or providers
- Restructuring GitHub Actions workflow architecture
- Adding new test frameworks or coverage tools

## Decisions

### Decision 1: Use Go 1.26 as the standard version
**Rationale:** Go 1.26 is the current stable release with active security support. Go follows a 2-version support policy, so with 1.25 and 1.26 available, older versions like 1.23 no longer receive security updates. Currently, go.mod incorrectly declares 1.25.3 (non-existent version) and Dockerfile uses 1.25-alpine. Standardizing on 1.26 ensures security patches and consistency across development, CI, and production.

**Alternatives considered:**
- Use Go 1.23: Rejected because it's past end-of-life and won't receive security patches
- Use Go 1.25: Rejected because 1.26 is the current stable release with latest fixes
- Keep matrix testing multiple versions: Rejected because it adds CI complexity and doesn't reflect production reality (Docker image uses single version)

**Implementation:**
- Update `go.mod` to `go 1.26`
- Update `Dockerfile` FROM to `golang:1.26-alpine`
- Update `.github/workflows/test.yml` matrix to test only on `1.26` (single version)
- Update `.github/workflows/test.yml` coverage job to explicitly use `1.26`

### Decision 2: Enforce 95% as hard minimum in coverage job
**Rationale:** Current implementation has tiered thresholds (90% minimum, 95% target) with different exit codes. To enforce 95% as production requirement, modify the threshold check logic to fail on coverage < 95%.

**Alternatives considered:**
- Add separate "production-ready" workflow: Rejected because it fragments coverage validation
- Use branch protection rules to enforce 95%: Rejected because workflow should be self-documenting

**Implementation:**
- In `.github/workflows/test.yml` coverage job, change `MIN_THRESHOLD` from 90.0 to 95.0
- Update messaging: error if < 95%, success if >= 95%
- Remove warning tier (90-95%) since 95% is now minimum

### Decision 3: Add Terraform validation as new job in test.yml
**Rationale:** Infrastructure validation should run on every PR/push to catch configuration errors early. Adding as separate job in test.yml keeps validation close to code tests.

**Alternatives considered:**
- Add to deploy.yml: Rejected because deploy runs only on main branch, too late to catch errors
- Create separate terraform.yml workflow: Rejected to avoid workflow fragmentation
- Use terraform plan: Rejected because it requires AWS credentials; validation only needs syntax/consistency checks

**Implementation:**
- Add new `terraform-validate` job to `.github/workflows/test.yml`
- Install Terraform CLI via hashicorp/setup-terraform action
- Run `terraform init` and `terraform validate` in `/infra` directory
- Run on same triggers as other test jobs (push/PR to main/develop branches)

## Risks / Trade-offs

**Risk:** Raising coverage to 95% may block PRs if existing code doesn't meet threshold
→ **Mitigation:** This is intentional - production-readiness gate. Contributors must add tests before merging.

**Risk:** Removing matrix testing (1.21, 1.22) means we don't validate backward compatibility
→ **Mitigation:** Acceptable trade-off. Project uses Docker for deployment, so only runtime version (1.23) matters. Development uses same version via go.mod.

**Risk:** Terraform validate won't catch runtime/credential issues
→ **Mitigation:** Expected. Validation checks syntax and consistency only. Deployment workflow still handles AWS integration testing.

**Risk:** Changing go.mod version may break local development for contributors on older Go versions
→ **Mitigation:** Go 1.26 is the current stable release (May 2026), widely available. Contributors can update via package managers or `go install golang.org/dl/go1.26@latest`.
