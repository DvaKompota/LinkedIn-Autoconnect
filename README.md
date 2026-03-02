# LinkedIn Autoconnect

[![Go Version](https://img.shields.io/badge/Go-1.24.2-00ADD8?logo=go)](https://go.dev/)
[![Playwright](https://img.shields.io/badge/Playwright-Go-45ba4b?logo=playwright)](https://playwright.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Automate LinkedIn networking by sending connection requests to people at target companies. Built with Go and Playwright for reliable, production-grade automation.

## Features

- **CLI-Based**: Command-line flags for config and feature selection - no code editing required
- **Browser State Persistence**: Login once, run headless forever
- **Multi-Company Targeting**: Process multiple companies in a single run
- **Smart Filtering**: Filter by job titles, blacklist specific names
- **Pagination Support**: Automatically navigate through search results
- **Invite Withdrawal**: Clean up old pending invitations
- **Anti-Bot Protection**: Random delays to avoid detection
- **Two Implementations**: Modern Go (current) + Legacy Python (Selenium)

## Quick Start

### Prerequisites

- Go 1.24+ installed
- LinkedIn account

### Installation

```bash
# Clone the repository
git clone https://github.com/DvaKompota/LinkedIn-Autoconnect.git
cd LinkedIn-Autoconnect

# Install Go dependencies
go mod download

# Install Playwright browsers (first time only)
go run github.com/playwright-community/playwright-go/cmd/playwright@latest install
```

### First Run - Manual Login

On initial execution, you'll need to log in manually to save browser state:

```bash
go run cmd/autoconnect.go
```

The app will:
1. Open a browser window
2. Wait for you to log in to LinkedIn
3. Save your session to `data/browser-state`
4. Run the default feature (invite)

### Subsequent Runs

After the first login, the app runs headlessly using saved credentials:

```bash
# Use default config and feature
go run cmd/autoconnect.go

# Use custom config
go run cmd/autoconnect.go --config data/config.dev.yaml

# Run withdrawal feature
go run cmd/autoconnect.go --feature withdraw

# Combine flags
go run cmd/autoconnect.go --config data/config.dev.yaml --feature invite
```

### Build Binary (Optional)

```bash
go build -o bin/autoconnect cmd/autoconnect.go
./bin/autoconnect --feature invite
```

## Configuration

The app uses YAML config files. **Personal configs are gitignored** - create from template:

```bash
# Copy example config and customize
cp data/config.yaml.example data/config.yaml
```

Edit `data/config.yaml` with your target companies, job titles, and preferences.

### Config Structure

```yaml
# data/config.yaml

headless: true              # Run browser in headless mode
search_level: 2             # LinkedIn search depth (1=1st, 2=2nd, 3=all connections)
connection_level: 2         # Target connection level
per_company_limit: 10       # Max invites per company per run

search_list:                # Companies to target
    - Tesla
    - SpaceX
    - Neuralink
    - Boring Company

job_titles:                 # Filter by job titles (partial match, case-insensitive)
    - Recruiter
    - Talent
    - Software

blacklist:                  # Names to skip (exact match)
    - Adolf Hitler
```

**Note:** `search_level` and `connection_level` may be redundant. See [ROADMAP.md](docs/ROADMAP.md) for ongoing evaluation.

### Config Fields

| Field | Type | Description |
|-------|------|-------------|
| `headless` | bool | Run browser without UI (true for automation, false for debugging) |
| `search_level` | int | LinkedIn search scope: 1 (1st circle), 2 (2nd circle), 3 (all) |
| `connection_level` | int | Target connection level |
| `per_company_limit` | int | Maximum invitations per company per run |
| `search_list` | []string | Companies to search for |
| `job_titles` | []string | Job title filters (partial, case-insensitive) |
| `blacklist` | []string | Names to exclude from invitations |

## CLI Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `data/config.yaml` | Path to YAML config file |
| `--feature` | `invite` | Feature to run: `invite` or `withdraw` |

### Examples

```bash
# Run with defaults
go run cmd/autoconnect.go

# Use dev config
go run cmd/autoconnect.go --config data/config.dev.yaml

# Withdraw old invitations
go run cmd/autoconnect.go --feature withdraw

# Help
go run cmd/autoconnect.go --help
```

## Features

### 1. Invite from Search

Sends connection requests to people at target companies based on your config.

**Workflow:**
1. Navigates to LinkedIn search with configured filters (1st/2nd/all connections)
2. Iterates through companies in `search_list`
3. Filters by company using LinkedIn's search filters
4. Scrolls to load all contact cards on current page
5. Validates each contact card:
   - Must have "Connect" button available
   - Name must not be in blacklist
   - Job title must match at least one term from `job_titles`
6. Sends invites up to `per_company_limit` per company
7. Moves to next page if available
8. Repeats for all companies

**Run:**
```bash
go run cmd/autoconnect.go --feature invite
```

### 2. Withdraw Old Invitations

Revokes pending connection requests older than a specified time period.

**Run:**
```bash
go run cmd/autoconnect.go --feature withdraw
```

## Architecture

The Go implementation uses a clean, layered architecture:

```
cmd/autoconnect.go              # Entry point, CLI flags, orchestration
internal/
  ├── browser/                  # Playwright browser wrapper + state persistence
  ├── config/                   # YAML config loading/management
  ├── feature/                  # High-level automation workflows
  │   ├── invite_from_search.go
  │   └── withdraw.go
  └── linkedin/                 # Page objects for LinkedIn pages
      ├── app.go                # Main App struct, page registry
      ├── login.go
      ├── search.go
      └── invitations.go
```

**Design Patterns:**
- **Page Object Model**: LinkedIn pages encapsulated as Go structs with locators and methods
- **Browser State Persistence**: Playwright's state API replaces credential management
- **Feature-Based**: High-level features compose page objects
- **Dependency Injection**: `App` struct provides access to all page objects

See [CLAUDE.md](CLAUDE.md) for detailed development guidance.

## Legacy Python Implementation

The repository includes a legacy Python implementation using Selenium WebDriver:

### Setup (Python)

```bash
# Install dependencies
./setup.sh
# or manually:
pip3 install -r requirements.txt

# Create credentials file
echo 'email = "your.email@gmail.com"' > data/credentials.py
echo 'password = "your_password"' >> data/credentials.py
```

### Running (Python)

```bash
# Withdraw old invitations
./modules/withdraw_old_invites.py

# Send invites from search
./modules/subscribe_from_search.py

# Send invites from 1st circle profiles
./modules/subscribe_from_profiles.py
```

**Note:** The Python implementation is deprecated. New development focuses on the Go version.

## Project Documentation

- **[CLAUDE.md](CLAUDE.md)** - Comprehensive guide for AI-assisted development
- **[ROADMAP.md](docs/ROADMAP.md)** - Project roadmap and planned features
- **[docs/completed/](docs/completed/)** - Completed feature implementation plans
  - [CLI Flags Implementation Plan](docs/completed/cli-flags-plan.md)

## Roadmap

**Current Version:** 0.1.0 (Personal Tool)
- ✅ 2 of 3 core features complete
- ⚠️ Missing: Invite from profiles (search limit workaround)

See **[Full Roadmap](docs/ROADMAP.md)** for detailed planning and implementation notes.

### Critical Priority: Feature Parity with Python

**Goal:** Complete Go implementation, then sunset Python version.

🔴 **1. Implement "Invite from Profiles" Feature**
- **Why:** LinkedIn limits search results for 2nd/3rd-degree connections
- **What:** When search limit hit, browse 1st-degree profiles and invite from "People also viewed" section
- **Status:** Not started (Python reference: `modules/subscribe_from_profiles.py`)

🔴 **2. Add Dry-Run Mode**
- **Why:** LinkedIn UI changes frequently, breaking locators
- **What:** `--dry-run` flag to test workflow without sending invites
- **Use case:** Catch UI breakage before production runs

### Next Release: v0.2.0 - Feature Parity

- ✅ Implement "Invite from Profiles"
- ✅ Add dry-run mode for all features
- ✅ Remove Python implementation from repo

### Reliability Improvements (Secondary Priority)

- Fix silent error swallowing (4+ locations)
- Add configuration validation
- Complete modal handling (rate limit warnings)
- Better logging (verbose/quiet modes)
- Progress tracking & resume capability

### Understanding LinkedIn Search Limits

LinkedIn limits search for 2nd/3rd-degree connections. The app has two strategies:

1. **`invite`** (Invite from Search) - Direct search until limit hit
2. **`profiles`** (Invite from Profiles) - Browse 1st-circle profiles (unlimited), use "People also viewed"
3. **`withdraw`** (Withdraw Old Invitations) - Cleanup pending invites

**Recommended workflow:** Run `invite` → hit limit → switch to `profiles` → periodically `withdraw`

## Development

### Adding New Features

1. Create feature in `internal/feature/`
2. Add switch case in `cmd/autoconnect.go`
3. Update CLAUDE.md and README.md
4. Test thoroughly in headed mode first

### Adding New Page Objects

1. Create new file in `internal/linkedin/`
2. Define page struct with Playwright `Page` and locators
3. Add constructor `NewXxxPage(page playwright.Page)`
4. Register in `App` struct in `app.go`

Example:
```go
type ProfilePage struct {
    page   playwright.Page
    header playwright.Locator
}

func NewProfilePage(page playwright.Page) *ProfilePage {
    return &ProfilePage{
        page:   page,
        header: page.Locator(".profile-header"),
    }
}
```

### Locator Updates

LinkedIn frequently updates their UI. When locators break, update:
- **Go**: Locator strings in `internal/linkedin/*.go`
- **Python**: Locator strings in `pages/*_page.py`

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Disclaimer

This tool is for educational purposes. Use responsibly and in compliance with LinkedIn's Terms of Service. Excessive automation may violate LinkedIn's policies.

## Author

**Sergey Kolokolov** - [GitHub](https://github.com/DvaKompota)

---

**Note:** This is a personal project for learning and automation practice. Not affiliated with LinkedIn.
