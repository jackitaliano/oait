import requests
import json

from utils.logger import logger

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
