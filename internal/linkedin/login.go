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

func NewLoginPage(page playwright.Page) *LoginPage {
	return &LoginPage{
		page:          page,
		usernameField: page.Locator("#username"),
		passwordField: page.Locator("#password"),
		loginButton:   page.Locator("button[type=submit]"),
	}
}

func (lp *LoginPage) Login(username, password string) error {
	if err := lp.usernameField.Fill(username); err != nil {
		return err
	}
	if err := lp.passwordField.Fill(password); err != nil {
		return err
	}
	return lp.loginButton.Click()
}

func (lp *LoginPage) WaitForLoad () error {
	return lp.page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
}

func (lp *LoginPage) Navigate() error {
	_, err := lp.page.Goto("https://www.linkedin.com/login")
	return err
}