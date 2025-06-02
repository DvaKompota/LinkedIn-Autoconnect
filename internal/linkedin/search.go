package linkedin

import (
	"fmt"
	"strings"
	"time"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/config"
	"github.com/playwright-community/playwright-go"
)

type SearchPage struct {
	page                  playwright.Page
	url                   string
	us1stCircle           string
	us2ndCircle           string
	header                playwright.Locator
	contactCard           playwright.Locator
	contactName           playwright.Locator
	contactTitle          playwright.Locator
	contactCircle         playwright.Locator
	connectButton         playwright.Locator
	companyFilter         playwright.Locator
	filterDropdown        playwright.Locator
	addCompanyField       playwright.Locator
	filterLabel           playwright.Locator
	filterSearchResult    playwright.Locator
	filterResetButton     playwright.Locator
	showResultsButton     playwright.Locator
	confirmationModal     playwright.Locator
	sendWithoutNoteButton playwright.Locator
	previousButton        playwright.Locator
	nextButton            playwright.Locator
}

//   var previousButton = await page.getByRole('button', {name: "Previous"})
//   var nextButton = await page.getByRole('button', {name: "Next"})

type ContactCard struct {
	locator       playwright.Locator
	Name          string
	Title         string
	Circle        string
	CanConnect    bool
	connectButton playwright.Locator
}

// NewSearchPage initializes a new SearchPage object
func NewSearchPage(page playwright.Page) *SearchPage {
	p := &SearchPage{
		page: page,
	}
	p.url = "https://www.linkedin.com/search/results/people/"
	p.us1stCircle = "?geoUrn=%5B%22103644278%22%5D&network=%5B%22F%22%5D&origin=FACETED_SEARCH"
	p.us2ndCircle = "?geoUrn=%5B%22103644278%22%5D&network=%5B%22S%22%5D&origin=FACETED_SEARCH"
	p.header = page.GetByLabel("Search filters")

	p.contactCard = page.Locator(".search-results-container ul[role='list'] li")
	p.contactName = page.Locator("//span[@dir]/span[@aria-hidden='true']")
	p.contactTitle = page.Locator(".mb1 > :nth-child(2)")
	p.contactCircle = page.Locator("//*[contains(@class, 'entity-result__badge ')]")
	p.connectButton = page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Connect"})
	p.previousButton = page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Previous"})
	p.nextButton = page.GetByRole("button", playwright.PageGetByRoleOptions{Name: "Next"})

	p.companyFilter = page.Locator("#searchFilter_currentCompany")
	p.filterDropdown = page.GetByRole("tooltip")
	p.addCompanyField = p.filterDropdown.GetByPlaceholder("Add a company")
	p.filterLabel = p.filterDropdown.GetByRole("listitem")
	p.filterSearchResult = p.filterDropdown.GetByRole("listbox").GetByRole("option")
	p.filterResetButton = p.filterDropdown.GetByRole("button", playwright.LocatorGetByRoleOptions{Name: "Reset"})
	p.showResultsButton = p.filterDropdown.GetByRole("button", playwright.LocatorGetByRoleOptions{Name: "Show results"})

	p.confirmationModal = page.GetByRole("dialog")
	p.sendWithoutNoteButton = p.confirmationModal.GetByRole("button", playwright.LocatorGetByRoleOptions{Name: "Send without a note"})
	return p
}

// FirstCircle goes to LinkedIn search within 1st-circle contacts in US
func (p *SearchPage) FirstCircle() error {
	url := p.url + p.us1stCircle
	_, err := p.page.Goto(url)
	return err
}

// SecondCircle goes to LinkedIn search within 2nd-circle contacts in US
func (p *SearchPage) SecondCircle() error {
	url := p.url + p.us2ndCircle
	_, err := p.page.Goto(url)
	return err
}

// WaitForLoad waits for search filters to be visible
func (p *SearchPage) WaitForLoad() error {
	err := p.page.WaitForLoadState()
	if err != nil {
		return err
	}
	return p.header.WaitFor()
}

// Counts contact cards
func (p *SearchPage) CountPeople() (int, error) {
	count, err := p.contactCard.Count()
	if err != nil {
		return 0, err
	}
	return count, nil
}

// FilterByCompany filters search results by company name
func (p *SearchPage) ResetFilterByCompany() error {
	buttonText, _ := p.companyFilter.InnerText()
	if buttonText == "Current company" {
		return nil // No filter applied, nothing to reset
	}
	if err := p.companyFilter.Click(); err != nil {
		return fmt.Errorf("failed to open company filter: %w", err)
	}
	if err := p.filterDropdown.WaitFor(); err != nil {
		return fmt.Errorf("company filter dropdown didn't appear: %w", err)
	}
	if resetVisible, _ := p.filterResetButton.IsVisible(); resetVisible {
		if err := p.filterResetButton.Click(); err != nil {
			return fmt.Errorf("failed to reset company filter: %w", err)
		}
	}
	return nil
}

// FilterByCompany filters search results by company name
func (p *SearchPage) FilterByCompany(company string) error {
	if err := p.ResetFilterByCompany(); err != nil {
		return err
	}
	if opened, _ := p.filterDropdown.IsVisible(); !opened {
		if err := p.companyFilter.Click(); err != nil {
			return fmt.Errorf("failed to open company filter: %w", err)
		}
	}
	label := p.filterLabel.GetByText(company, playwright.LocatorGetByTextOptions{Exact: playwright.Bool(true)})
	if labelVisible, _ := label.IsVisible(); labelVisible {
		if err := label.Click(); err != nil {
			return fmt.Errorf("failed to click company by label: %w", err)
		}
	} else {
		if err := p.addCompanyField.Fill(company); err != nil {
			return fmt.Errorf("failed to fill company name: %w", err)
		}
		if err := p.filterSearchResult.First().Click(); err != nil {
			return fmt.Errorf("failed to click on the top company in search: %w", err)
		}
	}
	if err := p.showResultsButton.Click(); err != nil {
		return fmt.Errorf("failed to click show results button: %w", err)
	}
	return nil
}

// Waits for contact cards to load
func (p *SearchPage) WaitForPeopleCountToBeMoreThan(count int) error {
	now := time.Now()
	for time.Since(now) < 5*time.Second { // Wait up to 5 seconds
		currentCount, err := p.contactCard.Count()
		if err != nil {
			return fmt.Errorf("failed to count invitations: %w", err)
		}
		if currentCount > count {
			return nil // Contact cards loaded
		}
		time.Sleep(500 * time.Millisecond) // Wait before retrying
	}
	return fmt.Errorf("contacts did not load within 5 seconds")
}

func (p *SearchPage) IsContactCardValid(index int) (bool, error) {
	contactCard := p.contactCard.Nth(index)
	if contactCard == nil {
		return false, fmt.Errorf("contact card at index %d not found", index)
	}
	hasName, err := contactCard.Locator(p.contactName).IsVisible(playwright.LocatorIsVisibleOptions{Timeout: playwright.Float(1000)})
	if err != nil {
		return false, fmt.Errorf("failed to check if name is visible for contact card at index %d: %w", index, err)
	}
	return hasName, nil
}

// GetContactCardByIndex retrieves a contact card by its index
func (p *SearchPage) GetContactCardByIndex(index int) (ContactCard, error) {
	contactCard := p.contactCard.Nth(index)
	if contactCard == nil {
		return ContactCard{}, fmt.Errorf("contact card at index %d not found", index)
	}
	name, err := contactCard.Locator(p.contactName).TextContent(playwright.LocatorTextContentOptions{Timeout: playwright.Float(1000)})
	if err != nil {
		return ContactCard{}, fmt.Errorf("failed to get name for contact card at index %d: %w", index, err)
	}
	title, err := contactCard.Locator(p.contactTitle).First().TextContent(playwright.LocatorTextContentOptions{Timeout: playwright.Float(1000)})
	if err != nil {
		return ContactCard{}, fmt.Errorf("failed to get title for contact card at index %d: %w", index, err)
	}
	circle, err := contactCard.Locator(p.contactCircle).TextContent(playwright.LocatorTextContentOptions{Timeout: playwright.Float(1000)})
	if err != nil {
		return ContactCard{}, fmt.Errorf("failed to get circle for contact card at index %d: %w", index, err)
	}
	connectButton := contactCard.Locator(p.connectButton)
	canConnect, _ := connectButton.IsVisible()
	return ContactCard{
		locator:       contactCard,
		Name:          strings.TrimSpace(name),
		Title:         strings.TrimSpace(title),
		Circle:        strings.TrimSpace(circle),
		CanConnect:    canConnect,
		connectButton: connectButton,
	}, nil
}

// ScrollToBottom scrolls to the bottom of the invitations list
func (p *SearchPage) ScrollToBottom() error {
	_, err := p.page.Evaluate("window.scrollTo(0, document.body.scrollHeight)")
	return err
}

func (p *SearchPage) WaitForConfirmationModal() error {
	_, err := p.page.Evaluate("window.scrollTo(0, document.body.scrollHeight)")
	return err
}

func (p *SearchPage) IsNextButtonEnabled() (bool, error) {
	if err := p.nextButton.WaitFor(playwright.LocatorWaitForOptions{Timeout: playwright.Float(1000)}); err != nil {
		return false, fmt.Errorf("next button is not present: %w", err)
	}
	enabled, err := p.nextButton.IsEnabled()
	if err != nil {
		return false, fmt.Errorf("failed to check if next button is enabled: %w", err)
	}
	return enabled, nil
}

func (p *SearchPage) NextPage() error {
	if err := p.nextButton.Click(); err != nil {
		return fmt.Errorf("failed to click next button: %w", err)
	}
	return nil
}

// ScrollIntoViewIfNeeded scrolls the contact card into view if needed
func (c *ContactCard) ScrollIntoViewIfNeeded() error {
	if c == nil {
		return fmt.Errorf("contact card is nil")
	}
	if err := c.locator.ScrollIntoViewIfNeeded(); err != nil {
		return fmt.Errorf("failed to scroll contact card into view: %w", err)
	}
	return nil
}

func (c *ContactCard) TitleMatches(cfg *config.Config) bool {
	if c == nil {
		return false
	}
	for _, searchedTitle := range cfg.JobTitles {
		if strings.Contains(c.Title, searchedTitle) {
			return true
		}
	}
	return false
}

func (c *ContactCard) Connect() error {
	if c == nil {
		return fmt.Errorf("contact card is nil")
	}
	if !c.CanConnect {
		return fmt.Errorf("cannot connect with %s, button not visible", c.Name)
	}
	if err := c.connectButton.Click(); err != nil {
		return fmt.Errorf("failed to click connect button for %s: %w", c.Name, err)
	}
	return nil
}

// Confirm confirms the connection request in the confirmation modal if it appears
func (p *SearchPage) Confirm() error {
	// TODO: add processing of other modals:
	// nice_job_text = "Nice job building your network!"
	// more_than_text = "You've sent more invitations than most"
	// close_to_text = "You're close to the weekly invitation limit"
	// limit_text = "You’ve reached the weekly invitation limit"

	if err := p.confirmationModal.WaitFor(playwright.LocatorWaitForOptions{Timeout: playwright.Float(3000)}); err != nil {
		return fmt.Errorf("confirmation modal did not appear: %w", err)
	}
	if err := p.sendWithoutNoteButton.Click(); err != nil {
		return fmt.Errorf("confirmation modal appeared, but failed to click send without note button: %w", err)
	}
	return nil
}
