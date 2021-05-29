from selenium.webdriver.common.action_chains import ActionChains
from selenium.common.exceptions import StaleElementReferenceException
from selenium.common.exceptions import NoSuchElementException
import selenium.webdriver.support.ui as ui
from time import time
from time import sleep


class BasePage:

    def __init__(self, data):
        self.driver = data["driver"]
        self.driver_wait = data["driver_wait"]

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
        try:
            wait.until(lambda driver: self.driver.find_element_by_xpath(self.get_element(locator)).is_displayed())
        except StaleElementReferenceException:
            sleep(0.5)
            wait.until(lambda driver: self.driver.find_element_by_xpath(self.get_element(locator)).is_displayed())

    def wait_element_not_displayed(self, locator):
        end_time = time() + self.driver_wait
        while time() < end_time:
            try:
                if not self.is_displayed(locator):
                    break
            except StaleElementReferenceException:
                sleep(0.5)
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
        try:
            return self.driver.find_element_by_xpath(self.get_element(locator)).is_displayed()
        except NoSuchElementException:
            return False

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

    def get_element_attribute(self, locator, attribute):
        self.wait_element_displayed(locator)
        return self.driver.find_element_by_xpath(self.get_element(locator)).get_attribute(attribute)

    def scroll_to_element(self, locator):
        self.wait_element_displayed(locator)
        element = self.driver.find_element_by_xpath(self.get_element(locator))
        ActionChains(self.driver).move_to_element(element).perform()

    def scroll_to_bottom(self):
        self.driver.execute_script("window.scrollTo(0, document.body.scrollHeight);")
