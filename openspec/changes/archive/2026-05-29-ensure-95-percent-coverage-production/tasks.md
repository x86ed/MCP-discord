## 1. Standardize Go Version to 1.26

- [x] 1.1 Update go.mod to declare `go 1.26` (currently 1.25.3)
- [x] 1.2 Update Dockerfile FROM clause to `golang:1.26-alpine` (currently 1.25-alpine)
- [x] 1.3 Update test.yml matrix to use single Go version `1.26` (remove 1.21, 1.22)
- [x] 1.4 Verify coverage job in test.yml uses Go 1.26 (already set, confirm unchanged)

## 2. Raise Coverage Threshold to 95%

- [x] 2.1 Update MIN_THRESHOLD variable in test.yml coverage job from `90.0` to `95.0`
- [x] 2.2 Remove TARGET_THRESHOLD variable and warning tier logic (90-95% tier no longer needed)
- [x] 2.3 Update failure message to indicate 95% required for production deployment
- [x] 2.4 Update success message to indicate production-ready status at 95%

## 3. Add Terraform Validation to CI

- [x] 3.1 Add new `terraform-validate` job to .github/workflows/test.yml
- [x] 3.2 Configure job to run on push to main/develop and on pull_request
- [x] 3.3 Add hashicorp/setup-terraform action step to install Terraform CLI
- [x] 3.4 Add step to run `terraform init` in infra/ directory
- [x] 3.5 Add step to run `terraform validate` in infra/ directory with failure on errors

## 4. Verification

- [x] 4.1 Run tests locally to ensure coverage calculation still works correctly
- [x] 4.2 Verify go.mod consistency: run `go mod tidy` with Go 1.26
- [x] 4.3 Test Docker build with updated golang:1.26-alpine base image
- [x] 4.4 Verify Terraform validate runs successfully: `cd infra && terraform init && terraform validate`
