from time import sleep

from utils import cli_utils
from utils.logger import logger

from . import openai_utils


def add_service(parser):
    run_parser = parser.add_parser('run', help="Run an assistant. See `oait asst run --help`")

    run_parser.description = "OpenAI Assistant Run Tools. Start run, or view status."
    run_parser.usage = "oait asst run (<assistant_id> <thread_id> | -r RUN_ID)"

    run_parser.add_argument('asst', type=str, help="Assistant id to run on thread")
    run_parser.add_argument('thread', type=str, help="Thread id to have assistant run on")
    run_parser.add_argument('--run', '-r', type=str, help="Run id to check status of")


def run_service(key: str , args):
    asst_id: str = args.asst
    thread_id: str = args.thread
    run_id: str = args.run

    if not asst_id:
        logger.fatal(f"Must provide assitant id for run.", method=run_service)
        exit(1)
    if not thread_id:
        logger.fatal(f"Must provide thread id for run.", method=run_service)
        exit(1)

    run_id: str = run_assistant(key, asst_id, thread_id)

    wait_for_run(key, thread_id, run_id)


def run_assistant(key: str, asst_id: str, thread_id: str):
    logger.info(f"Running assistant ('{asst_id}') on thread ('{thread_id}')")
    
    run_info = openai_utils.run_assistant_on_thread(key, asst_id, thread_id)

    logger.info(f"Run started.")
    logger.debug(f"OpenAI Run info: '{run_info}'", method=run_assistant)

    run_id = run_info.get('id')

    return run_id


def wait_for_run(key: str, thread_id: str, run_id: str, delay_s: int = 1):
    logger.info(f"Waiting for assistant run ('{run_id}')")

    run_data = openai_utils.get_run_status(key, thread_id, run_id)

    while(run_data.get('status') != "completed"):
        cli_utils.display_progress_dots(1)
        sleep(delay_s / 3)
        cli_utils.display_progress_dots(2)
        sleep(delay_s / 3)
        run_data = openai_utils.get_run_status(key, thread_id, run_id)
        cli_utils.display_progress_dots(3)
        sleep(delay_s / 3)

    logger.info(f"Run complete for run '{run_id}'")
