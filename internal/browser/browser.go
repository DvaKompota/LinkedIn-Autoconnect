package browser

import (
	"fmt"
	"os"

	"github.com/playwright-community/playwright-go"
)

type Browser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	Page    playwright.Page
}

func NewBrowser(headless bool, statePath string) (*Browser, error) {
	if _, err := os.Stat(statePath); os.IsNotExist(err) {
		// No state file, start fresh for login
		return newBrowserForSignIn(statePath)
	}
	return newBrowser(headless, statePath)
}

func newBrowserForSignIn(statePath string) (*Browser, error) {
	// Start in headed mode for manual login
	b, err := newBrowser(false, "")
	if err != nil {
		return nil, err
	}

	// Navigate to login page
	if _, err = b.Page.Goto("https://www.linkedin.com/login"); err != nil {
		b.Close()
		return nil, err
	}

	// Wait for login (profile icon to become visible, indicating a successful login)
	err = b.Page.Locator(".global-nav__me").WaitFor(playwright.LocatorWaitForOptions{
		State:   playwright.WaitForSelectorStateVisible,
		Timeout: playwright.Float(60000),
	})
	if err != nil {
		b.Close()
		return nil, fmt.Errorf("user failed to login within 60 seconds: %v", err)
	} 

	// Save the state only if login succeeds
	if err = b.SaveState(statePath); err != nil {
		b.Close()
		return nil, err
	}

	return b, nil
}

func newBrowser(headless bool, statePath string) (*Browser, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, err
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(headless),
	})
	if err != nil {
		pw.Stop()
		return nil, err
	}

	opts := playwright.BrowserNewContextOptions{}
	if statePath != "" {
		opts.StorageStatePath = playwright.String(statePath)
	}
	if !headless {
		opts.NoViewport = playwright.Bool(true)
	}

	context, err := browser.NewContext(opts)
	if err != nil {
		browser.Close()
		pw.Stop()
		return nil, err
	}

	page, err := context.NewPage()
	if err != nil {
		context.Close()
		browser.Close()
		pw.Stop()
		return nil, err
	}

	return &Browser{
		pw:      pw,
		browser: browser,
		context: context,
		Page:    page,
	}, nil
}

func (b *Browser) Close() error {
	if err := b.context.Close(); err != nil {
		return err
	}
	if err := b.browser.Close(); err != nil {
		return err
	}
	return b.pw.Stop()
}

func (b *Browser) SaveState(path string) error {
	_, err := b.context.StorageState(path)
	return err
}
