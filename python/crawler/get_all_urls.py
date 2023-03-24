import requests
from bs4 import BeautifulSoup


url = "https://www.geeksforgeeks.org/"
reqs = requests.get(url, verify=False, allow_redirects=False)
soup = BeautifulSoup(reqs.text, "html.parser")

urls = []
for link in soup.find_all("a"):
    print(link.get("href"))
