package linkedin

import (
	"github.com/playwright-community/playwright-go"
)

type LoginPage struct {
	page          playwright.Page
	usernameField playwright.Locator
	passwordField playwright.Locator
	loginButton   playwright.Locator
}

// NewLoginPage initializes a new LoginPage object
func NewLoginPage(page playwright.Page) *LoginPage {
	return &LoginPage{
		page:          page,
		usernameField: page.Locator("#username"),
		passwordField: page.Locator("#password"),
		loginButton:   page.Locator("button[type=submit]"),
	}
}

// Login fills in the username and password and submits (optional for manual login)
func (lp *LoginPage) Login(username, password string) error {
	if err := lp.usernameField.Fill(username); err != nil {
		return err
	}
	if err := lp.passwordField.Fill(password); err != nil {
		return err
	}
	return lp.loginButton.Click()
}

// Navigate goes to the LinkedIn login page
func (lp *LoginPage) Navigate() error {
	_, err := lp.page.Goto("https://www.linkedin.com/login")
	return err
}

// WaitForLogin waits for the user to log in manually
func (lp *LoginPage) WaitForLogin(timeout float64) error {
	return lp.page.Locator(".global-nav__me").WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(timeout),
	})
}
