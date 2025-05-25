package main

import (
	"log"
	"os"
	"strings"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/browser"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/linkedin"
)

func main() {
	statePath := "data/browser-state"
	configPath := "data/config.yaml"
	config, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}
	headless := config.Headless

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

	// Navigate to the invitations page and wait for invitations to load
	if err := a.Invitations.Navigate(); err != nil {
		log.Fatalf("could not navigate to invitations: %v", err)
	}
	if err := a.Invitations.WaitForLoad(); err != nil {
		log.Fatalf("could not wait for invitations to load: %v", err)
	}
	if err := a.Invitations.OpenSentTab(); err != nil {
		log.Fatalf("could not open sent invitations tab: %v", err)
	}
	if err := a.Invitations.WaitForInvitationsCountToBeMoreThan(1); err != nil {
		log.Fatalf("There was less than 10 invitations: %v", err)
	}

	// Iterate through the invitations and withdraw those older than a month, starting from the bottom
	count, _ := a.Invitations.CountInvitations()
	log.Printf("Found %d invitations", count)
	if err := a.Invitations.ScrollToBottom(); err != nil {
		log.Fatalf("could not scroll to bottom: %v", err)
	}
	for i := count - 1; i >= 0; i-- {
		invitation := a.Invitations.GetInvitationLocatorByIndex(i)
		if err := invitation.ScrollIntoViewIfNeeded(); err != nil {
			log.Printf("could not scroll invitation %d into view: %v", i, err)
			continue
		}
		name, err := a.Invitations.GetInvitationName(invitation)
		if err != nil {
			log.Printf("could not get name for invitation %d: %v", i, err)
			continue
		}
		time, err := a.Invitations.GetInvitationTime(invitation)
		if err != nil {
			log.Printf("could not get time for invitation %d: %v", i, err)
			continue
		}
		if strings.Contains(time, "month") {
			a.Sleep(1) // Slows down the process to avoid being detected as a bot
			if err := a.Invitations.WithdrawInvitation(invitation); err != nil {
				log.Printf("could not withdraw invitation %d: %v", i, err)
			}
			log.Printf("Withdrawn invitation for %s (%s)", name, time)
			// Adding person to the blacklist to prevent sending automatic invites to them in the future
			if err := appendToConfigList(configPath, "blacklist", name); err != nil {
				log.Printf("could not add %s to the blacklist: %v", name, err)
			}
		}
	}
}
