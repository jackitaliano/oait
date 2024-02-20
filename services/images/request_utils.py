import requests
from utils.logger import logger


def get_image_from_url(url: str):
    logger.debug(f"Getting image from url: '{url}'")
    response = requests.get(url)

    if response.status_code != 200:
        error_message: str = str(response)

        print("ERROR: " + error_message)
        return None

    logger.debug(f"Response from get: '{response}'")
    return response.content

