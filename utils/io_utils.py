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
        line = file
        json_object = json.loads(line)

        print("line: " + json_object)


def display_image_from_bytes(file_name: str, image_bytes: str):
    
    image = Image.open(BytesIO(image_bytes))
    image.save(file_name)

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

    large_break: int = 70
    small_break: int = 20

    total_threads: int = len(data)

    for i, d in enumerate(data):

        thread_id: str = d['thread_id']
        thread: list = d['thread']

        thread_progress_str = f"[{i+1}/{total_threads}]"
        thread_id_str = f"{'Thread ID: ' + thread_id: <{large_break - len(thread_progress_str)}}"

        output(large_break * "-" + deli)
        output(thread_id_str + thread_progress_str + deli)
        output(large_break * "-" + deli)
        output(deli)

        total_messages: int = len(thread)
        for j, msg in enumerate(thread):
            role = msg['role']
            text = msg['text']

            progress_str = f"[{j+1}/{total_messages}]" 
            role_str = f"{role.upper(): <{small_break - len(progress_str)}}" 

            output(small_break * "-" + deli)
            output(role_str + progress_str + deli)
            output(small_break * "-" + deli)
            output(text + "\n" + deli)

        output(deli)
