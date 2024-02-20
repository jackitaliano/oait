import requests
from utils.logger import logger


def get_thread_messages(key: str, thread_id: str, limit: int) -> dict | None:
    logger.debug(f"Getting OpenAI thread messages with thread_id: '{thread_id}' and limit: '{limit}'")

    req_url: str = f"https://api.openai.com/v1/threads/{thread_id}/messages?limit={limit}"

    headers: dict[str, str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
        "OpenAI-Beta": "assistants=v1"
    }

    response: requests.Response = requests.get(url=req_url, headers=headers)
    res = response.json()

    if response.status_code != 200:
        error_message: str = ""
        if res['error']['type'] == 'invalid_request_error' and res['error']['message'][0:25] == "No thread found with id '":
            error_message = res['error']['message']
        else:
            error_message = str(res)

        logger.error("ERROR: " + error_message)
        return None

    res['thread_id'] = thread_id

    logger.debug(f"OpenAI Response: '{res}'", method=get_thread_messages)
    return res


def get_session_threads(session_id: str, limit: int):
    logger.debug(f"Getting OpenAI session threads with limit: '{limit}'")

    req_url: str = f"https://api.openai.com/v1/threads?limit={limit}"

    headers: dict[str, str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {session_id}",
        "OpenAI-Beta": "assistants=v1"
    }

    response: requests.Response = requests.get(url=req_url, headers=headers)
    res = response.json()

    if response.status_code != 200:
        error_message: str = ""
        if res['error']['type'] == 'invalid_request_error' and res['error']['message'][0:25] == "":
            error_message = res['error']['message']
        else:
            error_message = str(res)

        logger.error(error_message)
        return None

    logger.debug(f"OpenAI response: '{res}'", method=get_session_threads)
    return res


