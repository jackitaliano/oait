from argparse import Namespace

from utils import openai_utils
from utils import io_utils

def run_image_service(key: str, args: Namespace):

    image_file_id: str = args.retrieve

    if image_file_id:
        url: str = get_image_by_file_id(key, image_file_id)

    else:
        print("ERROR (fatal): Must choose method of image service") 
        exit(1)


def add_image_service(subparsers):
    images_parser = subparsers.add_parser('images', help='Options for images service. See `oait images --help`')

    images_parser.add_argument('--retrieve', '-r', type=str, help="Retrive image by file id")


def get_image_by_file_id(key: str, image_file_id: str) -> str:

    file_object: dict = openai_utils.get_file_by_id(key, image_file_id)

    if file_object is None:
        print(f"ERROR: No file by id '{image_file_id}'")

    else:
        io_utils.output_jsonl("test.jsonl", file_object)
        io_utils.read_jsonl("test.jsonl")
        # io_utils.display_image_from_bytes(data)

    return file_object

