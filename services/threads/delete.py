import os

from utils import cli_utils, io_utils
from utils.logger import logger
from . import openai_utils


def add_service(parser):
    delete_parser = parser.add_parser('del', help='Options for Threads Deletion service. See `oait threads del --help`')

    delete_parser.description = "OpenAI Thread Deletion Tools."
    delete_parser.usage = "oait threads del ( <thread_id> [<thread_id> ...] | -f FILE | -s SESSION ) [-l LIMIT] [-Ml MAXLEN]"

    delete_parser.add_argument('thread_ids', nargs='*', type=str, help="Read messages from provided list of thread_ids (space separated)")
    delete_parser.add_argument('--file', '-f', nargs='?', type=str, help="Read thread_ids from file path. (json or newline separated txt)\n(Pass only -f for default: input.txt)", const="input.txt", default=None)
    delete_parser.add_argument('--session', '-s', type=str, help="Get thread messages from session id and delete (Navigate to https://platform.openai.com/assistants. Copy Authorization header 'sess-')")
    delete_parser.add_argument('--limit', '-l', type=int, help="Limit for number of session threads include in deletion.(Default=1)", default=1)
    delete_parser.add_argument('--Maxlen', '-Ml', type=int, help="Maximum length of threads to delete. (Default=1)", default=1)


def run_service(key: str, args):

    thread_ids: list[str] = args.thread_ids
    file: str = args.file
    session: str = args.session
    limit: int = args.limit
    maxlen: int = args.Maxlen

    if thread_ids:
        logger.info("Getting threads from list...")
        del_threads_from_list(key, thread_ids)

    elif file:
        del_threads_from_file(key, file, maxlen)

    elif session:
        del_threads_from_session_id(key, session, maxlen=maxlen, limit=limit)

    else:
        logger.fatal("Must pass list of thread_ids (space separated), an input file (json list or newline separated txt), or session.")
        exit(1)


def del_threads_from_list(key: str, thread_ids: list[str]):
    logger.debug(f"Getting threads from list.", method=del_threads_from_list)

    progress: int = 0
    total: int = len(thread_ids)

    if total <= 0:
        logger.info("No threads to delete")
        return

    if not confirm_delete_threads(total):
        return

    logger.info(f"Deleting threads...", method=del_threads_from_list)
    print("\rDeleting threads...")

    cli_utils.update_progress(0, total)
    for thread_id in thread_ids:
        
        res = openai_utils.delete_thread(key, thread_id)

        progress += 1
        cli_utils.update_progress(progress, total)

    logger.info(f"Threads deleted from list...", method=del_threads_from_list)


def del_threads_from_file(key: str, file: str, maxlen: int = None):
    file_type: str = os.path.splitext(file)[1][1:]

    if file_type == "json":
        logger.info(f"Getting threads from json file: {file}.", method=del_threads_from_file)
        threads: list[dict] = io_utils.read_json(file)
        thread_ids = [thread.get('thread_id') for thread in threads]

        if maxlen:
            thread_ids = filter_threads_from_list(key, thread_ids, maxlen)

        del_threads_from_list(key, thread_ids)

    elif file_type == "txt":
        logger.info(f"Getting threads from txt file: {file}.", method=del_threads_from_file)
        thread_ids: list[str] = io_utils.read_txt(file)

        if maxlen:
            thread_ids = filter_threads_from_list(key, thread_ids, maxlen)

        del_threads_from_list(key, thread_ids)

    else:
        logger.fatal(f"Must pass json or txt input file. Got: {file_type}", method=del_threads_from_file)
        exit(1)


def del_threads_from_session_id(key: str, session_id: str, maxlen: int = None, limit=None):
    if limit is None:
        limit = 1

    logger.info(f"Getting threads from session_id...")
    logger.debug(f"Session id: '{session_id}'", method=del_threads_from_session_id)
    sess_threads: dict = openai_utils.get_session_threads(session_id, limit)
    thread_ids: list[str] = parse_session_threads(sess_threads)

    if thread_ids:
        logger.info(f"Thread ids retrieved from session, getting threads from list.")

        if maxlen:
            thread_ids = filter_threads_from_list(key, thread_ids, maxlen)

        del_threads_from_list(key, thread_ids)
        # return del_threads_from_list(key, thread_ids, maxlen)
    
    else:
        logger.warning(f'No threads found for session_id="{session_id}"')
        exit(0)


def parse_session_threads(session_threads: dict) -> list[str]:
    logger.debug(f"Parsing session_threads.", method=parse_session_threads)
    data: list[dict] = session_threads['data']
    thread_ids: list[str] = [thread['id'] for thread in data]

    logger.debug(f"Parsed ids: '{thread_ids}'.", method=parse_session_threads)
    return thread_ids


def filter_threads_from_list(key: str, thread_ids: list[str], maxlen: int):
    logger.debug(f"Getting threads from list.", method=filter_threads_from_list)

    filtered_thread_ids: list = []
    progress: int = 0
    total: int = len(thread_ids)

    logger.info("Getting threads for filtering...")
    print("Filtering threads for deletion...")
    cli_utils.update_progress(0, total)
    for thread_id in thread_ids:
        
        thread = openai_utils.get_thread_messages(key, thread_id, limit=100)

        if thread:
            thread_data = thread.get('data')
            if not thread_data is None and len(thread_data) <= maxlen:
                filtered_thread_ids.append(thread_id)

        progress += 1
        cli_utils.update_progress(progress, total)

    logger.info("Threads filtered.")

    return filtered_thread_ids


def confirm_delete_threads(num_threads: int):
    print(f"Confirm deleting {num_threads} [y/n]: ", end="")
    confirm = input()

    if confirm.lower() in ('y', 'yes'):
        print("\r", end="\r")
        return True

    return False
