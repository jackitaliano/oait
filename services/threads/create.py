from utils.logger import logger

from . import openai_utils


def add_service(parser):
    create_parser = parser.add_parser('create', help='Options for Threads Creation service. See `oait threads create --help`')

    create_parser.description = "OpenAI Thread Creation Tools."
    create_parser.usage = "oait threads create [-m MESSAGES]"

    create_parser.add_argument('--messages', '-m', type=str, help="Initial messages for thread", default=1)


def run_service(key: str, args):
    
    messages: str = args.messages

    handle_threads_create(key, messages)


def handle_threads_create(key: str, messages: str):



