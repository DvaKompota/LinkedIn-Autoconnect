import selenium.webdriver.support.ui as ui


def wait_for_element_displayed(driver, locator):
    wait = ui.WebDriverWait(driver, 10)
    wait.until(lambda driver: driver.find_element_by_xpath(locator).is_displayed())


def is_displayed(driver, locator):
    return driver.find_element_by_xpath(locator).is_displayed()

