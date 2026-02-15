# 041: Homebrew tap setup with formula template

## Objective
Create the `zarlcorp/homebrew-tap` repository with a formula template and automation so that `brew install zarlcorp/tap/<tool>` works for every zarlcorp tool.

## Context
The manifesto defines the distribution goal: "brew install zarlcorp/tap/zburn" as the validation of the full pipeline. Homebrew taps are the standard distribution channel for Go CLI tools on macOS and Linux.

GoReleaser can auto-update Homebrew formulas on release. This spec sets up the tap repo and the formula structure so GoReleaser has a target.

No dependencies on other specs — can be done any time. But it's most useful after tool repos exist (specs 037-039) and have a release workflow (spec 040).
Issue #27.

## Requirements

### Repository: zarlcorp/homebrew-tap

### Formula template
Each tool gets a formula file. Start with a placeholder for zburn:

```ruby
# Formula/zburn.rb
class Zburn < Formula
  desc "Disposable identity generator — burner emails, names, passwords"
  homepage "https://github.com/zarlcorp/zburn"
  license "MIT"

  # GoReleaser populates these fields automatically on release
  # version, url, sha256 are filled by goreleaser

  def install
    bin.install "zburn"
  end

  test do
    assert_match "zburn", shell_output("#{bin}/zburn version")
  end
end
```

### GoReleaser integration
The formula is updated automatically by GoReleaser's `brews` section in each tool repo's `.goreleaser.yml`:

```yaml
brews:
  - repository:
      owner: zarlcorp
      name: homebrew-tap
    directory: Formula
    homepage: "https://github.com/zarlcorp/zburn"
    description: "Disposable identity generator"
    license: "MIT"
```

This spec only creates the tap repo and placeholder formula. The GoReleaser config lives in each tool repo.

### Repository structure
```
zarlcorp/homebrew-tap/
├── Formula/
│   └── zburn.rb             # placeholder — GoReleaser overwrites on release
├── README.md                # usage: brew install zarlcorp/tap/zburn
└── LICENSE                  # MIT
```

### Automation token
GoReleaser needs push access to the tap repo. Set up a `HOMEBREW_TAP_TOKEN` secret (or use the org-level PAT) in tool repos that GoReleaser can use to push formula updates.

## Target Repo
zarlcorp/homebrew-tap

## Agent Role
devops

## Files to Modify
All files are new — this is a new repository.
- `Formula/zburn.rb` — placeholder formula
- `README.md` — usage instructions
- `LICENSE` — MIT

## Notes
- Homebrew tap repos must be named `homebrew-tap` (or `homebrew-<name>`) for `brew tap zarlcorp/tap` to work.
- The tap repo is mostly automated — GoReleaser pushes formula updates on every release. Manual edits are rare.
- zvault and zshield formulas are added when those tools ship their first release. Don't create placeholder formulas for tools that don't exist yet.
