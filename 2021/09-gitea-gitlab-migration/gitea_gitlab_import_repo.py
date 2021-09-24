from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import os
import time

GITLAB_URL="https://gitlab.example.com/"
GITLAB_USER="user"
GITLAB_PASSWORD="pa$$word"
GITED_URL="https://gitea.example.com/"
GITED_TOKEN="superSecret"


driver = webdriver.Firefox(os.getcwd()+os.path.sep)
driver.get(GITLAB_URL)

# GitLab login
user = driver.find_element_by_id("user_login")
user.send_keys(GITLAB_USER)
pas = driver.find_element_by_id("user_password")
pas.send_keys(GITLAB_PASSWORD)
login = driver.find_element_by_name("commit").click()

# starting import process
driver.get(GITLAB_URL+"/import/gitea/new")
gitea_host = driver.find_element_by_name("gitea_host_url")
gitea_host.send_keys(GITED_URL)
gitea_token = driver.find_element_by_name("personal_access_token")
gitea_token.send_keys(GITED_TOKEN)
process = driver.find_element_by_name("commit").click()

# iterate over table
wait = WebDriverWait(driver, 10)
table =  wait.until(EC.presence_of_element_located((By.XPATH, '//table')))
for row in table.find_elements_by_xpath(".//tr"):
  group=row.get_attribute("data-qa-source-project").split("/")[0]
  # clicking select button to show dropdown menu and activate buttons
  row.find_element_by_class_name("gl-dropdown-toggle").click()
  time.sleep(1)
  # finding project group
  for btn in row.find_elements_by_class_name("dropdown-item"):
    if btn.get_attribute("data-qa-group-name") == group:
      btn.click()
  time.sleep(1)
  # starting import
  import_button = row.find_element(By.XPATH, "//button[@data-qa-selector='import_button']")
  import_button.click()
  # wait until import is ready
  while True:
    time.sleep(10)
    status = row.find_elements_by_class_name("gl-p-4")[-1].text
    if status == "Complete":
      break

time.sleep(60)

# closing session
driver.quit()
