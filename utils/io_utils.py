import json
import os
import base64
from PIL import Image
from io import BytesIO

def read_txt(file_path: str):
    with open(file_path, 'r') as file:
        lines: list = [line.strip() for line in file.readlines()]
    
    return lines


def read_json(file_path: str):
    with open(file_path, 'r') as file:
        data: list = json.load(file)

    return data


def output_jsonl(file_path: str, data):
    with open(file_path, 'w') as file:
        for item in data:
            line = json.dumps(item)
            file.write(line + 'n')

def read_jsonl(file_path: str):
    with open(file_path, 'r') as file:
        for line in file:
            json_object = json.loads(line)

        print("line: " + json_object)


def display_image_from_bytes(bytes: str):
    decoded_content = base64.b64decode(bytes)
    image = Image.open(BytesIO(decoded_content))

    image.show()


def output_thread_to_file(file_path: str, data):
    file_type: str = os.path.splitext(file_path)[1][1:]

    if file_type == "json":
        output_to_json(file_path, data)

    elif file_type == "txt":
        pretty_print_thread_to_txt(file_path, data)

    else:
        print(f"ERROR (fatal): Must pass json or txt output file. Got: {file_type}")
        exit(1)


def output_to_json(file_path: str, data: dict | list):
    with open(file_path, 'w') as file:
        json.dump(data, file, indent=2)


def print_json(data):
    pretty = json.dumps(data, indent=2)

    print(pretty)


def pretty_print_thread_to_txt(file_path: str, data: list):
    with open(file_path, 'w') as file:
        pretty_print_thread(data, output=file.write)


def pretty_print_thread(data: list, output=None, deli="\n"):
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
