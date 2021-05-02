import selenium.webdriver.support.ui as ui


webdriver_wait = 10


def wait_for_element_displayed(driver, locator):
    wait = ui.WebDriverWait(driver, webdriver_wait)
    wait.until(lambda driver: driver.find_element_by_xpath(locator).is_displayed())


def is_displayed(driver, locator):
    return driver.find_element_by_xpath(locator).is_displayed()


def click(driver, locator):
    wait_for_element_displayed(driver, locator)
    driver.find_element_by_xpath(locator).click()


def enter_text(driver, locator, text):
    wait_for_element_displayed(driver, locator)
    driver.find_element_by_xpath(locator).send_keys(text)


