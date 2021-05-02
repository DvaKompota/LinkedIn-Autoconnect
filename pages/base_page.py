import selenium.webdriver.support.ui as ui


class BasePage:

    def __init__(self, driver):
        self.driver = driver
        self.webdriver_wait = 10

    def wait_element_displayed(self, locator):
        wait = ui.WebDriverWait(self.driver, self.webdriver_wait)
        wait.until(lambda driver: self.driver.find_element_by_xpath(locator).is_displayed())

    def is_displayed(self, locator):
        return self.driver.find_element_by_xpath(locator).is_displayed()

    def click(self, locator):
        self.wait_element_displayed(locator)
        self.driver.find_element_by_xpath(locator).click()

    def enter_text(self, locator, text):
        self.wait_element_displayed(locator)
        self.driver.find_element_by_xpath(locator).send_keys(text)
