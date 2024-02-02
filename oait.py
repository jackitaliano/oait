import argparse
import os
from services import threads_service

    
def run(args: argparse.Namespace):
    key: str | None = args.key

    if not key:
        key = os.environ.get("OPENAI_API_KEY")

        if not key:
            print('ERROR (fatal): Must provide OpenAI API Key (-k followed by key), or set environment variable (export OPENAI_API_KEY="your key")')
            exit(1)

    service = args.services
    
    if service == 'threads':
        threads_service.run_thread_service(key, args)

def main():
    parser: argparse.ArgumentParser = argparse.ArgumentParser(description="OpenAI Threads Retrieval. Read from cli args or file input. Write to stdout or file.")

    parser.add_argument('--key', '-k', type=str, help="Provide OpenAI key. Defaults to process.env.OPENAI_API_KEY")

    subparsers = parser.add_subparsers(title='Available services', dest='services')

    threads_service.add_thread_service(subparsers)

    subparsers.add_parser('test', help='Test')
    args = parser.parse_args()

    if not args.services:
        print("ERROR: Must enter a service (e.g. 'threads'). See oait --help")
        exit(1)

    run(args)


if __name__ == "__main__":
    main()


