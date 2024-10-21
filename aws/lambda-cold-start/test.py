import time
import requests


while True:
    url = ""
    res = requests.get(url)
    if res.status_code != 200:
        print(f"Error: {res.status_code}")

    print(res.json())
    time.sleep(1)
