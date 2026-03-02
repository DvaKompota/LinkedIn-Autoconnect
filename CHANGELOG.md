# Changelog

All notable changes to this project will be documented in this file.

Format based on [Keep a Changelog](https://keepachangelog.com/).

## [Unreleased]

## [0.1.0] - 2026-03-01

### Added
- **Dry-run mode** (`--dry-run` flag) — validate entire workflow without side effects
  - Config overrides for safety: 2 companies, 2 pages/company, level 2, headed mode
  - Logs `[DRY-RUN] Would send invite to...` / `[DRY-RUN] Would withdraw invitation for...`
  - Invite: limits to 2 pages per company; Withdraw: limits to 2 withdrawals
- **CLI flags** (`--config`, `--feature`, `--dry-run`) — runtime configuration without editing code
- **Invite from Search** feature — send connection requests to people at target companies
  - Multi-company targeting from `search_list`
  - Smart filtering by job titles (partial match) and blacklist (exact match)
  - Pagination support across search result pages
  - Anti-bot random delays (1-2s between actions)
- **Withdraw Old Invitations** feature — revoke pending requests older than a month
  - Lazy loading support (clicks "Load more" until all invitations visible)
  - Auto-blacklists withdrawn names to prevent re-inviting
- **Browser state persistence** — login once manually, run headless forever
- **YAML configuration** — `config.yaml.example` template with gitignored personal configs

### Fixed
- LinkedIn "Sent" tab locator (changed from link to button role)
- Invitation count threshold (from 10 to 0)
