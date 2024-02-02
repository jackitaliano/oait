import argparse
import json
import os
from re import error

import requests


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



def parse_thread_data(thread_data: list) -> list[dict]:
    messages: list[dict[str,str]] = []

    for data in thread_data:
        role: str = data['role']

        data_content: list[dict] = data['content']

        for content in data_content:
            if content['type'] == 'text':
                text: str = content['text']['value']
                msg: dict[str, str] = { 'role': role, 'text': text }
                messages.append(msg)

            elif len(content['text']['annotations']) > 0:
                text: str = str(content['text']['annotations'])
                msg: dict[str, str] = { 'role': role, 'text': text }
                messages.append(msg)

    return messages


def parse_thread(thread_messages: dict) -> dict:
    thread_data: list[dict] = thread_messages['data']
    thread_id: str = thread_messages['thread_id']

    parsed_thread: list[dict[str, str]] = parse_thread_data(thread_data)

    # they're given in reverse order by openai
    parsed_thread = parsed_thread[::-1]

    return {'thread_id': thread_id, 'thread': parsed_thread}


def parse_session_threads(session_threads: dict) -> list[str]:
    data: list[dict] = session_threads['data']
    thread_ids: list[str] = [thread['id'] for thread in data]

    return thread_ids


def get_threads_from_list(key: str, thread_ids: list[str]):

    threads: list[dict] = []
    progress: int = 0
    total: int = len(thread_ids)

    for thread_id in thread_ids:
        update_progress(progress, total)
        
        thread = get_thread_messages(key, thread_id)

        if thread:
            threads.append(thread)

        progress += 1

    parsed_threads: list[dict] = [parse_thread(thread) for thread in threads]

    return parsed_threads


def get_threads_from_file(key: str, file: str):
    file_type: str = os.path.splitext(file)[1][1:]

    if file_type == "json":
        thread_ids: list[str] = read_json(file)
        return get_threads_from_list(key, thread_ids)

    elif file_type == "txt":
        thread_ids: list[str] = read_txt(file)
        return get_threads_from_list(key, thread_ids)
    else:
        print(f"ERROR (fatal): Must pass json or txt input file. Got: {file_type}")
        exit(1)


def get_threads_from_session_id(key: str, session_id: str, limit=None):
    if limit is None:
        limit = 50

    sess_threads: dict = get_session_threads(session_id, limit)
    thread_ids: list[str] = parse_session_threads(sess_threads)

    if thread_ids:
        return get_threads_from_list(key, thread_ids)
    
    else:
        print(f'No threads found for session_id="{session_id}"')
        exit(0)


def update_progress(progress, total):
    print(f"\rProgress:[{progress}/{total}]", end="") 

def read_txt(file_path: str):
    with open(file_path, 'r') as file:
        lines: list = [line.strip() for line in file.readlines()]
    
    return lines


def read_json(file_path: str):
    with open(file_path, 'r') as file:
        data: list = json.load(file)

    return data


def output_to_file(file_path: str, data):
    file_type: str = os.path.splitext(file_path)[1][1:]

    if file_type == "json":
        output_to_json(file_path, data)

    elif file_type == "txt":
        pretty_print_to_txt(file_path, data)

    else:
        print(f"ERROR (fatal): Must pass json or txt output file. Got: {file_type}")
        exit(1)


def output_to_json(file_path: str, data: dict | list):
    with open(file_path, 'w') as file:
        json.dump(data, file, indent=2)


def print_json(data):
    pretty = json.dumps(data, indent=2)

    print(pretty)


def pretty_print_to_txt(file_path: str, data: list):
    with open(file_path, 'w') as file:
        pretty_print(data, output=file.write)


def pretty_print(data: list, output=None, deli="\n"):
    if output is None:
        output=print
        deli=""

    large_break = 50
    small_break = 20

    for d in data:

        thread_id: str = d['thread_id']
        thread: list = d['thread']

        output(large_break * "-" + deli)
        output(f"Thread ID: {thread_id}" + deli)
        output(large_break * "-" + deli)
        output(deli)

        for msg in thread:
            role = msg['role']
            text = msg['text']

            output(small_break * "-" + deli)
            output(f"{role.upper():^{small_break}}" + deli)
            output(small_break * "-" + deli)
            output(text + "\n" + deli)

        output(deli)

    
def run(args: argparse.Namespace):
    key: str | None = args.key

    if not key:
        key = os.environ.get("OPENAI_API_KEY")

        if not key:
            print('ERROR (fatal): Must provide OpenAI API Key (-k followed by key), or set environment variable (export OPENAI_API_KEY="your key")')
            exit(1)

    thread_ids: list[str] = args.thread_ids
    file: str = args.file
    output: str = args.output
    session_id:str = args.session
    limit:str = args.limit

    if thread_ids and file:
        print("ERROR (fatal): Must pass thread_ids as a list of cli args OR a file input. Not both.")
        exit(1)

    parsed_threads: list[dict] = []
    if thread_ids:
        parsed_threads = get_threads_from_list(key, thread_ids)

    elif file:
        parsed_threads = get_threads_from_file(key, file)

    elif session_id:
        parsed_threads = get_threads_from_session_id(key, session_id, limit)

    else:
        print("ERROR (fatal): Must pass list of thread_ids (space separated), an input file (json list or newline separated txt), or session.")
        exit(1)

    if output:
        output_to_file(output, parsed_threads)

    else:
        pretty_print(parsed_threads)


def main():
    parser: argparse.ArgumentParser = argparse.ArgumentParser(description="OpenAI Threads Retrieval. Read from cli args or file input. Write to stdout or file.")

    parser.add_argument('thread_ids', nargs='*', type=str, help="Read messages from provided list of thread_ids (space separated)")
    parser.add_argument('--key', '-k', type=str, help="Provide OpenAI key. Defaults to process.env.OPENAI_API_KEY")
    parser.add_argument('--file', '-f', nargs='?', type=str, help="Read thread_ids from file path. (json or newline separated txt)\n(Pass only -f for default: input.txt)", const="input.txt", default=None)
    parser.add_argument('--output', '-o', nargs='?', type=str, help="Output file (json). (Pass only -o for default: output.json)", const="output.txt", default=None)
    parser.add_argument('--session', '-s', type=str, help="Get thread messages from session id (Navigate to https://platform.openai.com/assistants. Copy Authorization header 'sess-')")
    parser.add_argument('--limit', '-l', type=int, help="Limit for number of session threads")

    args = parser.parse_args()
    
    run(args)


if __name__ == "__main__":
    main()


