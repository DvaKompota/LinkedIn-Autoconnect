package linkedin

import (
	"fmt"
	"time"

	"github.com/playwright-community/playwright-go"
)

type InvitationsPage struct {
	page                  playwright.Page
	url                   string
	header                playwright.Locator
	received              playwright.Locator
	sent                  playwright.Locator
	invitation            playwright.Locator
	name                  playwright.Locator
	timeBadge             playwright.Locator
	withdraw              playwright.Locator
	withdrawDialog        playwright.Locator
	confirmWithdraw       playwright.Locator
	withdrawnConfirmation playwright.Locator
}

// NewInvitationsPage initializes a new InvitationsPage object
func NewInvitationsPage(page playwright.Page) *InvitationsPage {
	p := &InvitationsPage{
		page: page,
	}
	p.url = "https://www.linkedin.com/mynetwork/invitation-manager/"
	p.header = page.GetByText("Manage invitations")

	// First iteration of locators (keep in case LinkedIn reverts to this strategy)
	// p.received = page.GetByRole(playwright.AriaRole("tab"), playwright.PageGetByRoleOptions{Name: "Received"})
	// p.sent = p.page.GetByRole(playwright.AriaRole("tab"), playwright.PageGetByRoleOptions{Name: "Sent"})
	// p.invitation = page.Locator("ul li.invitation-card")   // Targets each invitation card in the list
	// p.name = p.page.Locator(".invitation-card__tvm-title") // Targets the name within each invitation card
	// p.timeBadge = p.page.Locator(".time-badge")            // Targets the time badge within each invitation card

	// Second iteration of locators (currently working)
	p.received = page.GetByRole(playwright.AriaRole("link"), playwright.PageGetByRoleOptions{Name: "Received"})
	p.sent = p.page.GetByRole(playwright.AriaRole("link"), playwright.PageGetByRoleOptions{Name: "Sent"})
	p.invitation = page.Locator(`[componentkey="InvitationManagerPage_InvitationsList"]`).GetByRole("listitem") // Targets each invitation card in the list
	p.name = p.page.Locator("p a").First()                                                                      // Targets the name within each invitation card
	p.timeBadge = p.page.GetByText(" ago")                                                                      // Targets the time badge within each invitation card

	p.withdraw = p.page.GetByRole(playwright.AriaRole("button"), playwright.PageGetByRoleOptions{Name: "Withdraw"})
	p.withdrawDialog = p.page.GetByRole(playwright.AriaRole("alertdialog"), playwright.PageGetByRoleOptions{Name: "Withdraw"})
	p.confirmWithdraw = p.withdrawDialog.GetByRole(playwright.AriaRole("button"), playwright.LocatorGetByRoleOptions{Name: "Withdraw"})
	p.withdrawnConfirmation = p.page.GetByRole(playwright.AriaRole("alert")).GetByText("withdrawn")
	return p
}

// Navigate goes to the LinkedIn login page
func (p *InvitationsPage) Navigate() error {
	_, err := p.page.Goto(p.url)
	return err
}

// WaitForLoad waits for Manage Invitations header to be visible
func (p *InvitationsPage) WaitForLoad() error {
	err := p.page.WaitForLoadState()
	if err != nil {
		return err
	}
	return p.header.WaitFor()
}

// Open Sent Invitations tab
func (p *InvitationsPage) OpenSentTab() error {
	if err := p.sent.Click(); err != nil {
		return fmt.Errorf("failed to open sent invitations tab: %w", err)
	}
	return nil
}

// Waits for the invitations to load
func (p *InvitationsPage) WaitForInvitationsCountToBeMoreThan(count int) error {
	now := time.Now()
	for time.Since(now) < 10*time.Second { // Wait up to 10 seconds
		currentCount, err := p.invitation.Count()
		if err != nil {
			return fmt.Errorf("failed to count invitations: %w", err)
		}
		if currentCount > count {
			return nil // Invitations loaded
		}
		time.Sleep(500 * time.Millisecond) // Wait before retrying
	}
	return fmt.Errorf("invitations did not load within 10 seconds")
}

// Counts invitations
func (p *InvitationsPage) CountInvitations() (int, error) {
	count, err := p.invitation.Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetInvitationLocatorByIndex retrieves the Locator for the specified invitation index
func (p *InvitationsPage) GetInvitationLocatorByIndex(index int) playwright.Locator {
	return p.invitation.Nth(index)
}

// GetInvitationName retrieves the name from the specified invitation Locator
func (p *InvitationsPage) GetInvitationName(invitation playwright.Locator) (string, error) {
	name, err := invitation.Locator(p.name).InnerText()
	if err != nil {
		return "", fmt.Errorf("failed to get name for invitation: %w", err)
	}
	return name, nil
}

// GetInvitationTime retrieves the time from the specified invitation Locator
func (p *InvitationsPage) GetInvitationTime(invitation playwright.Locator) (string, error) {
	timeBadge, err := invitation.Locator(p.timeBadge).InnerText()
	if err != nil {
		return "", fmt.Errorf("failed to get time for invitation: %w", err)
	}
	return timeBadge, nil
}

// WithdrawInvitation clicks the Withdraw button for the specified invitation Locator
func (p *InvitationsPage) WithdrawInvitation(invitation playwright.Locator) error {
	withdrawButton := invitation.Locator(p.withdraw)
	if err := withdrawButton.Click(); err != nil {
		return fmt.Errorf("failed to click withdraw button: %w", err)
	}
	if err := p.withdrawDialog.WaitFor(); err != nil {
		return fmt.Errorf("withdraw dialog did not appear: %w", err)
	}
	if err := p.confirmWithdraw.Click(); err != nil {
		return fmt.Errorf("failed to confirm withdrawal: %w", err)
	}
	return nil
}

// ScrollToBottom scrolls to the bottom of the invitations list
func (p *InvitationsPage) ScrollToBottom() error {
	_, err := p.page.Evaluate("window.scrollTo(0, document.body.scrollHeight)")
	return err
}
