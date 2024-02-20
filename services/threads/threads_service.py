from argparse import Namespace

from utils.logger import logger
from services.threads import ret


def add_service(parser):
    threads_parser = parser.add_parser('threads', help='Options for Threads service. See `oait thread --help`')
    threads_parser.description = "OpenAI Tools for Threads"
    threads_parser.usage="oait threads ret ..."

    threads_subparsers = threads_parser.add_subparsers(title='Available Threads services', dest='threads_service')
    ret.add_service(threads_subparsers)


def run_service(key: str, args: Namespace):

    logger.debug(f"Received arguments: {args}", method=run_service)

    service = args.threads_service

    if service == 'ret':
        ret.run_service(key, args)

    else:
        logger.fatal("Must choose threads service (`oait thread ret`).", method=run_service) 
        exit(1)
