from argparse import Namespace
from time import sleep
from services.images_service import cli_utils

from utils import openai_utils
from utils.logger import logger


def run_assistants_service(key: str, args: Namespace):

    logger.debug(f"Received arguments: {args}", method=run_assistants_service)

    service = args.asst_service

    if service == 'create':
        handle_assistant_create(key, args)

    elif service == 'mod':
        handle_assistant_modify(key, args)

    elif service == 'run':
        handle_run_assistant_on_thread(key, args)

    else:
        logger.fatal("Must choose assistant service (oait assts ['create'/'run']).", method=run_assistants_service) 
        exit(1)


def add_assistant_service(parser):
    asst_parser = parser.add_parser('asst', help='Options for Assistants service. See `oait asst --help`')
    asst_parser.description = "OpenAI Tools for Assistants"
    asst_parser.usage="oait asst <create, mod, run> ..."

    asst_subparsers = asst_parser.add_subparsers(title='Available assistant services', dest='asst_service')
    add_create_service(asst_subparsers)
    add_mod_service(asst_subparsers)
    add_run_service(asst_subparsers)


def add_create_service(parser):
    create_parser = parser.add_parser('create', help="Create an assistant. See `oait asst create --help`")
    create_parser.add_argument('--name', '-n', type=str, help="Name of created assistant.", default=None)
    create_parser.add_argument('--instructions', '-i', type=str, help="Instructions for created assistant", default=None)
    create_parser.add_argument('--description', '-d', type=str, help="Description of created assistant", default=None)
    create_parser.add_argument('--model', '-m', type=str, help="Model used by created assistant (Default 'gpt-3.5-turbo-0125')", default="gpt-3.5-turbo-0125")
    create_parser.add_argument('--tools', '-t', type=str, help="Tools to provide to created assistant (code, retrieval, functions)", default=None)


def add_mod_service(parser):
    mod_parser = parser.add_parser('mod', help="Modify an assistant. See `oait asst mod --help`")
    mod_parser.add_argument('--asst', '-a', type=str, help="Assistant id to modify.")
    mod_parser.add_argument('--name', '-n', type=str, help="Name of updated assistant.", default=None)
    mod_parser.add_argument('--instructions', '-i', type=str, help="Instructions for updated assistant ", default=None)
    mod_parser.add_argument('--description', '-d', type=str, help="Description of updated assistant")
    mod_parser.add_argument('--tools', '-t', type=str, help="Tools for updated assistant (code, retrieval, functions)")   
    mod_parser.add_argument('--model', '-m', type=str, help="Model to use for updated assistant", default=None)


def add_run_service(parser):
    run_parser = parser.add_parser('run', help="Run an assistant. See `oait asst run --help`")

    run_parser.description = "OpenAI Assistant Run Tools. Start run, or view status."
    run_parser.usage = "oait asst run <-a 'assistant_id' -t 'thread_id' / -r 'run_id'>"

    run_parser.add_argument('--asst', '-a', type=str, help="Assistant id to run on thread")
    run_parser.add_argument('--thread', '-t', type=str, help="Thread id to have assistant run on")
    run_parser.add_argument('--run', '-r', type=str, help="Run id to check status of")


def handle_assistant_create(key, args):
    name: str = args.name
    instructions: str = args.instructions
    description: str = args.description
    model: str =  args.model
    tools: str = args.tools

    if not name:
        logger.fatal(f"Assistant name must be passed for creation.", method=handle_assistant_create)
        exit(1)
    if not model:
        logger.fatal(f"Assistant model must be passed for creation.", method=handle_assistant_create)
        exit(1)
    if not instructions:
        logger.fatal(f"Assistant instructions must be passed for creation.", method=handle_assistant_create)
        exit(1)

    assistant_config: dict = { 'name': name, 'model': model, 'instructions': instructions }

    if description:
        assistant_config['description'] = description
    if tools:
        assistant_config['tools'] = tools

    logger.info(f"Creating assistant with config: '{assistant_config}'...")

    asst_id: str = create_assistant_from_config(key, assistant_config)

    logger.info(f"Created assistant (assistant_id: '{asst_id}').")


def handle_assistant_modify(key, args):
    asst: str = args.asst
    name: str = args.name
    instructions: str = args.instructions
    description: str = args.description
    model: str =  args.model
    tools: str = args.tools

    if not asst:
        logger.fatal("Must provide assistant id to modify.", method=handle_assistant_modify)
        exit(1)

    assistant_config = {}
    if name:
        assistant_config['name'] = name
    if instructions:
        assistant_config['instructions'] = instructions
    if description:
        assistant_config['description'] = description
    if model:
        assistant_config['model'] = model
    if tools:
        assistant_config['tools'] = tools

    if len(assistant_config) <= 0:
        logger.fatal("Must pass a field to update assistant. See `oait asst mod -h`")
        exit(0)

    logger.info(f"Creating assistant with config: '{assistant_config}'...")

    asst_id: str = modify_assistant_from_config(key, asst, assistant_config)

    logger.info(f"Created assistant (assistant_id: '{asst_id}').")


def handle_run_assistant_on_thread(key: str , args):
    asst_id: str = args.asst
    thread_id: str = args.thread
    run_id: str = args.run

    if not asst_id:
        logger.fatal(f"Must provide assitant id for run.", method=handle_run_assistant_on_thread)
        exit(1)
    if not thread_id:
        logger.fatal(f"Must provide thread id for run.", method=handle_run_assistant_on_thread)
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

    
def create_assistant_from_config(key: str, assistant_config: dict):
    logger.info("Creating assistant...")

    response = openai_utils.create_assistant(key, assistant_config)

    if not response:
        logger.fatal("No data returned from OpenAI.", method=create_assistant_from_config)

    asst_id: str | None = response.get('id')

    if not asst_id:
        logger.fatal("No assistant id returned from OpenAI", method=create_assistant_from_config)

    logger.info(f"OpenAI assistant id created: '{asst_id}'")

    return asst_id


def modify_assistant_from_config(key: str, asst_id: str, assistant_config: dict):
    logger.info("Updating assistant...")

    response = openai_utils.modify_assistant(key, asst_id, assistant_config)

    if not response:
        logger.fatal("No data returned from OpenAI.", method=modify_assistant_from_config)

    res_asst_id: str | None = response.get('id')

    if not res_asst_id:
        logger.fatal("No assistant id returned from OpenAI", method=modify_assistant_from_config)

    logger.info(f"OpenAI assistant id updated: '{res_asst_id}'")

    return res_asst_id


