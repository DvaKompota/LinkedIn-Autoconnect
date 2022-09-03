from pages.search_page import SearchPage
from pages.profile_page import ProfilePage
from modules.base import get_data
from modules.login import login


data = get_data()
login(data)
page = SearchPage(data)
url = page.make_search_url(1)
for company in data["search_list"]:
    page.open_url(url)
    print(f'Searching for employees of {company}')
    page.search_company(company)
    search_pages_count = page.get_search_pages_count()
    results_count = page.get_results_count()
    print(f'Search returned {results_count} people from {company} within 1 circle shown on {search_pages_count} pages')
    invites_sent = 0
    for page_no in range(1, search_pages_count + 1):
        if invites_sent == data["per_company_limit"]:
            print(f'Already sent {invites_sent} invites to {company} employees which is maximum for one company')
            break
        people_count = page.get_people_count()
        for person in range(1, people_count + 1):
            if invites_sent == data["per_company_limit"]:
                print(f'Already sent {invites_sent} invites to {company} employees which is maximum for one company')
                page.go_back()
                page = SearchPage(data)
                break
            page.wait_all_people_loaded(page_no, search_pages_count)
            page.go_to_persons_profile(person)
            page = ProfilePage(data)
            people_shown = page.show_more_people_also_viewed()
            if people_shown is False:
                page.go_back()
                page = SearchPage(data)
                continue
            invites_sent = page.send_invites(company, invites_sent, data, connection_level=2)
            if data["connection_level"] == 3:
                invites_sent = page.send_invites(company, invites_sent, data, connection_level=3)
            page.go_back()
            page = SearchPage(data)
page.close_browser()
