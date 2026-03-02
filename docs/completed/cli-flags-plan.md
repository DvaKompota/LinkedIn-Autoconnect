# ✅ COMPLETED - Implementation Plan: CLI Flags for Config and Feature Selection

**Status:** Implemented (pending commit)
**Implementation Date:** January 2026
**Files Modified:** `cmd/autoconnect.go`

This plan has been fully implemented. The CLI flags feature is working as designed.

---

## Overview
Add command-line flag support to `cmd/autoconnect.go` to allow dynamic config file selection and feature selection without editing source code.

## Current State
- Config path hardcoded to `data/config.yaml`
- Feature selection done by commenting/uncommenting lines in main()
- No CLI argument parsing exists
- Two available features: `InviteFromSearch` and `WithdrawOldInvitations`

## Proposed Changes

### 1. Add CLI Flags
Use Go's standard `flag` package (no external dependencies).

**Flags to add:**
- `--config` (default: `"data/config.yaml"`) - Path to config YAML file
- `--feature` (default: `"invite"`) - Feature to run: `invite` or `withdraw`

**Example usage:**
```bash
# Use dev config with invite feature (default)
go run cmd/autoconnect.go --config data/config.dev.yaml

# Use prod config with withdrawal feature
go run cmd/autoconnect.go --config data/config.yaml --feature withdraw

# Use defaults (current behavior)
go run cmd/autoconnect.go
```

### 2. File to Modify

**`cmd/autoconnect.go`** - Only file requiring changes

Changes:
1. Import `flag` package
2. Add flag definitions at start of `main()`
3. Call `flag.Parse()`
4. Replace hardcoded `configPath := "data/config.yaml"` with flag value
5. Replace commented feature selection with switch statement based on feature flag

### 3. Implementation Details

**Before (lines 14-15, 52-54):**
```go
statePath := "data/browser-state"
configPath := "data/config.yaml"
...
// Call the feature
// feature.WithdrawOldInvitations(a, cfg)
feature.InviteFromSearch(a, cfg)
```

**After:**
```go
// Parse command-line flags
configPath := flag.String("config", "data/config.yaml", "Path to config YAML file")
featureName := flag.String("feature", "invite", "Feature to run: invite or withdraw")
flag.Parse()

statePath := "data/browser-state"
cfg, err := config.LoadConfig(*configPath)
...
// Execute selected feature
switch *featureName {
case "invite":
    feature.InviteFromSearch(a, cfg)
case "withdraw":
    feature.WithdrawOldInvitations(a, cfg)
default:
    log.Fatalf("Unknown feature: %s. Valid options: invite, withdraw", *featureName)
}
```

### 4. Backward Compatibility

✅ **100% backward compatible**
- Default values match current hardcoded paths
- Running `go run cmd/autoconnect.go` behaves identically to current version
- Default feature is `invite` (current active feature)

### 5. Error Handling

Follow existing pattern:
- Invalid feature name → `log.Fatalf` with helpful message
- Invalid config path → `config.LoadConfig()` already returns error, handled by existing code

### 6. Code Style

- Use existing `log.Fatalf` error handling pattern
- Maintain current dependency injection (App, Config)
- Keep simple, procedural style
- No additional abstractions needed

## Files Changed

1. **`cmd/autoconnect.go`**
   - Add `import "flag"`
   - Add 2 flag definitions (6 lines)
   - Replace hardcoded configPath with flag
   - Replace feature comments with switch statement (~10 lines)
   - Total: ~15 lines changed/added

## Verification Steps

1. **Test default behavior (backward compatibility):**
   ```bash
   go run cmd/autoconnect.go
   # Should use data/config.yaml and run invite feature
   ```

2. **Test dev config:**
   ```bash
   go run cmd/autoconnect.go --config data/config.dev.yaml
   # Should load dev config with SpaceX only
   ```

3. **Test feature selection:**
   ```bash
   go run cmd/autoconnect.go --feature withdraw
   # Should run withdrawal feature instead of invite
   ```

4. **Test invalid inputs:**
   ```bash
   go run cmd/autoconnect.go --feature invalid
   # Should exit with error message about valid options
   ```

5. **Test help flag:**
   ```bash
   go run cmd/autoconnect.go --help
   # Should display flag usage information
   ```

## Out of Scope (Future Enhancements)

Not included in this implementation:
- `--headless` override flag (config value is sufficient for now)
- `--state` path override (single state file is fine)
- Subcommands (`autoconnect invite` vs flags)
- Config validation
- Verbose/debug logging flags
