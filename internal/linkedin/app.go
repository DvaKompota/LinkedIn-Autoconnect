package linkedin

import (
	"time"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/browser"
)

type App struct {
	Browser   *browser.Browser
	LoginPage *LoginPage
	// Add other page objects as needed
}

func NewApp(headless bool, statePath string) (*App, error) {
	// Browser provisioning handles state check internally
	b, err := browser.NewBrowser(headless, statePath)
	if err != nil {
		return nil, err
	}
	page := b.Page
	return &App{
		Browser:      b,
		LoginPage:    NewLoginPage(page),
		// Initialize other page objects here
	}, nil
}

func (a *App) Sleep(duration int) {
	time.Sleep(time.Duration(duration) * time.Second)
}

func (a *App) Close() error {
	return a.Browser.Close()
}