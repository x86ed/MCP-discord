## ADDED Requirements

### Requirement: Repository provides a default pull request template
The repository SHALL define a default pull request template that GitHub automatically loads when contributors open a new pull request.

#### Scenario: New pull request starts with template content
- **WHEN** a contributor opens a pull request against the repository
- **THEN** the pull request body includes the repository-defined default template

### Requirement: Template captures review-critical change context
The pull request template MUST require structured fields for summary, motivation/context, linked tracking references, and testing evidence so reviewers can assess risk and correctness.

#### Scenario: Author fills required context sections
- **WHEN** a contributor prepares a pull request description
- **THEN** the template presents explicit sections for summary, context, references, and testing that the contributor can complete

### Requirement: Template supports OpenSpec-linked contributions
The pull request template SHALL include guidance for linking relevant OpenSpec change artifacts when a contribution originates from an OpenSpec workflow.

#### Scenario: OpenSpec contribution includes change linkage
- **WHEN** a contributor submits a pull request based on an OpenSpec change
- **THEN** the template includes a dedicated prompt to reference the related `openspec/changes/<name>/` artifacts

### Requirement: Template includes readiness checklist
The pull request template MUST include checklist items that confirm testing completion and reviewer-readiness declarations.

#### Scenario: Contributor confirms readiness before review
- **WHEN** a contributor finalizes a pull request
- **THEN** the template provides checklist items for test execution and readiness confirmations
