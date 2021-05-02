from selenium import webdriver
from data.credentials import email
from data.credentials import password


def get_data():
    data = {"driver": get_driver(), "email": email, "password": password}
    return data


def get_driver():
    opts = webdriver.ChromeOptions()
    opts.add_argument("--start-maximized")
    driver = webdriver.Chrome(options=opts)
    return driver
