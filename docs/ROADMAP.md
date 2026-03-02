# LinkedIn-Autoconnect Roadmap

This roadmap outlines planned improvements and features for the Go implementation of LinkedIn-Autoconnect, a personal automation tool for expanding LinkedIn networks.

## Project Status

**Current Version:** 0.1.0 (Personal Tool)
- ✅ CLI flags implementation
- ✅ Browser state persistence
- ✅ Invite from search (2/3 features complete)
- ✅ Invite withdrawal
- ⚠️ **Missing:** Invite from profiles (search limit workaround)

**Project Scope:** Personal automation tool for individual use, not enterprise software.

**Primary Goal:** Achieve feature parity with Python implementation, then sunset Python version.

---

## Critical Priority: Feature Parity with Python

### 🔴 1. Implement "Invite from Profiles" Feature
**Priority:** CRITICAL | **Effort:** Medium | **Status:** Not Started

**Problem:** Go implementation is missing the third core feature that exists in Python.

**What it does:**
- When LinkedIn search limit is exhausted, this is the workaround
- Searches for 1st-degree connections at target companies
- Opens each connection's profile
- Expands "People also viewed" section (shows 2nd/3rd-degree connections)
- Sends invites to people from target companies in that section
- Continues until per-company limit reached

**Why it's needed:**
- LinkedIn limits search results for 2nd/3rd-degree connections
- Searching within 1st circle (browsing your contacts) is unlimited
- This feature enables continued invitations when search is exhausted

**Implementation Requirements:**
- [ ] Create `ProfilePage` in `internal/linkedin/profile.go`
- [ ] Locators for "People also viewed" section
- [ ] "Show more" button handling
- [ ] Extract contact cards from "People also viewed"
- [ ] Reuse existing filtering logic (job titles, blacklist)
- [ ] Create `InviteFromProfiles` feature in `internal/feature/`
- [ ] Add `profiles` option to `--feature` flag
- [ ] Update CLAUDE.md and README.md

**Python Reference:**
- `modules/subscribe_from_profiles.py` - Main workflow
- `pages/profile_page.py` - Page object

**Success Criteria:**
- Feature works identically to Python version
- Can sunset Python implementation and remove from repo

---

### 🔴 2. Dry-Run Mode for All Features
**Priority:** CRITICAL | **Effort:** Small | **Status:** Not Started

**Problem:** LinkedIn frequently changes UI, breaking locators. No way to test without sending real invites.

**What it does:**
- Run entire workflow without actually sending invites or withdrawing
- Validates all locators still work
- Logs what *would* happen (e.g., "Would send invite to John Doe")
- Catches UI breakage before production runs

**Implementation:**
- [ ] Add `--dry-run` flag (boolean)
- [ ] Pass dry-run flag through App and features
- [ ] Replace `Connect()` with log statement in dry-run mode
- [ ] Replace `Withdraw()` with log statement in dry-run mode
- [ ] All page navigation and locators execute normally
- [ ] Add dry-run indicator to logs
- [ ] Override config settings in dry-run mode for safety and speed:
  - `per_company_limit: 1-2` (speed up testing)
  - `search_level: 3` and `connection_level: 3` (works for everyone, no dependency on real connections)
  - `headless: false` (visual confirmation, up for debate)

**Usage:**
```bash
# Test invite feature without sending
go run cmd/autoconnect.go --feature invite --dry-run

# Test profiles feature without sending
go run cmd/autoconnect.go --feature profiles --dry-run

# Test withdrawal without withdrawing
go run cmd/autoconnect.go --feature withdraw --dry-run
```

**Success Criteria:**
- Can validate entire workflow without side effects
- Catches locator breakage immediately
- Safe to run frequently to detect LinkedIn UI changes
- Config overrides make it fast and universally runnable

---

## High Priority: Reliability Improvements

### 🟡 3. Fix Silent Error Swallowing
**Priority:** HIGH | **Effort:** Small | **Status:** Not Started

**Problem:** Errors discarded with `_` cause silent failures.

**Locations to fix:**
```go
// internal/feature/invite_from_search.go
count, _ := a.Search.CountPeople()  // Line 57
if bool, _ := a.Search.IsContactCardValid(i); !bool { }  // Line 67
if bool, _ := a.Search.IsNextButtonEnabled(); bool { }  // Line 104

// internal/feature/withdraw.go
count, _ := a.Invitations.CountInvitations()  // Line 27
```

**Fix:**
- [ ] Replace `_` with proper error handling
- [ ] Log errors and continue (don't crash on non-critical errors)
- [ ] Return errors for critical failures
- [ ] Add retry logic for transient errors (timeouts)

**Impact:** Prevents silent failures, improves debuggability

---

### 🟡 4. Configuration Validation
**Priority:** HIGH | **Effort:** Small | **Status:** Not Started

**Problem:** Invalid configs accepted without validation.

**Tasks:**
- [ ] Validate `search_level` is 1-3
- [ ] Validate `per_company_limit > 0`
- [ ] Require `search_list` non-empty
- [ ] Require `job_titles` non-empty
- [ ] Add validation tests
- [ ] Return clear error messages

**Impact:** Fails fast with clear errors instead of silent no-ops

---

### 🟡 4.5. Evaluate `search_level` vs `connection_level` Configuration
**Priority:** MEDIUM | **Effort:** Small | **Status:** Analysis Needed

**Problem:** Two similar config fields that may be redundant or poorly named.

**Current State:**
- `search_level`: LinkedIn search depth (1=1st circle, 2=2nd circle, 3=all)
- `connection_level`: Target connection level

**Questions:**
- Are both necessary or is this redundant?
- Should they serve different purposes for different features?
  - `search_level` for `invite` (search-based invites)
  - `connection_level` for `profiles` (profile-based invites)
- Should naming be more explicit if they serve different purposes?

**Potential Use Case:**
- Search and invite only 2nd-degree from search (`search_level: 2`)
- But invite everyone (2nd/3rd) from profiles (`connection_level: 3`)
- Or vice versa

**Tasks:**
- [ ] Analyze Python implementation usage
- [ ] Determine if both are needed
- [ ] Rename for clarity if they serve different purposes
- [ ] Consolidate if redundant
- [ ] Update docs with final decision

**Impact:** Clearer configuration, less confusion

---

### 🟡 5. Complete Modal Handling
**Priority:** MEDIUM-HIGH | **Effort:** Medium | **Status:** Partially Implemented

**Problem:** App crashes on LinkedIn rate limit or congratulatory modals.

**Missing Modals:**
```go
// TODO in internal/linkedin/search.go:
// "Nice job building your network!"
// "You've sent more invitations than most"
// "You're close to the weekly invitation limit"  ⚠️
// "You've reached the weekly invitation limit"   ⚠️ Crashes here
```

**Tasks:**
- [ ] Detect and handle all modal types
- [ ] Log when rate limit is approached
- [ ] Stop gracefully at rate limit (don't crash)
- [ ] Optional: Track weekly invite count locally

**Impact:** Prevents crashes, enables safe long-running automation

---

## Medium Priority: Quality of Life

### 🟢 6. Better Logging
**Priority:** MEDIUM | **Effort:** Small | **Status:** Not Started

**Current:** Plaintext printf logs

**Improvements:**
- [ ] Add `--verbose` flag for detailed logging
- [ ] Add `--quiet` flag for minimal output
- [ ] Structured log format (timestamp, level, message)
- [ ] Summary report at end (total invites sent, skipped, errors)

**Impact:** Better visibility into automation progress

---

### 🟢 7. Progress Tracking
**Priority:** MEDIUM | **Effort:** Medium | **Status:** Not Started

**Problem:** If app crashes, no way to resume. Re-running may re-invite people.

**Tasks:**
- [ ] Log processed companies to file (`data/progress.json`)
- [ ] Add `--resume` flag to skip already-processed companies
- [ ] Track sent invites to prevent duplicates
- [ ] Clear progress on successful completion

**Impact:** Safe resumption after crashes

---

### 🟢 8. Locator Resilience
**Priority:** MEDIUM | **Effort:** Medium-Large | **Status:** Not Started

**Problem:** Brittle locators (nth-child, CSS utility classes) break frequently.

**Tasks:**
- [ ] Replace nth-child selectors with semantic locators
- [ ] Add fallback locators (try primary, then fallback)
- [ ] Extract locators to separate file for easy updates
- [ ] Document locator update process

**Impact:** Reduces breakage on LinkedIn UI changes

---

## Low Priority: Nice-to-Have

### 🔵 9. Enhanced Filtering
**Priority:** LOW | **Effort:** Small-Medium | **Status:** Backlog

**Possible Enhancements:**
- [ ] Regex support for job titles
- [ ] Location filtering (if LinkedIn exposes it)
- [ ] Fuzzy name matching for blacklist
- [ ] Whitelist mode (only invite specific people)

---

### 🔵 10. Invite Customization
**Priority:** LOW | **Effort:** Medium | **Status:** Backlog

**Enhancements:**
- [ ] Add optional custom invite message
- [ ] Template system with variables (name, company)
- [ ] Per-company custom messages

---

### 🔵 11. Multi-Account Support
**Priority:** LOW | **Effort:** Medium | **Status:** Backlog

**Tasks:**
- [ ] Support multiple browser state files
- [ ] Profile-based config (`--profile account1`)
- [ ] Account switching

---

## Version Milestones

### v0.2.0 - Feature Parity (Target: Next Release)
**Focus:** Complete Python feature parity, sunset Python

- ✅ Implement "Invite from Profiles" feature
- ✅ Add dry-run mode
- ✅ Remove Python implementation from repo
- ✅ Update documentation

**Exit Criteria:** All 3 features work in Go, Python deleted

---

### v0.3.0 - Reliability (Target: Future)
**Focus:** Improve reliability and error handling

- ✅ Fix silent error swallowing
- ✅ Config validation
- ✅ Complete modal handling
- ✅ Better logging

**Exit Criteria:** Can run reliably without crashes

---

### v0.4.0 - Quality of Life (Target: Future)
**Focus:** Polish and UX improvements

- Progress tracking & resume
- Locator resilience
- Enhanced filtering

**Exit Criteria:** Pleasant to use, easy to maintain

---

## LinkedIn Search Limits - Important Context

**Why Multiple Features Exist:**

LinkedIn limits search for 2nd and 3rd-degree connections. When you hit the limit:

1. **1st Strategy:** `invite` (Invite from Search)
   - Search directly for 2nd/3rd-degree connections at target companies
   - Works until LinkedIn search limit is hit

2. **2nd Strategy:** `profiles` (Invite from Profiles)
   - Search within 1st-degree connections (unlimited)
   - Open each 1st-degree profile
   - Use "People also viewed" section to find 2nd/3rd-degree connections
   - Workaround when search is exhausted

3. **Cleanup:** `withdraw` (Withdraw Old Invitations)
   - Revoke old pending invites
   - Keeps inbox clean

**Recommended Workflow:**
1. Run `invite` until search limit hit
2. Switch to `profiles` to continue
3. Periodically run `withdraw` to clean up

---

## Contributing to the Roadmap

This is a personal project, but ideas are welcome! Open an issue or PR.

**Prioritization:**
- 🔴 **CRITICAL** - Blocking feature parity or data loss
- 🟡 **HIGH** - Reliability, crash prevention
- 🟢 **MEDIUM** - Quality of life, UX
- 🔵 **LOW/BACKLOG** - Nice-to-have, future ideas

---

**Last Updated:** 2026-03-01
