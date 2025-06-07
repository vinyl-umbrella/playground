import logging
import requests


def get_http_status(code: int) -> requests.Response:
    url = f"https://httpbin.org/status/{code}"
    return requests.get(url)


def main():
    logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")

    status_ok = get_http_status(200)
    logging.info(f"Status code 200: {status_ok.status_code}")

    status_not_found = get_http_status(404)
    logging.info(f"Status code 404: {status_not_found.status_code}")


if __name__ == "__main__":
    main()
