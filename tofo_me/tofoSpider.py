# coding:utf-8
from selenium import webdriver
import re
from selenium.webdriver.common.by import By
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.support.ui import WebDriverWait
import os
import urllib2
driver = webdriver.Chrome()

driver.get('https://tofo.me/')

file_base_url = os.path.dirname(__file__)
print file_base_url
driver.find_element_by_class_name("card").click()
imageRoot = 'https://x.gto.cc/'
user_agent = 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.181 Safari/537.36'

while True:
    txt = driver.find_element_by_class_name("modals").get_attribute("innerHTML")
    btns = driver.find_elements_by_class_name('black')
    if len(btns) > 1:
        btns[1].click()
    for item in re.findall('<(img|video) src="(https://x.gto.cc/.*?)"\\s+', txt):
        try:
            wait = WebDriverWait(driver, 10)
            element = wait.until(EC.presence_of_element_located((By.CLASS_NAME, 'modals')))
            print item[1][-10:0]
            print item[1]
            req = urllib2.Request(url=item[1])
            req.headers = {
                'user-agent': user_agent,
                'Referer': 'https://tofo.me/'
            }

            if item[1][-4:] != '.mp4':
                response = urllib2.urlopen(req, timeout=2)
            else:
                response = urllib2.urlopen(req)
            with open(file_base_url+'/images/'+item[1][-20:], "wb") as f:
                f.write(response.read())
        except:
            continue



