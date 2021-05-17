from pages.search_page import SearchPage
from modules.base import get_data
from modules.login import login
import data.companies as search_list


data = get_data()
login(data)
page = SearchPage(data)
url = page.make_search_url()
page.open_url(url)
for company in search_list.companies:
    page.search_company(company)
    results_count = page.get_results_count()
    results_pages = int(results_count/10)
    for i in range(results_pages):
        page.connect_all_2nd()
        page.go_to_next_search_page()
page.close_browser()
