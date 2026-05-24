## Context

The repository currently lacks a standardized pull request template, so PR descriptions vary and often omit information needed for efficient review (scope, test evidence, and rollout impact). GitHub supports repository-level templates via markdown files in recognized locations, which makes this a low-risk metadata change with immediate workflow benefits.

## Goals / Non-Goals

**Goals:**
- Provide a single default PR template that appears automatically for new pull requests.
- Require core review context (summary, related issue or OpenSpec change, testing evidence, and checklist confirmations).
- Keep contributor friction low with concise, clearly labeled sections and examples.

**Non-Goals:**
- Enforce PR quality via bots or CI checks in this change.
- Introduce multiple template variants (feature, bugfix, chore) in this change.
- Modify runtime code paths, binaries, or release automation.

## Decisions

1. Use a repository-level default template file.
- Decision: Add `.github/pull_request_template.md`.
- Rationale: This is natively recognized by GitHub and applies broadly with minimal maintenance.
- Alternative considered: `.github/PULL_REQUEST_TEMPLATE/*.md` for multiple templates. Rejected for now to avoid premature complexity.

2. Standardize mandatory author-provided sections.
- Decision: Include sections for summary, motivation/context, linked issue/OpenSpec reference, testing, and reviewer checklist.
- Rationale: These fields directly reduce back-and-forth during review and make changes auditable.
- Alternative considered: Very short freeform template. Rejected because it does not reliably improve review signal.

3. Include OpenSpec-aware guidance without making it universally required.
- Decision: Add a prompt to link `openspec/changes/<name>/` artifacts when the PR originates from an OpenSpec change.
- Rationale: Supports this repository's workflow while remaining compatible with non-OpenSpec maintenance PRs.
- Alternative considered: Hard requirement for all PRs. Rejected because not all contributions map to an OpenSpec change.

## Risks / Trade-offs

- [Template is too long and gets ignored] -> Keep sections concise and use checkboxes plus placeholders.
- [Template is too strict for small fixes] -> Include N/A-friendly prompts and avoid hard automation in this change.
- [Future workflow drift] -> Keep the file easy to revise and reference it in contributor docs in a follow-up if needed.

## Migration Plan

1. Add `.github/pull_request_template.md` with the agreed sections.
2. Open a test PR to verify GitHub auto-load behavior and formatting.
3. Adjust wording if maintainers identify friction during first usage.

Rollback strategy:
- Remove or simplify `.github/pull_request_template.md`; GitHub immediately returns to freeform PR bodies.

## Open Questions

- Should a follow-up add dedicated template variants in `.github/PULL_REQUEST_TEMPLATE/`?
- Should CI later validate that PR bodies include minimum metadata fields?
