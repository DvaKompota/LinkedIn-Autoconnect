from pages.base_page import BasePage
from pages.login_page import LoginPage
from pages.my_network_page import MyNetworkPage
from modules.base import get_data
from modules.login import login
import re


data = get_data()
login(data)

page = LoginPage(data)
page.wait_element_displayed("my_network")
page.click("my_network")

page = MyNetworkPage(data)
page.wait_element_displayed("manage_network")
page.click("manage_network")
page.wait_element_displayed("sent_button")
page.click("sent_button")
page.wait_element_displayed("pages_button")
page.wait_element_displayed("events_button")
page.wait_element_displayed("people_button")
invites_count = int(re.search(r"\d+", page.get_element_text("people_button"))[0])
for invite_page in range(int(invites_count / 100) + 1, 0, -1):
    page.go_to_invites_page(invite_page) if invite_page > 1 else None
    cards_count = len(page.get_invite_cards_count())
    for invite_card in range(cards_count, 0, -1):
        if "month" in page.get_invite_age(invite_card):
            page.withdraw_invite(invite_card)
BasePage(data).close_browser()
