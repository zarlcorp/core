# 109: Add PR review agent to delegation pipeline

## Objective
Update the `/delegate` command to insert an automated PR review step between PR creation and merge. Currently, the git-workflow-manager agent creates a PR and immediately merges it. Instead, a review agent should evaluate the PR against the spec before deciding to merge or request changes.

## Context
The current pipeline is: work agent → gitops agent (commit + push + PR + merge). This auto-merges without review, which means code quality issues or spec deviations slip through. Adding a review step creates a proper gate.

## Requirements

### Update delegate.md Step 6: gitops (no merge)
- Remove step 7 (`gh pr merge`) from the git-workflow-manager prompt
- The gitops agent should stop after creating the PR and commenting on the issue
- It should report the PR number and URL back

### Add delegate.md Step 7: PR review agent
After the gitops agent completes and reports the PR URL, launch a review agent:

**Review agent responsibilities:**
1. Read the spec from `.manager/specs/<id>-<name>.md`
2. Read the PR diff via `gh pr diff <pr-number> --repo <target-repo>`
3. Read the full files changed (not just diff) for context
4. **Spec compliance check** — for each requirement in the spec:
   - Is it implemented?
   - Does it meet acceptance criteria?
   - Is it within scope (no unrelated changes)?
5. **Code quality check** against CLAUDE.md standards:
   - Error handling (no "failed to" prefixes)
   - Naming conventions (scope-based)
   - Early returns over if/else chains
   - No branch duplication
   - Test coverage for new code
   - No mocking when fakes work
6. **Decision:**
   - **Approve**: comment review summary on PR, then merge (`gh pr merge --squash --delete-branch`)
   - **Request changes**: comment detailed feedback on PR, do NOT merge, report to PM

### Review agent prompt template
```
You are a code review agent. Review PR #<pr-number> on <target-repo> against the spec.

Spec file: <spec-path>
Target repo: zarlcorp/<repo-name>
PR number: <pr-number>
GitHub issue: #<issue-number>

Steps:
1. Read the spec to understand requirements
2. Read the PR diff: gh pr diff <pr-number> --repo zarlcorp/<repo-name>
3. For each file changed, read the full file for context
4. Check every spec requirement is implemented
5. Check code quality against project standards (CLAUDE.md)
6. If approved:
   - gh pr comment <pr-number> --repo zarlcorp/<repo-name> --body "## Review: Approved\n\n<summary>"
   - gh pr merge <pr-number> --repo zarlcorp/<repo-name> --squash --delete-branch
7. If changes needed:
   - gh pr comment <pr-number> --repo zarlcorp/<repo-name> --body "## Review: Changes Requested\n\n<detailed feedback>"
   - Do NOT merge
   - Report what needs fixing

Review criteria:
- Every spec requirement must be implemented
- Error messages: no "failed to", "unable to", "could not" prefixes
- Naming: scope-based (short names for small scope, descriptive for large)
- Code quality: early returns, no branch duplication, extract common ops
- Tests: table-driven preferred, real implementations over mocks
- No unrelated changes or scope creep

Do NOT add Co-Authored-By lines to commits.
Report the outcome (approved + merged, or changes requested + details).
```

### Update delegate.md Step 8 (renumbered): Report
- If review approved: report merge complete
- If review rejected: report what needs fixing, suggest re-delegation with fixes

### Review agent type
Use `subagent_type: "general-purpose"` since the review agent needs to:
- Read files (spec, source code)
- Run `gh` commands (pr diff, pr comment, pr merge)
- Make judgement calls about code quality

## Target Repo
zarlcorp/core

## Agent Role
PM task — update command files directly

## Files to Modify
- .claude/commands/delegate.md (modify — split Step 6, add Step 7 review, renumber Step 8)

## Notes
This is a PM workflow change, not application code. The PM can update delegate.md directly. The review agent should be strict but practical — don't block on style nitpicks, focus on spec compliance and correctness. The review agent runs as a background task like the gitops agent, so the PM can continue monitoring other work.
