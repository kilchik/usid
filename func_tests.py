from selenium import webdriver


def setup_driver_for_main_page():
    driver = webdriver.Chrome()
    driver.get("http://127.0.0.1:9632/usid")
    return driver


def test_main_page_contains_input():
    driver = setup_driver_for_main_page()
    elem = driver.find_element_by_id("input")
    assert elem
    driver.close()
