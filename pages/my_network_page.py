from pages.base_page import BasePage
from time import sleep


class MyNetworkPage(BasePage):

    login_url = "https://www.linkedin.com/"
    manage_network = '//a[@data-control-name="manage_all_invites"]'
    sent_button = '//button[text()="Sent"]'
    settings_button = '//button[@data-control-name="invitation_settings_entrypoint"]'
    people_button = '//button[contains(., "People")]'
    pages_button = '//button[contains(., "Pages")]'
    events_button = '//button[contains(., "Events")]'
    invitation_card = '//li[contains(@class, "invitation-card")]'
    invitation_age = '//time'
    withdraw_button = '//button[@data-control-name="withdraw_single"]'
    withdraw_dialog = '//*[@data-test-modal]'
    dialog_withdraw_button = f'{withdraw_dialog}//button//span[text()="Withdraw"]'
    withdraw_success = '//li//span[text()="Invitation withdrawn"]'

    def go_to_invites_page(self, page_no):
        locator = f'//li[@data-test-pagination-page-btn="{page_no}"]'
        self.click(locator)
        self.wait_element_selected(locator)
        sleep(1)

    def get_invite_age(self, card_no):
        card = f'{self.invitation_card}[{card_no}]'
        locator = f'({card}){self.invitation_age}'
        return self.get_element_text(locator)

    def get_invite_cards_count(self):
        return self.driver.find_elements_by_xpath(self.get_element(self.invitation_card))

    def withdraw_invite(self, card_no):
        card = f'{self.invitation_card}[{card_no}]'
        locator = f'({card}){self.withdraw_button}'
        self.click(locator)
        self.wait_element_displayed(self.withdraw_dialog)
        self.wait_element_displayed(self.dialog_withdraw_button)
        self.click(self.dialog_withdraw_button)
        self.wait_element_displayed(self.withdraw_success)
