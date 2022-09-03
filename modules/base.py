from selenium import webdriver
from data import config
from data import credentials
import chromedriver_autoinstaller


def get_data():
    data = {
        "driver": get_driver(config.headless),
        "driver_wait": config.driver_wait,
        "email": credentials.email,
        "password": credentials.password,
        "search_level": config.search_level,
        "connection_level": config.connection_level,
        "per_company_limit": config.per_company_limit,
        "search_list": config.search_list,
        "job_titles": config.job_titles,
        }
    return data


def get_driver(headless=True):
    chromedriver_autoinstaller.install()
    opts = webdriver.ChromeOptions()
    opts.add_argument("--start-maximized")
    opts.add_argument("--headless") if headless else None
    driver = webdriver.Chrome(options=opts)
    return driver
