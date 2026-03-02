# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

LinkedIn Autoconnect is a LinkedIn automation tool for expanding your professional network by automatically sending connection requests to people working at target companies. The project contains two implementations:

- **Python (Legacy)**: Uses Selenium WebDriver with the Page Object Model pattern
- **Go (Current)**: Uses Playwright with a clean architecture approach

## Project Structure

```
LinkedIn-Autoconnect/
├── cmd/                    # Entry points
│   └── autoconnect.go     # Main CLI application
├── internal/              # Private application code
│   ├── browser/           # Browser wrapper
│   ├── config/            # Config management
│   ├── feature/           # Feature implementations
│   └── linkedin/          # Page objects
├── modules/               # Python scripts (legacy)
├── pages/                 # Python page objects (legacy)
├── data/                  # Config files, browser state (gitignored)
├── docs/                  # Project documentation
│   └── completed/         # Completed implementation plans
├── README.md              # Public-facing documentation
└── CLAUDE.md             # AI development guide (this file)
```

**Documentation Structure:**
- [README.md](README.md) - User-facing documentation for GitHub
- [CLAUDE.md](CLAUDE.md) - Comprehensive AI development guide
- [docs/completed/](docs/completed/) - Archived completed feature plans
  - [cli-flags-plan.md](docs/completed/cli-flags-plan.md) - CLI flags implementation details

## Go Implementation (Current)

### Architecture

The Go codebase follows a clean, layered architecture:

- `cmd/autoconnect.go` - Main entry point and orchestration
- `internal/browser/` - Playwright browser wrapper with state persistence
- `internal/config/` - YAML configuration loading and management
- `internal/feature/` - High-level feature implementations (invite workflows, withdrawal)
- `internal/linkedin/` - Page objects and LinkedIn-specific automation (app, login, search, invitations)

### LinkedIn Search Limits - Important Context

LinkedIn limits search results for 2nd and 3rd-degree connections. This is why the app has **three distinct features**:

1. **`invite` (Invite from Search)** - `internal/feature/invite_from_search.go`
   - Searches directly for 2nd/3rd-degree connections at target companies
   - Works until LinkedIn's search limit is hit
   - Most efficient when you haven't exhausted search

2. **`profiles` (Invite from Profiles)** - ⚠️ **NOT YET IMPLEMENTED IN GO**
   - Python reference: `modules/subscribe_from_profiles.py`
   - Searches within 1st-degree connections (unlimited)
   - Opens each 1st-degree profile
   - Expands "People also viewed" section (shows 2nd/3rd-degree connections)
   - Sends invites to people from target companies in that section
   - **Workaround when search is exhausted**

3. **`withdraw` (Withdraw Old Invitations)** - `internal/feature/withdraw.go`
   - Revokes pending connection requests older than threshold
   - Keeps inbox clean

**Recommended workflow:**
1. Run `invite` until search limit hit
2. Switch to `profiles` to continue (when implemented)
3. Periodically run `withdraw` to clean up old pending invites

**Current Status (v0.1.0):**
- ✅ `invite` - Implemented
- ✅ `withdraw` - Implemented
- ⚠️ `profiles` - **Missing** (blocks feature parity with Python)

### Running the Application

```bash
# Install dependencies
go mod download

# Install Playwright browsers (first time only)
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install

# Build the binary (optional)
go build -o bin/autoconnect cmd/autoconnect.go

# Run with default config (data/config.yaml) - invite feature
go run cmd/autoconnect.go

# Run with custom config
go run cmd/autoconnect.go --config data/config.dev.yaml

# Run specific feature
go run cmd/autoconnect.go --feature withdraw  # Withdraw old invitations
go run cmd/autoconnect.go --feature invite    # Send invites (default)

# Run from built binary
./bin/autoconnect --feature invite
```

**First Run**: On initial execution, the app opens a browser window for you to manually log in to LinkedIn. After successful login, it saves browser state to `data/browser-state` for future automated runs.

**Subsequent Runs**: The app loads the saved browser state and runs headlessly (or headed, depending on config).

### CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `data/config.yaml` | Path to YAML config file |
| `--feature` | `invite` | Feature to run: `invite` or `withdraw` |

### Configuration

**Config files are gitignored** - `data/config.yaml` is in `.gitignore`. Create your own local copy:

```bash
# Copy the template from the repo (if it exists) or create from scratch
cp data/config.yaml data/config.dev.yaml  # For local development
# Edit your local config
```

Config structure (`data/config.yaml` or `data/config.dev.yaml`):

```yaml
headless: true              # Run browser in headless mode
search_level: 3             # LinkedIn search depth (1=1st circle, 2=2nd circle, 3=all)
connection_level: 2         # Target connection level
per_company_limit: 10       # Max invites per company
search_list:                # Companies to target
    - Tesla
    - SpaceX
job_titles:                 # Filter by job titles (partial match)
    - Recruiter
    - Software
blacklist:                  # Names to skip
    - John Doe
```

**Programmatic config updates**: Use `config.AppendToList()` to add items to `search_list`, `job_titles`, or `blacklist` at runtime. Changes are persisted back to the YAML file.

## Python Implementation (Legacy)

### Architecture

Uses the Page Object Model pattern:

- `modules/` - Feature scripts (executable entry points)
- `pages/` - Page objects (login_page, search_page, profile_page, my_network_page)
- `utils/` - Utility functions
- `data/config.py` - Configuration
- `data/credentials.py` - LinkedIn credentials (not in repo, create manually)

### Setup

```bash
# Install dependencies
./setup.sh
# or manually:
pip3 install -r requirements.txt

# Create credentials file
echo 'email = "your.email@gmail.com"' > data/credentials.py
echo 'password = "your_password"' >> data/credentials.py
```

### Running Scripts

```bash
# Withdraw old pending invitations
./modules/withdraw_old_invites.py

# Send invites from LinkedIn search
./modules/subscribe_from_search.py

# Send invites from 1st circle profiles (when search limit reached)
./modules/subscribe_from_profiles.py
```

## Key Differences Between Implementations

| Aspect | Python (Legacy) | Go (Current) |
|--------|----------------|--------------|
| Browser | Selenium + ChromeDriver | Playwright |
| Login | Credentials in file | Browser state persistence |
| Execution | Separate scripts | Single binary with feature toggle |
| Config | Python file | YAML file |
| Headless | ChromeOptions | Playwright headless mode |

## Important Implementation Notes

### Browser State Management (Go)

The Go implementation uses Playwright's state persistence instead of credentials. The browser state file (`data/browser-state`) contains cookies and auth tokens. Delete this file to force a fresh login.

### Page Object Pattern

Both implementations use page objects to encapsulate LinkedIn page interactions:
- **Locators**: Defined as constants/methods within page objects
- **Actions**: Methods like `Connect()`, `FilterByCompany()`, `WaitForLogin()`
- **Queries**: Methods like `CountPeople()`, `CanConnect()`, `TitleMatches()`

### Anti-Bot Measures

The Go implementation includes random sleep delays (`1 + rand.Float64()` seconds, i.e., 1-2 seconds) between connection requests to avoid detection. See [invite_from_search.go:86](internal/feature/invite_from_search.go#L86).

### Search Pagination and Scrolling

The `InviteFromSearch` feature processes search results with:
1. **Per-company iteration**: Filters by one company at a time from `cfg.SearchList`
2. **Scroll to bottom**: Loads all contact cards on the current page
3. **Per-company limit**: Stops at `cfg.PerCompanyLimit` connections per company
4. **Pagination**: Continues to next page if `Next` button is enabled
5. **Full reset**: Resets company filter between companies

See [invite_from_search.go:34-119](internal/feature/invite_from_search.go#L34-L119) for the full workflow.

### Contact Card Validation Logic

Before sending an invite, the app validates contact cards against multiple criteria (see [invite_from_search.go:81-100](internal/feature/invite_from_search.go#L81-L100)):

1. **Card validity**: Must be a valid contact card (has required elements)
2. **Connect button**: Must have an active "Connect" button
3. **Blacklist**: Name must not be in `cfg.Blacklist`
4. **Job title**: Title must contain at least one term from `cfg.JobTitles` (partial match, case-insensitive)

All four conditions must be true to send the invite. Skipped cards are logged with the reason.

## Development Workflow

When modifying features:

1. **Go**: Changes to `internal/feature/` affect automation logic
2. **Go**: Changes to `internal/linkedin/` affect page interactions and locators
3. **Go**: The `App` struct in `internal/linkedin/app.go` provides access to all page objects
4. **Python**: Changes to `modules/` affect feature scripts
5. **Python**: Changes to `pages/` affect page object implementations

### Adding New LinkedIn Page Objects (Go)

Create a new file in `internal/linkedin/` following the established pattern:

```go
package linkedin

import "github.com/playwright-community/playwright-go"

type ProfilePage struct {
    page     playwright.Page
    header   playwright.Locator
    // ... other locators
}

func NewProfilePage(page playwright.Page) *ProfilePage {
    return &ProfilePage{
        page:   page,
        header: page.Locator(".profile-header"),
        // Initialize all locators in constructor
    }
}

// Navigation, actions, and query methods
func (p *ProfilePage) Navigate(url string) error { ... }
```

Then register the page in [app.go:9-14](internal/linkedin/app.go#L9-L14):
- Add field to `App` struct
- Initialize in `NewApp()` with `NewProfilePage(page)`

### Locator Updates

LinkedIn frequently changes their UI. When locators break:
- **Go**: Update locator strings in `internal/linkedin/*.go`
- **Python**: Update locator strings in `pages/*_page.py`

Look for CSS selectors, XPath expressions, or text-based locators that may need updating.
