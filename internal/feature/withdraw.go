package feature

import (
	"log"
	"strings"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/config"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/linkedin"
)

func WithdrawOldInvitations(a *linkedin.App, cfg *config.Config) {
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
	if err := a.Invitations.WaitForInvitationsCountToBeMoreThan(0); err != nil {
		log.Fatalf("There were no invitations: %v", err)
	}
	if err := a.Invitations.LoadAllInvitations(); err != nil {
		log.Printf("Could not load all invitations: %v", err)
	}

	// Iterate through the invitations and withdraw those older than a month, starting from the bottom
	count, _ := a.Invitations.CountInvitations()
	log.Printf("Found %d invitations", count)
	if err := a.Invitations.ScrollToBottom(); err != nil {
		log.Fatalf("could not scroll to bottom: %v", err)
	}
	withdrawn := 0
	for i := count - 1; i >= 0; i-- {
		if a.DryRun && withdrawn >= 2 {
			log.Printf("[DRY-RUN] Validated 2 withdrawals, skipping remaining invitations")
			break
		}

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
			if err := a.Invitations.WithdrawInvitation(invitation, name, a.DryRun); err != nil {
				log.Printf("could not withdraw invitation %d: %v", i, err)
				continue
			}
			if !a.DryRun {
				log.Printf("Withdrawn invitation for %s (%s)", name, time)
				// Adding person to the blacklist to prevent sending automatic invites to them in the future
				if err := cfg.AppendToList("blacklist", name); err != nil {
					log.Printf("could not add %s to the blacklist: %v", name, err)
				}
			}
			withdrawn++
		}
	}
}
