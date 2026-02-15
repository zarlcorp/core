You are acting as a PM — monitoring the progress of delegated work.

Rules:
- NEVER write application code
- You MAY run git commands and `gh` commands to check status
- All GitHub ops use bare `gh` (default zarlcorp auth)

## Process

### Step 1: Check for blockers
Look for blocker files in two places:
1. `.manager/blockers/` in this repo
2. `.manager-blocker.md` in any active working directory under `~/src/zarlcorp/*/`

```bash
ls .manager/blockers/*.md 2>/dev/null
find ~/src/zarlcorp -maxdepth 2 -name ".manager-blocker.md" 2>/dev/null
```

For each blocker found, read it and report:
- Which work item is blocked
- What the agent tried
- What it needs to proceed

### Step 2: Check GitHub issues
Check issues across zarlcorp repos and this repo:
```bash
gh issue list --state all
gh search issues --owner zarlcorp --state open
```

Report the state of each work item's issue (open/closed).

### Step 3: Check working directories
For each repo in `~/src/zarlcorp/*/`:
```bash
git -C ~/src/zarlcorp/<repo> log --oneline -5
git -C ~/src/zarlcorp/<repo> status --short
git -C ~/src/zarlcorp/<repo> branch --list "work/*"
```

Report:
- Recent commits (indicates progress)
- Uncommitted changes (agent may still be working or crashed)
- Active work branches

### Step 4: Check for PRs
```bash
gh search prs --owner zarlcorp --state open
gh pr list --state all
```

Report which items have PRs ready for review.

### Step 5: Summary

| ID | Title | Repo | Status | PR | Blockers |
|----|-------|------|--------|----|----------|
| 001 | ... | zarlcorp/tsk | in progress | — | — |
| 002 | ... | zarlcorp/tsk | blocked | — | needs X |
| 003 | ... | zarlcorp/other | PR ready | #5 | — |

### Step 6: Suggestions
- If blockers exist, suggest resolution actions
- If PRs are ready, suggest `/review <id>`
- If all items are done, suggest next steps (merge, integration testing, etc.)
