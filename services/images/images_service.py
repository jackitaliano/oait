from argparse import Namespace

from utils.logger import logger
from services.images import gen, ret


def add_service(parser):
    image_parser = parser.add_parser('images', help='Options for Images service. See `oait images --help`')
    image_parser.description = "OpenAI Tools for Images"
    image_parser.usage="oait images <gen, ret> ..."

    image_subparsers = image_parser.add_subparsers(title='Available Images services', dest='image_service')
    gen.add_service(image_subparsers)
    ret.add_service(image_subparsers)


def run_service(key: str, args: Namespace):

    logger.debug(f"Received arguments: {args}", method=run_service)

    service = args.image_service

    if service == 'gen':
        gen.run_service(key, args)

    elif service == 'ret':
        ret.run_service(key, args)

    else:
        logger.fatal("Must choose image service (oait image (get | ret))", method=run_service) 
        exit(1)


