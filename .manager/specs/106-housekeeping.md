# 106: Housekeeping — prune stale branches and close issues

## Objective
Clean up stale work branches from completed specs across all repos, and close resolved issues.

## Context
After completing the zburn feature expansion (095-105), there are 7 stale work branches across 5 repos from earlier completed work. These should be pruned to keep the repos clean.

## Requirements

### Prune stale work branches
- core: `work/094-core-github-pages`, `work/095-remove-go-gets`, `work/096-docs-syntax-highlight`
- zburn: `work/097-readme-links`
- zvault: `work/098-readme-links`
- zshield: `work/099-readme-links`
- dot-github: `work/100-readme-links`

### Close resolved issues
- zburn #25 (068: tool page nav consistency) — check if the work was completed in earlier PRs, close if resolved

## Target Repo
zarlcorp/core (tracking)

## Agent Role
PM task — not delegated to agents

## Notes
This is a PM-executed task, not delegated. Branch pruning and issue triage are git/gh operations that don't require code changes.
