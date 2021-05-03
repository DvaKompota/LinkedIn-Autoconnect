import selenium.webdriver.support.ui as ui
from selenium.common.exceptions import StaleElementReferenceException
from time import time
from time import sleep


class BasePage:

    def __init__(self, data):
        self.driver = data["driver"]
        self.driver_wait = 10

    def get_element(self, locator_name):
        try:
            return getattr(self, locator_name)
        except AttributeError:
            return locator_name

    def open_url(self, url):
        self.driver.get(self.get_element(url))

    def close_browser(self):
        self.driver.close()

    def wait_element_displayed(self, locator):
        wait = ui.WebDriverWait(self.driver, self.driver_wait)
        wait.until(lambda driver: self.driver.find_element_by_xpath(self.get_element(locator)).is_displayed())

    def wait_element_not_displayed(self, locator):
        end_time = time() + self.driver_wait
        while time() < end_time:
            if not self.is_displayed(locator):
                break

    def wait_element_selected(self, locator):
        end_time = time() + self.driver_wait
        while time() < end_time:
            try:
                if self.is_selected(locator):
                    break
            except StaleElementReferenceException:
                sleep(0.5)
                if self.is_selected(locator):
                    break

    def is_displayed(self, locator):
        return self.driver.find_element_by_xpath(self.get_element(locator)).is_displayed()

    def is_selected(self, locator):
        return "selected" in self.driver.find_element_by_xpath(self.get_element(locator)).get_attribute("class")

    def click(self, locator):
        self.wait_element_displayed(locator)
        self.driver.find_element_by_xpath(self.get_element(locator)).click()

    def enter_text(self, locator, text):
        self.wait_element_displayed(locator)
        self.driver.find_element_by_xpath(self.get_element(locator)).send_keys(text)

    def get_element_text(self, locator):
        self.wait_element_displayed(locator)
        return self.driver.find_element_by_xpath(self.get_element(locator)).text
