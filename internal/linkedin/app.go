package linkedin

import (
	"time"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/browser"
)

type App struct {
	Browser     *browser.Browser
	Login       *LoginPage
	Invitations *InvitationsPage
	Search      *SearchPage
}

// NewApp creates a new App instance with an initialized browser
func NewApp(headless bool, statePath string) (*App, error) {
	b, err := browser.NewBrowser(headless, statePath)
	if err != nil {
		return nil, err
	}
	page := b.Page
	page.SetDefaultTimeout(10000)

	return &App{
		Browser:     b,
		Login:       NewLoginPage(page),
		Invitations: NewInvitationsPage(page),
		Search:      NewSearchPage(page),
	}, nil
}

// Sleep pauses execution for the specified number of seconds
func (a *App) Sleep(duration float64) {
	time.Sleep(time.Duration(duration) * time.Second)
}

// Close shuts down the app and browser
func (a *App) Close() error {
	return a.Browser.Close()
}
