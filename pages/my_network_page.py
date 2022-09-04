from selenium.webdriver.common.by import By
from pages.base_page import BasePage
from time import sleep


class MyNetworkPage(BasePage):

    login_url = "https://www.linkedin.com/"
    manage_network = '//a[@href="/mynetwork/invitation-manager/"]'
    sent_button = '//button[text()="Sent"]'
    settings_button = '//button[@data-control-name="invitation_settings_entrypoint"]'
    people_button = '//button[contains(., "People")]'
    pages_button = '//button[contains(., "Pages")]'
    events_button = '//button[contains(., "Events")]'
    invitation_card = '//li[contains(@class, "invitation-card")]'
    invitation_age = '//time'
    invitation_name = '//*[contains(@class, "invitation-card__title")]'
    withdraw_button = '//button[contains(., "Withdraw")]'
    withdraw_dialog = '//*[@data-test-modal]'
    dialog_withdraw_button = f'{withdraw_dialog}//button//span[text()="Withdraw"]'
    withdraw_success = '//*[text()="Invitation withdrawn"]'

    def go_to_invites_page(self, page_no):
        locator = f'//li[@data-test-pagination-page-btn="{page_no}"]'
        self.click(locator)
        self.wait_element_selected(locator)
        sleep(1)

    def get_invite_age(self, card_no):
        card = f'{self.invitation_card}[{card_no}]'
        locator = f'({card}){self.invitation_age}'
        return self.get_element_text(locator)

    def get_invite_name(self, card_no):
        card = f'{self.invitation_card}[{card_no}]'
        locator = f'({card}){self.invitation_name}'
        return self.get_element_text(locator)

    def get_invite_cards_count(self):
        return self.driver.find_elements(By.XPATH, self.get_element(self.invitation_card))

    def withdraw_invite(self, card_no):
        card = f'{self.invitation_card}[{card_no}]'
        locator = f'({card}){self.withdraw_button}'
        self.click(locator)
        self.wait_element_displayed(self.withdraw_dialog)
        self.wait_element_displayed(self.dialog_withdraw_button)
        self.click(self.dialog_withdraw_button)
        self.wait_element_displayed(self.withdraw_success)

    def add_to_blacklist(self, invite_name):
        file = open("../data/blacklist.txt", "r+")
        lines = file.readlines()
        blacklisted = []
        for line in lines:
            line = line.replace("\n", "")
            blacklisted.append(line)
        file.write(f'{invite_name}\n') if invite_name not in blacklisted else None
        file.close()
