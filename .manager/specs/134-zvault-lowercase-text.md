# 134: Lowercase all TUI text in zvault

## Objective
Convert all user-facing title-case text in zvault's TUI to lowercase, matching the terse lowercase style used across zarlcorp tools.

## Context
The zvault TUI currently displays text like "Create New Vault", "Unlock Vault", "Secrets", "Password", etc. in title case. The zarlcorp style is lowercase everywhere. This spec fixes all title-case strings in the TUI.

## Requirements

### 1. Password screen (password.go)
- "Create New Vault" → "create new vault"
- "Unlock Vault" → "unlock vault"
- "Choose a master password to protect your vault." → "choose a master password to protect your vault."
- "Enter your master password." → "enter your master password."
- "Password" (field label) → "password"
- "Confirm" (field label) → "confirm"

### 2. View titles (navigation.go)
- "Unlock" → "unlock"
- "Menu" → "menu"
- "Secrets" → "secrets"
- "Secret" → "secret"
- "Edit Secret" → "edit secret"
- "Tasks" → "tasks"
- "Task" → "task"
- "Edit Task" → "edit task"

### 3. Menu items (menu.go)
- "Secrets" → "secrets"
- "Tasks" → "tasks"

### 4. Secret detail labels (secret_detail.go)
- "Name" → "name"
- "Type" → "type"
- "URL" → "url"
- "Username" → "username"
- "Password" → "password"
- "TOTP Secret" → "totp secret"
- "Notes" → "notes"
- "Service" → "service"
- "Key" → "key"
- "Label" → "label"
- "Private Key" → "private key"
- "Public Key" → "public key"
- "Passphrase" → "passphrase"
- "Content" → "content"
- "Tags" → "tags"
- "Created" → "created"
- "Updated" → "updated"

### 5. Secret form (secret_form.go)
- "Type" label → "type"
- Type options: "Password" → "password", "API Key" → "api key", "SSH Key" → "ssh key", "Note" → "note"
- Field labels in addInput calls: "Name" → "name", "Username" → "username", "Password" → "password", "Notes" → "notes", "Service" → "service", "Key" → "key", "Label" → "label", "Passphrase" → "passphrase", "Content" → "content", "Tags" → "tags"

### 6. Task detail (task_detail.go)
- "Status" → "status"
- "Priority" → "priority"
- "Due" → "due"
- "Tags" → "tags"
- "Created" → "created"
- "Completed" → "completed"
- Status values: "Pending" → "pending", "Done" → "done"
- Priority values: "High" → "high", "Medium" → "medium", "Low" → "low", "None" → "none"

### 7. Task form (task_form.go)
- "Title" → "title"
- "Priority" → "priority"
- "Due Date" → "due date"
- "Tags" → "tags"

### 8. Secret list filter labels (secret_list.go)
- "All" → "all"
- "Password" → "password"
- "API Key" → "api key"
- "SSH Key" → "ssh key"
- "Note" → "note"

### 9. Update tests
All test assertions that check for title-case strings must be updated to match the new lowercase values.

### What NOT to change
- Go identifiers (ZvaultAccent, etc.) — these follow Go naming conventions
- ZVAULT_PASSWORD environment variable — follows Unix convention
- Import paths
- Error messages (already lowercase)

## Target Repo
zarlcorp/zvault

## Agent Role
backend

## Files to Modify
- internal/tui/password.go
- internal/tui/navigation.go
- internal/tui/menu.go
- internal/tui/secret_detail.go
- internal/tui/secret_form.go
- internal/tui/secret_list.go
- internal/tui/task_detail.go
- internal/tui/task_form.go
- internal/tui/tui_test.go
- internal/tui/secret_detail_test.go
- internal/tui/secret_form_test.go
- internal/tui/secret_list_test.go
- internal/tui/task_detail_test.go
- internal/tui/task_form_test.go
- internal/tui/task_list_test.go

## Notes
This is a mechanical find-and-replace across the TUI layer. The key risk is missing test assertions that check for the old casing — grep for all title-case strings in test files and update them.
