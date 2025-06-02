package feature

import (
	"log"
	"math/rand"
	"slices"

	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/config"
	"github.com/DvaKompota/LinkedIn-Autoconnect/internal/linkedin"
)

func InviteFromSearch(a *linkedin.App, cfg *config.Config) {
	// Navigate to the search page with the configured search criteria
	if cfg.SearchLevel == 1 {
		if err := a.Search.FirstCircle(); err != nil {
			log.Fatalf("Navigation to first circle search failed: %v", err)
		}
	}
	if cfg.SearchLevel == 2 {
		if err := a.Search.SecondCircle(); err != nil {
			log.Fatalf("Navigation to second circle search failed: %v", err)
		}
	}
	if err := a.Search.WaitForLoad(); err != nil {
		log.Fatalf("Search header didn't load: %v", err)
	}
	if err := a.Search.WaitForPeopleCountToBeMoreThan(10); err != nil {
		log.Fatalf("There were no contact cards in the search resutls: %v", err)
	}
	if err := a.Search.ResetFilterByCompany(); err != nil {
		log.Fatalf("Failed to reset company filter: %v", err)
	}

	// Iterate through the companies from the config.SearchList
	companies := cfg.SearchList
	totalConnected := 0
	for _, company := range companies {
		if err := a.Search.FilterByCompany(company); err != nil {
			log.Fatalf("Failed to reset company filter: %v", err)
		}
		if err := a.Search.WaitForPeopleCountToBeMoreThan(0); err != nil {
			log.Printf("There were no contact cards for %s: %v", company, err)
			continue
		}
		if err := a.Search.ScrollToBottom(); err != nil {
			log.Fatalf("could not scroll to bottom: %v", err)
		}

		connected := 0
		for connected < cfg.PerCompanyLimit {

			// Iterate through the contact cards and connect with people who meet invitation criteria
			if err := a.Search.WaitForPeopleCountToBeMoreThan(0); err != nil {
				log.Printf("There were no contact cards for %s: %v", company, err)
				continue
			}
			count, _ := a.Search.CountPeople()
			log.Printf("Found %d contact cards in %s", count, company)

			for i := range count {
				if connected >= cfg.PerCompanyLimit {
					log.Printf("Reached the limit of %d connections for %s", connected, company)
					break
				}

				// Get the contact card by index and check if it can be connected
				if bool, _ := a.Search.IsContactCardValid(i); !bool {
					log.Printf("Skipping card #%d — it is not valid", i)
					continue
				}
				contactCard, err := a.Search.GetContactCardByIndex(i)
				if err != nil {
					log.Printf("Error getting contact card #%d: %v", i, err)
					continue
				}
				if err := contactCard.ScrollIntoViewIfNeeded(); err != nil {
					log.Printf("could not scroll contact card %d into view: %v", i, err)
					continue
				}
				canConnect := contactCard.CanConnect
				blacklisted := slices.Contains(cfg.Blacklist, contactCard.Name)
				titleMatches := contactCard.TitleMatches(cfg)

				// If the contact matches all the criteria, connect with them
				if canConnect && !blacklisted && titleMatches {
					a.Sleep(1 + rand.Float64()) // Slows down the process to avoid being detected as a bot
					if err := contactCard.Connect(); err != nil {
						log.Printf("could not connect with %s: %v", contactCard.Name, err)
						continue
					}
					if err := a.Search.Confirm(); err != nil {
						log.Printf("could not confirm connection with %s: %v", contactCard.Name, err)
						continue
					}
					log.Printf("Connected with %s (%s) from %s", contactCard.Name, contactCard.Title, company)
					connected++
				} else {
					log.Printf("Skipping %s (%s) from %s - canConnect: %t, blacklisted: %t, titleMatches: %t",
						contactCard.Name, contactCard.Title, company, canConnect, blacklisted, titleMatches)
				}
			}

			// Check if there are more pages of search results to process
			if bool, _ := a.Search.IsNextButtonEnabled(); bool {
				if err := a.Search.NextPage(); err != nil {
					log.Printf("could not go to the next page: %v", err)
					break
				}
			} else {
				log.Println("No more pages to process.")
				break
			}
			a.Sleep(1) // One extra second to allow the page to load
		}
		log.Printf("Finished processing %s, connected with %d people", company, connected)
		totalConnected += connected
	}
	log.Printf("Finished processing all companies, total connected: %d", totalConnected)
}
