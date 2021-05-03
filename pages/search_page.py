from pages.base_page import BasePage
from selenium.webdriver.common.keys import Keys
from time import sleep


class SearchPage(BasePage):

    search_2nd_us_only = "https://www.linkedin.com/search/results/people/" \
                         "?geoUrn=%5B%22103644278%22%5D&network=%5B%22S%22%5D&origin=FACETED_SEARCH"

    # ============================== LOCATORS ============================== #
    people_button = '//button[@aria-label="People"]'
    connections_button = '//button[contains(@aria-label, "Connections filter.")]'
    locations_button = '//button[contains(@aria-label, "Locations filter.")]'

    company_button = '//button[contains(@aria-label, "Current company filter.")]'
    add_company_input = '//input[@placeholder="Add a company"]'
    company_filter = f'{add_company_input}/ancestor::fieldset'
    company_show_results_button = f'{company_filter}//button[contains(., "Show results")]'
    company_reset_button = f'{company_filter}//button[contains(., "Reset")]'

    def search_company(self, company_name):
        self.wait_element_displayed(self.people_button)
        self.wait_element_displayed(self.connections_button)
        self.wait_element_displayed(self.locations_button)
        self.wait_element_displayed(self.company_button)
        self.click(self.company_button)
        self.wait_element_displayed(self.add_company_input)
        self.click(self.company_reset_button) if self.is_displayed(self.company_reset_button) else None
        self.enter_text(self.add_company_input, company_name)
        sleep(1)
        self.driver.find_element_by_xpath(self.get_element(self.add_company_input)).send_keys(Keys.DOWN)
        self.driver.find_element_by_xpath(self.get_element(self.add_company_input)).send_keys(Keys.ENTER)
        self.wait_element_displayed(self.company_show_results_button)
        self.click(self.company_show_results_button)
