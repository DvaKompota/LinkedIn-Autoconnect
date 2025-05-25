package main

import (
	"log"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/linkedin"
)

func main() {
	a, err := linkedin.NewApp(false, "data/browser-state")
	if err != nil {
		log.Fatalf("could not create app: %v", err)
	}
	defer a.Close()

	a.LoginPage.Navigate()
	a.LoginPage.WaitForLoad()
	a.Sleep(1)
}
