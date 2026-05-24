## Why

Contributors currently have no repository-standard pull request template, which leads to inconsistent PR descriptions and missing review context. Adding a template now improves review quality, keeps maintenance overhead low, and aligns contribution workflow with the project's growing OpenSpec and testing practices.

## What Changes

- Add a repository-level pull request template in GitHub's recognized location.
- Define required PR sections for summary, linked change/spec context, testing evidence, and checklist items.
- Include guidance for OpenSpec-based work so contributors link relevant change artifacts when applicable.
- Keep the template lightweight and markdown-only with no runtime/code behavior impact.

## Capabilities

### New Capabilities
- `pull-request-template`: Standardizes pull request authoring requirements and review-ready metadata for all contributions.

### Modified Capabilities
- None.

## Impact

- Affected code: Repository metadata and documentation conventions only (no production Go code changes).
- APIs: None.
- Dependencies: None.
- Systems: GitHub pull request authoring and maintainer review workflow.
