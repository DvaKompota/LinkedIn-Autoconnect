from pages.base_page import BasePage


class LoginPage(BasePage):

    login_url = "https://www.linkedin.com/"

    # ============================== LOCATORS ============================== #
    email_field = '//input[@id="session_key"]'
    pass_field = '//input[@id="session_password"]'
    submit_button = '//button[@type="submit"]'
    avatar = '//img[contains(@alt, "Photo of ")]'
    my_network = '//a[@href="https://www.linkedin.com/mynetwork/?"]'
