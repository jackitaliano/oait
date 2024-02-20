import requests
import json

from utils.logger import logger

def get_file_by_id(key: str, file_id: str) -> dict | None:
    logger.debug(f"Getting OpenAI file by file_id: '{file_id}'")

    req_url: str = f"https://api.openai.com/v1/files/{file_id}/content"

    headers: dict[str, str] = {
        "Authorization": f"Bearer {key}",
    }

    response: requests.Response = requests.get(url=req_url, headers=headers)
    res = response.content

    if response.status_code != 200:
        res = response.json()
        error_message: str = ""
        if res['error']['type'] == 'invalid_request_error':
            error_message = res['error']['message']
        else:
            error_message = str(res)

        logger.fatal(error_message)
        return None

    logger.debug(f"OpenAI response: '{res}'", method=get_file_by_id)
    return res


def generate_image(key: str, prompt: str):

    req_url: str = "https://api.openai.com/v1/images/generations"

    headers: dict[str,str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}"
    }

    payload: dict = {
        "model": "dall-e-3",
        "prompt": prompt,
        "n": 1,
        "size": "1024x1024"
    }

    json_payload = json.dumps(payload)

    logger.debug(f"OpenAI with post payload for image generation: '{json_payload}'", method=generate_image)

    response: requests.Response = requests.post(url=req_url, headers=headers, data=json_payload)
    res = response.json()

    if response.status_code != 200:
        error_message: str = ""
        if res['error']['type'] == 'invalid_request_error' and res['error']['message'][0:25] == "":
            error_message = res['error']['message']
        else:
            error_message = str(res)

        logger.error(error_message)
        return None

    logger.debug(f"OpenAI response: '{res}'", method=generate_image)
    return res
