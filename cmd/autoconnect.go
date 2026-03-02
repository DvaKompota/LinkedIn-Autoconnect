package main

import (
	"flag"
	"log"
	"os"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/browser"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/config"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/feature"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/linkedin"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "data/config.yaml", "Path to config YAML file")
	featureName := flag.String("feature", "invite", "Feature to run: invite or withdraw")
	dryRun := flag.Bool("dry-run", false, "Test mode: validate workflow without sending invites or withdrawing")
	flag.Parse()

	statePath := "data/browser-state"
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Override config settings in dry-run mode for safety and speed
	if *dryRun {
		log.Println("[DRY-RUN MODE] Overriding config for testing:")
		log.Println("  - search_list: first 2 companies only")
		log.Println("  - page limit: 2 pages per company")
		log.Println("  - search_level: 2, connection_level: 2")
		log.Println("  - headless: false (visual confirmation)")
		if len(cfg.SearchList) > 2 {
			cfg.SearchList = cfg.SearchList[:2]
		}
		cfg.SearchLevel = 2
		cfg.ConnectionLevel = 2
		cfg.Headless = false
	}

	headless := cfg.Headless

	// Check if state file exists for persistent login
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// State file doesn't exist, perform initial login
		b, err := browser.NewBrowser(false, "") // Headed mode, no state
		if err != nil {
			log.Fatalf("could not create browser: %v", err)
		}
		defer b.Close()

		page := b.Page
		loginPage := linkedin.NewLoginPage(page)
		if err := loginPage.Navigate(); err != nil {
			log.Fatalf("could not navigate to login: %v", err)
		}
		if err := loginPage.WaitForLogin(60000); err != nil { // Wait 60 seconds for manual login
			log.Fatalf("login failed: %v", err)
		}
		if err := b.SaveState(statePath); err != nil {
			log.Fatalf("could not save state: %v", err)
		}
		b.Close()
	}

	// Create the app with the saved state
	a, err := linkedin.NewApp(headless, statePath, *dryRun)
	if err != nil {
		log.Fatalf("could not create app: %v", err)
	}
	defer a.Close()

	// Execute selected feature
	switch *featureName {
	case "invite":
		feature.InviteFromSearch(a, cfg)
	case "withdraw":
		feature.WithdrawOldInvitations(a, cfg)
	default:
		log.Fatalf("Unknown feature: %s. Valid options: invite, withdraw", *featureName)
	}
}
