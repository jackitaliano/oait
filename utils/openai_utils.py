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


def create_assistant(key: str, assistant_config: dict):
    req_url: str = "https://api.openai.com/v1/assistants"

    headers: dict[str,str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
        "OpenAI-Beta": "assistants=v1"
    }

    payload: dict = assistant_config

    json_payload = json.dumps(payload)

    logger.debug(f"OpenAI with post payload for assistant creation '{json_payload}'", method=create_assistant)

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

    logger.debug(f"OpenAI response: '{res}'", method=create_assistant)
    return res


def modify_assistant(key: str, asst_id: str, assistant_config: dict):
    req_url: str = f"https://api.openai.com/v1/assistants/{asst_id}"

    headers: dict[str,str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
        "OpenAI-Beta": "assistants=v1"
    }

    payload: dict = assistant_config

    json_payload = json.dumps(payload)

    logger.debug(f"OpenAI with post payload for assistant modification'{json_payload}'", method=modify_assistant)

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

    logger.debug(f"OpenAI response: '{res}'", method=modify_assistant)
    return res





def run_assistant_on_thread(key: str, assistant_id: str, thread_id: str):
    req_url: str = f"https://api.openai.com/v1/threads/{thread_id}/runs"

    headers: dict[str,str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
        "OpenAI-Beta": "assistants=v1"
    }

    payload: dict = { 'assistant_id': assistant_id }

    json_payload = json.dumps(payload)

    logger.debug(f"OpenAI with post payload for adding message to thread_id: '{thread_id}' and paylaod: '{json_payload}'", method=run_assistant_on_thread)

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

    logger.debug(f"OpenAI response: '{res}'", method=run_assistant_on_thread)
    return res


def get_run_status(key: str, thread_id: str, run_id: str):
    logger.debug(f"OpenAI checking run status of thread_id: '{thread_id}', run_id: '{run_id}'", method=get_run_status)

    req_url: str = f"https://api.openai.com/v1/threads/{thread_id}/runs/{run_id}"

    headers: dict[str,str] = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {key}",
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

    logger.debug(f"OpenAI response: '{res}'", method=get_run_status)
    return res


