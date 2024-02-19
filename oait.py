import argparse
import os
from typing import TextIO
from utils.logger import logger, setup_logs
from services import threads_service, images_service, assitants_service
    

def run(args: argparse.Namespace):
    key: str | None = args.key

    if not key:
        key = os.environ.get("OPENAI_API_KEY")

        if not key:
            logger.fatal('Must provide OpenAI API Key (-k followed by key), or set environment variable (export OPENAI_API_KEY="your key")')
            exit(1)

    service = args.services
    logger.debug(f"Service selected: '{service}'", method=run)
    logger.debug(f"Args: '{args}'", method=run)
    
    if service == 'threads':
        threads_service.run_thread_service(key, args)

    elif service == 'images':
        images_service.run_image_service(key, args)

    elif service == 'asst':
        assitants_service.run_assistants_service(key, args)


def main():

    parser: argparse.ArgumentParser = argparse.ArgumentParser(prog="oait")
    parser.description="""OpenAI Tools.\n
    Tools for interacting with the OpenAI API via:\n
      - cli (`oait ...`)\n
      - python cli (`python oait.py ...`)\n
      - copy/pasting source into your own project\n
    """
    parser.usage=f"{parser.prog} [-v/-d/-s/-S/-l] [-k] <threads, images, asst> ..."

    parser.add_argument('--key', '-k', type=str, help="Provide OpenAI key. Defaults to process.env.OPENAI_API_KEY")
    parser.add_argument('--verbose', '-v', action="store_true", help="Enable verbosity of logging (info level)")
    parser.add_argument('--debug', '-d', action="store_true", help="Set logging level to debug (extremely verbose).")
    parser.add_argument('--silent', '-s', action="store_true", help="Silence all logs (other than FATAL).")
    parser.add_argument('--Silent', '-S', action="store_true", help="Silence ALL logs.")
    parser.add_argument('--logs', '-l', action="store_true", help="Save logs to a log file.")

    subparsers = parser.add_subparsers(title='Available services', dest='services')

    threads_service.add_thread_service(subparsers)
    images_service.add_image_service(subparsers)
    assitants_service.add_assistant_service(subparsers)

    args = parser.parse_args()

    logs_file: TextIO | None = setup_logs(args)

    try:
        if not args.services:
            logger.fatal("Must enter a service (e.g. 'threads'). See oait --help")

        run(args)

    except Exception as e:
        logger.fatal(str(e), method=main)
        raise e
        
    finally:
        if logs_file and not logs_file.closed:
            logs_file.close()


if __name__ == "__main__":
    main()


