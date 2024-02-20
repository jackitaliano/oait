from argparse import Namespace

from utils.logger import logger
from services.assistants import create, modify, run


def add_service(parser):
    asst_parser = parser.add_parser('assts', help='Options for Assistants service. See `oait assts --help`')
    asst_parser.description = "OpenAI Tools for Assistants"
    asst_parser.usage="oait assts <create, mod, run> ..."

    asst_subparsers = asst_parser.add_subparsers(title='Available assistant services', dest='asst_service')
    create.add_service(asst_subparsers)
    modify.add_service(asst_subparsers)
    run.add_service(asst_subparsers)


def run_service(key: str, args: Namespace):

    logger.debug(f"Received arguments: {args}", method=run_service)

    service = args.asst_service

    if service == 'create':
        create.run_service(key, args)

    elif service == 'modify':
        modify.run_service(key, args)

    elif service == 'run':
        run.run_service(key, args)

    else:
        logger.fatal("Must choose assistant service (oait assts (create | modify | run)).", method=run_service) 
        exit(1)
