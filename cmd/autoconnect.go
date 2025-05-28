package main

import (
	"log"
	"os"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/browser"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/config"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/feature"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/linkedin"
)

func main() {
	statePath := "data/browser-state"
	configPath := "data/config.yaml"
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
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
	a, err := linkedin.NewApp(headless, statePath)
	if err != nil {
		log.Fatalf("could not create app: %v", err)
	}
	defer a.Close()

    // Call the feature
    feature.WithdrawOldInvitations(a, cfg)
}
