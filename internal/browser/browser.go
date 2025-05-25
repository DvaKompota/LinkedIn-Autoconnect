package browser

import (
	"github.com/playwright-community/playwright-go"
)

type Browser struct {
	pw      *playwright.Playwright
	browser playwright.Browser
	context playwright.BrowserContext
	Page    playwright.Page
}

// NewBrowser creates a new browser instance with the given options
func NewBrowser(headless bool, statePath string) (*Browser, error) {
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

// Close shuts down the browser and cleans up resources
func (b *Browser) Close() error {
	if err := b.context.Close(); err != nil {
		return err
	}
	if err := b.browser.Close(); err != nil {
		return err
	}
	return b.pw.Stop()
}

// SaveState saves the browser's state to a file
func (b *Browser) SaveState(path string) error {
	_, err := b.context.StorageState(path)
	return err
}
