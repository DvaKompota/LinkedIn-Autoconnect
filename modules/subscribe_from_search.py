from pages.search_page import SearchPage
from modules.base import get_data
from modules.login import login
from warnings import simplefilter


data = get_data()
login(data)
page = SearchPage(data)
url = page.make_search_url(data["search_level"])
page.open_url(url)
for company in data["search_list"]:
    print(f'Searching for employees of {company}')
    page.search_company(company)
    search_pages_count = page.get_search_pages_count()
    print(f'Search returned {search_pages_count} pages of potential contacts within the {data["search_level"]} circle')
    invites_sent = 0
    for page_no in range(1, search_pages_count + 1):
        if invites_sent == data["per_company_limit"]:
            print(f'Already sent {invites_sent} invites to {company} employees which is maximum for one company')
            break
        page.wait_all_people_loaded(page_no, search_pages_count)
        invites_sent = page.send_invites(company, invites_sent, data["per_company_limit"], connection_level=2)
        if data["connection_level"] == 3:
            invites_sent = page.send_invites(company, invites_sent, data["per_company_limit"], connection_level=3)
        page.go_to_next_search_page(search_pages_count)
page.close_browser()
simplefilter("ignore", ResourceWarning)
