import selenium.webdriver as webdriver


def get_driver():
    opts = webdriver.ChromeOptions()
    opts.add_argument("--start-maximized")
    driver = webdriver.Chrome(options=opts)
    return driver
