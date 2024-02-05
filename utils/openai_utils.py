import requests

def get_file_by_id(key: str, file_id: str) -> dict | None:
    req_url: str = f"https://api.openai.com/v1/files/{file_id}/content"

    headers: dict[str, str] = {
        "Authorization": f"Bearer {key}",
    }

    response: requests.Response = requests.get(url=req_url, headers=headers)
    res = response.content

    if response.status_code != 200:
        error_message: str = ""
        if res['error']['type'] == 'invalid_request_error':
            error_message = res['error']['message']
        else:
            error_message = str(res)

        print("ERROR: " + error_message)
        return None

    return res

def get_thread_messages(key: str, thread_id: str) -> dict | None:
    req_url: str = f"https://api.openai.com/v1/threads/{thread_id}/messages"

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

        print("ERROR: " + error_message)
        return None

    res['thread_id'] = thread_id

    return res


def get_session_threads(session_id: str, limit: int):
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

        print("ERROR: " + error_message)
        return None

    return res
