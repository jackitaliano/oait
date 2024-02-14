import argparse
import os
from typing import TextIO
from utils.logger import logger
from services import threads_service, images_service
    

def run(args: argparse.Namespace):
    key: str | None = args.key

    if not key:
        key = os.environ.get("OPENAI_API_KEY")

        if not key:
            logger.fatal('Must provide OpenAI API Key (-k followed by key), or set environment variable (export OPENAI_API_KEY="your key")')
            exit(1)

    service = args.services
    
    if service == 'threads':
        threads_service.run_thread_service(key, args)

    elif service == 'images':
        images_service.run_image_service(key, args)


def main():

    parser: argparse.ArgumentParser = argparse.ArgumentParser(description="OpenAI Threads Retrieval. Read from cli args or file input. Write to stdout or file.")

    parser.add_argument('--key', '-k', type=str, help="Provide OpenAI key. Defaults to process.env.OPENAI_API_KEY")
    parser.add_argument('--verbose', '-v', action="store_true", help="Enable verbosity of logging (info level)")
    parser.add_argument('--debug', '-d', action="store_true", help="Set logging level to debug (extremely verbose).")
    parser.add_argument('--silent', '-s', action="store_true", help="Silence all logs (other than FATAL).")
    parser.add_argument('--Silent', '-S', action="store_true", help="Silence ALL logs.")
    parser.add_argument('--logs', '-l', action="store_true", help="Save logs to a log file.")

    subparsers = parser.add_subparsers(title='Available services', dest='services')

    threads_service.add_thread_service(subparsers)
    images_service.add_image_service(subparsers)

    args = parser.parse_args()

    logs_file: TextIO | None = logger.setup_logs(args)

    try:
        if not args.services:
            logger.fatal("Must enter a service (e.g. 'threads'). See oait --help")

        run(args)

    except Exception as e:
        logger.fatal(str(e), method=main)
        
    finally:
        if logs_file and not logs_file.closed:
            logs_file.close()


if __name__ == "__main__":
    main()


