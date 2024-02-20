import requests
import json

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


def create_thread(key: str, initial_messages: str):
    logger.debug(f"OpenAI thread creation running.", method=create_thread)

    req_url: str = f"https://api.openai.com/v1/threads"

    headers: dict[str, str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
        "OpenAI-Beta": "assistants=v1"
    }

    payload: dict = initial_messages 

    json_payload = json.dumps(payload)

    logger.debug(f"OpenAI with post payload for creating thread '{json_payload}'", method=create_thread)

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

    logger.debug(f"OpenAI response: '{res}'", method=create_thread)
    return res


def add_message_to_thread(key: str, thread_id: str, message: dict):
    req_url: str = f"https://api.openai.com/v1/threads/{thread_id}/messages"

    headers: dict[str,str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
        "OpenAI-Beta": "assistants=v1"
    }

    payload: dict = message

    json_payload = json.dumps(payload)

    logger.debug(f"OpenAI with post payload for adding message to thread '{json_payload}'", method=add_message_to_thread)

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

    logger.debug(f"OpenAI response: '{res}'", method=add_message_to_thread)
    return res


def delete_thread(key: str, thread_id: str):
    logger.debug(f"Deleting OpenAI thread_id: '{thread_id}'")

    req_url: str = f"https://api.openai.com/v1/threads/{thread_id}"

    headers: dict[str, str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
        "OpenAI-Beta": "assistants=v1"
    }

    response: requests.Response = requests.delete(url=req_url, headers=headers)
    res = response.json()

    if response.status_code != 200:
        error_message: str = ""
        if res['error']['type'] == 'invalid_request_error' and res['error']['message'][0:25] == "":
            error_message = res['error']['message']
        else:
            error_message = str(res)

        logger.error(error_message)
        return None

    logger.debug(f"OpenAI response: '{res}'", method=delete_thread)
    return res
