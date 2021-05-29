from selenium.common.exceptions import TimeoutException
from pages.search_page import SearchPage


class ProfilePage(SearchPage):

    people_also_viewed_section = '//*[@class="pv-browsemap-section"]'
    show_more_button = '//button//*[text()="Show more"]'

    search_result_card_locator = f'{people_also_viewed_section}//li[contains(@class, "member")]'
    name_on_card = '//*[contains(@class, "name ")]'
    connection_circle_locator = '//*[contains(@class, "dist-value")]'
    title_and_company_locator = '//*[contains(@class, "member-headline")]'
    connect_button_locator = '//button/*[text()="Connect"]/..'

    def show_more_people_also_viewed(self):
        try:
            self.wait_element_displayed(f'{self.people_also_viewed_section}')
            button = f'{self.people_also_viewed_section}/..{self.show_more_button}'
            self.scroll_to(100)
            self.click(button)
            return True
        except TimeoutException:
            return False
