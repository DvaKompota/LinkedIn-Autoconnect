package main

import (
	"os"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/browser"
)

func main() {
	if _, err := os.Stat("data/browser-state"); os.IsNotExist(err) {
		browser.Login()
	} else {
		browser.NewBrowser("data/browser-state")
	}
}
