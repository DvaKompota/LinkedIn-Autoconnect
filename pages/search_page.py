from selenium.webdriver.common.by import By
from pages.base_page import BasePage
from selenium.webdriver.common.keys import Keys
from time import time
from time import sleep
import re


class SearchPage(BasePage):

    search_url = "https://www.linkedin.com/search/results/people/"
    us_1st = "?geoUrn=%5B%22103644278%22%5D&network=%5B%22F%22%5D&origin=FACETED_SEARCH"
    us_2nd = "?geoUrn=%5B%22103644278%22%5D&network=%5B%22S%22%5D&origin=FACETED_SEARCH"
    us_2nd_and_3rd = "?geoUrn=%5B%22103644278%22%5D&network=%5B%22S%22%2C%22O%22%5D&origin=FACETED_SEARCH"

    # ============================== LOCATORS ============================== #
    people_button = '//button[@aria-label="People"]'
    connections_button = '//button[contains(@aria-label, "Connections filter.")]'
    locations_button = '//button[contains(@aria-label, "Locations filter.")]'

    company_button = '//button[contains(@aria-label, "Current company filter.")]'
    add_company_input = '//input[@placeholder="Add a company"]'
    company_filter = f'{add_company_input}/ancestor::fieldset'
    company_show_results_button = f'{company_filter}//button[contains(., "Show results")]'
    company_reset_button = f'{company_filter}//button[contains(., "Reset")]'

    no_results_locator = '//*[contains(@class, "no-results")]'
    search_results_quantity_locator = '//*[@class="search-results-container"]/*[contains(@class, "pb2")]'
    search_result_card_locator = '//li[contains(@class, "result")]'
    name_on_card = '//span[@dir]/span[@aria-hidden="true"]'
    cards_default_quantity = 10
    connection_circle_locator = '//*[contains(@class, "entity-result__badge ")]'
    connect_button_locator = '//button/*[text()="Connect"]/..'
    title_and_company_locator = '//*[contains(@class, "primary-subtitle")]'

    dialog_locator = '//div[@role="dialog"]'
    dialog_heading_locator = f'{dialog_locator}//div[contains(@class, "modal__header")]'
    nice_job_text = "Nice job building your network!"
    more_than_text = "You've sent more invitations than most"
    close_to_text = "You're close to the weekly invitation limit"
    limit_text = "You’ve reached the weekly invitation limit"
    send_button_locator = f'{dialog_locator}//button/*[text()="Send"]/..'
    add_note_button_locator = f'{dialog_locator}//button/*[text()="Add a note"]/..'
    got_it_button_locator = f'{dialog_locator}//button/*[text()="Got it"]/..'
    dismiss_button_locator = f'{dialog_locator}//button[contains(@class, "dismiss")]'

    current_page_locator = '//button[@aria-current="true"]/..'
    pagination_li_locator = '//li[@data-test-pagination-page-btn]'
    pagination_button_locator = f'{pagination_li_locator}/button'
    pagination_dots_button_locator = f'//li/button/span[.="…"]/..'

    def make_search_url(self, circle: int) -> str:
        if circle == 1:
            url = f'{self.search_url}{self.us_1st}'
        elif circle == 2:
            url = f'{self.search_url}{self.us_2nd}'
        else:
            url = f'{self.search_url}{self.us_2nd_and_3rd}'
        return url

    def search_company(self, company_name: str):
        self.wait_element_displayed(self.people_button)
        self.wait_element_displayed(self.connections_button)
        self.wait_element_displayed(self.locations_button)
        self.wait_element_displayed(self.company_button)
        self.click(self.company_button)
        self.wait_element_displayed(self.add_company_input)
        self.click(self.company_reset_button) if self.is_displayed(self.company_reset_button) else None
        self.enter_text(self.add_company_input, company_name)
        sleep(1)
        self.driver.find_element(By.XPATH, self.get_element(self.add_company_input)).send_keys(Keys.DOWN)
        self.driver.find_element(By.XPATH, self.get_element(self.add_company_input)).send_keys(Keys.ENTER)
        self.wait_element_displayed(self.company_show_results_button)
        self.click(self.company_show_results_button)

    def get_search_pages_count(self):
        end_time = time() + self.driver_wait
        while time() < end_time:
            if self.is_displayed(self.no_results_locator) or self.is_displayed(self.search_results_quantity_locator):
                break
            sleep(0.5)
        if self.is_displayed(self.no_results_locator):
            return 0
        else:
            results_count = self.get_results_count()
            if results_count < 10:
                pages_count = 1
            elif results_count % 10 != 1:
                pages_count = int(results_count / self.cards_default_quantity) + 1
            else:
                pages_count = int(results_count / self.cards_default_quantity)
            return pages_count

    def get_people_count(self):
        return len(self.driver.find_elements(By.XPATH, self.get_element(self.search_result_card_locator)))

    def get_results_count(self):
        results_string = self.get_element_text(self.search_results_quantity_locator).replace(",", "")
        return int(re.search("\\d+", results_string).group(0))

    def wait_all_people_loaded(self, current_page, last_page):
        if current_page == last_page:
            self.wait_element_displayed(f'{self.search_result_card_locator}[1]')
            sleep(1)
        else:
            end_time = time() + self.driver_wait
            while time() < end_time:
                cards_count = len(self.driver.find_elements(By.XPATH, self.get_element(self.search_result_card_locator)))
                if cards_count >= 10:
                    self.wait_element_displayed(f'{self.search_result_card_locator}[10]')
                    sleep(0.5)
                    break

    def send_invites(self, company, invites_sent, data, connection_level):
        cards_count = len(self.driver.find_elements(By.XPATH, self.get_element(self.search_result_card_locator)))
        for card in range(1, cards_count + 1):
            if invites_sent == data["per_company_limit"]:
                return invites_sent
            current_card = f'{self.search_result_card_locator}[{card}]'
            connection_circle_badge = f'({current_card}){self.connection_circle_locator}'
            title_and_company = f'({current_card}){self.title_and_company_locator}'
            connect_button = f'({current_card}){self.connect_button_locator}'
            connection_level_text = self.circle_3_text if connection_level == 3 else self.circle_2_text
            if self.is_displayed(connect_button):
                badge_text = self.get_element_text(connection_circle_badge).lower()
                title_and_company_text = self.get_element_text(title_and_company).lower()
                connection_name = self.get_element_text(f'({current_card}){self.name_on_card}')
                company_is_desired = company.lower() in title_and_company_text
                title_is_desired = any([title.lower() in title_and_company_text for title in data["job_titles"]])
                connection_level_is_correct = connection_level_text.lower() in badge_text
                if connection_level_is_correct and company_is_desired and title_is_desired:
                    self.click(connect_button)
                    self.wait_element_displayed(self.dialog_locator)
                    self.click(self.send_button_locator)
                    self.check_dialog_header()
                    self.check_dialog_header()
                    print(f"Invitation sent to {connection_name} from {company} in {connection_level}d circle")
                    invites_sent += 1
        return invites_sent

    def go_to_persons_profile(self, card_no: int):
        person_name_locator = f'({self.search_result_card_locator}[{card_no}]){self.name_on_card}'
        self.click(person_name_locator)

    def go_to_next_search_page(self, last_page: int):
        self.scroll_to_bottom()
        current_page = int(self.get_element_attribute(self.current_page_locator, 'data-test-pagination-page-btn'))
        next_page_locator = f'{self.pagination_li_locator}[@data-test-pagination-page-btn={current_page + 1}]'
        if self.is_displayed(next_page_locator):
            self.click(next_page_locator)
        elif current_page != last_page:
            self.click(self.pagination_dots_button_locator)

    def check_dialog_header(self):
        end_time = time() + 1
        while time() < end_time:
            if self.is_displayed(self.dialog_locator):
                break
        if self.is_displayed(self.dialog_locator):
            heading = self.get_element_text(self.dialog_heading_locator)
            if self.limit_text in heading:
                self.close_browser()
                print("I'm deeply sorry, my Lord, but we've reached the weekly invitation limit. Let's wait 1 week.")
            elif self.nice_job_text in heading or self.more_than_text in heading or self.close_to_text in heading:
                self.click(self.got_it_button_locator)
                self.wait_element_not_displayed(self.dialog_locator)
            elif "Connect" in heading:
                self.click(self.dismiss_button_locator)
                self.wait_element_not_displayed(self.dialog_locator)
            else:
                self.wait_element_not_displayed(self.dialog_locator)
