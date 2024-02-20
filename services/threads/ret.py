from argparse import Namespace
import os

from utils import cli_utils, io_utils
from utils.logger import logger

from . import openai_utils

def run_service(key: str, args: Namespace):

    thread_ids: list[str] = args.thread_ids
    file: str = args.file
    output: str = args.output
    session_id:str = args.session
    limit:str = args.limit
    minlen: int = args.minlen

    if thread_ids and file:
        logger.fatal("Must pass thread_ids as a list of cli args OR a file input. Not both.")
        exit(1)

    parsed_threads: list[dict] = []
    if thread_ids:
        logger.info("Getting threads from list...")
        parsed_threads = get_threads_from_list(key, thread_ids, minlen)

    elif file:
        logger.info("Getting threads from file...")
        parsed_threads = get_threads_from_file(key, file, minlen)

    elif session_id:
        logger.info("Getting threads from session id...")
        parsed_threads = get_threads_from_session_id(key, session_id, minlen, limit)

    else:
        logger.fatal("Must pass list of thread_ids (space separated), an input file (json list or newline separated txt), or session.")
        exit(1)

    if output:
        logger.info(f"Outputing threads to file: {output}")
        io_utils.output_thread_to_file(output, parsed_threads)

    else:
        logger.info("Outputing threads to std...")
        io_utils.pretty_print_thread(parsed_threads)


def add_service(parser):
    threads_parser = parser.add_parser('ret', help='Options for Threads Retrieval service. See `oait thread ret --help`')

    threads_parser.description = "OpenAI Thread Retrieval Tools."
    threads_parser.usage = "oait thread ret ( <thread_id> [<thread_id> ...] | -f FILE | -s SESSION ) [-o OUTPUT] [-l LIMIT] [-ml MINLEN]"

    threads_parser.add_argument('thread_ids', nargs='*', type=str, help="Read messages from provided list of thread_ids (space separated)")
    threads_parser.add_argument('--file', '-f', nargs='?', type=str, help="Read thread_ids from file path. (json or newline separated txt)\n(Pass only -f for default: input.txt)", const="input.txt", default=None)
    threads_parser.add_argument('--session', '-s', type=str, help="Get thread messages from session id (Navigate to https://platform.openai.com/assistants. Copy Authorization header 'sess-')")
    threads_parser.add_argument('--output', '-o', nargs='?', type=str, help="Output file (json). (Pass only -o for default: output.json)", const="output.txt", default=None)
    threads_parser.add_argument('--limit', '-l', type=int, help="Limit for number of session threads")
    threads_parser.add_argument('--minlen', '-ml', type=int, help="Minimum length of threads to include in output.", default=1)


def parse_thread_data(thread_data: list) -> list[dict]:
    logger.debug(f"Parsing thread data: '{thread_data}'", method=parse_thread_data)
    messages: list[dict[str,str]] = []

    for data in thread_data:
        role: str = data['role']

        data_content: list[dict] = data['content']

        for content in data_content:
            text: str = ""

            if content.get('type') == 'text':
                text: str = content['text']['value']

            elif content.get('type') == 'image_file':
                file_id: str = content['image_file']['file_id']
                text: str = f"(Image file) file_id: {file_id}"

            elif len(content.get('text').get('annotations')) > 0:
                text: str = str(content['text']['annotations'])

            msg: dict[str, str] = { 'role': role, 'text': text }
            messages.append(msg)

    logger.debug(f"Parsed message data: {messages}", method=parse_thread_data)
    return messages


def parse_thread(thread_messages: dict, minlen) -> dict:
    logger.debug(f"Parsing thread messages. with minlen: '{minlen}'", method=parse_thread)
    thread_id: str = thread_messages['thread_id']
    try: 
        thread_data: list[dict] = thread_messages['data']

        parsed_thread: list[dict[str, str]] = parse_thread_data(thread_data)

        if len(parsed_thread) < minlen:
            return None
        # they're given in reverse order by openai
        parsed_thread = parsed_thread[::-1]

        parsed_thread_obj = {'thread_id': thread_id, 'thread': parsed_thread}
        logger.debug(f"Parsed thread: '{parsed_thread_obj}'", method=parse_thread)
        return parsed_thread_obj
    except Exception as e:
        logger.error(f"Parsing thread_id: {thread_id}. error: {e}")
        return {'thread_id': thread_id, 'thread': 'error'}


def parse_session_threads(session_threads: dict) -> list[str]:
    logger.debug(f"Parsing session_threads.", method=parse_session_threads)
    data: list[dict] = session_threads['data']
    thread_ids: list[str] = [thread['id'] for thread in data]

    logger.debug(f"Parsed ids: '{thread_ids}'.", method=parse_session_threads)
    return thread_ids


def get_threads_from_list(key: str, thread_ids: list[str], minlen: int):
    logger.debug(f"Getting threads from list.", method=get_threads_from_list)

    threads: list[dict] = []
    progress: int = 0
    total: int = len(thread_ids)

    cli_utils.update_progress(0, total)
    for thread_id in thread_ids:
        
        thread = openai_utils.get_thread_messages(key, thread_id, limit=100)

        if thread:
            threads.append(thread)

        progress += 1
        cli_utils.update_progress(progress, total)

    logger.info(f"Threads retrieved from list...", method=get_threads_from_list)
    logger.info(f"Parsing threads...", method=get_threads_from_list)

    parsed_threads: list[dict] = []
    for thread in threads:
        parsed_thread = parse_thread(thread, minlen)
        if parsed_thread:
            parsed_threads.append(parsed_thread)

    logger.debug(f"Parsed threads from list", method=get_threads_from_list)
    return parsed_threads


def get_threads_from_file(key: str, file: str, minlen: int):
    file_type: str = os.path.splitext(file)[1][1:]

    if file_type == "json":
        logger.info(f"Getting threads from json file: {file}.", method=get_threads_from_file)
        thread_ids: list[str] = io_utils.read_json(file)
        return get_threads_from_list(key, thread_ids, minlen)

    elif file_type == "txt":
        logger.info(f"Getting threads from txt file: {file}.", method=get_threads_from_file)
        thread_ids: list[str] = io_utils.read_txt(file)
        return get_threads_from_list(key, thread_ids, minlen)
    else:
        logger.fatal(f"Must pass json or txt input file. Got: {file_type}", method=get_threads_from_file)
        exit(1)


def get_threads_from_session_id(key: str, session_id: str, minlen: int, limit=None):
    if limit is None:
        limit = 10

    logger.info(f"Getting threads from session_id...")
    logger.debug(f"Session id: '{session_id}'", method=get_threads_from_session_id)
    sess_threads: dict = openai_utils.get_session_threads(session_id, limit)
    thread_ids: list[str] = parse_session_threads(sess_threads)

    if thread_ids:
        logger.info(f"Thread ids retrieved from session, getting threads from list.")
        return get_threads_from_list(key, thread_ids, minlen)
    
    else:
        logger.warning(f'No threads found for session_id="{session_id}"')
        exit(0)
