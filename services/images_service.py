from argparse import Namespace

from utils import openai_utils
from utils import io_utils

def run_image_service(key: str, args: Namespace):

    image_file_id: str = args.retrieve
    image_file_name: str = args.output

    if image_file_id:
        url: str = get_image_by_file_id(key, image_file_id, image_file_name)

    else:
        print("ERROR (fatal): Must choose method of image service") 
        exit(1)


def add_image_service(subparsers):
    images_parser = subparsers.add_parser('images', help='Options for images service. See `oait images --help`')

    images_parser.add_argument('--retrieve', '-r', type=str, help="Retrive image by file id")
    images_parser.add_argument('--output', '-o', nargs='?', type=str, help="Output file (png). (Pass only -o for default: [file-id].png)", const="default", default=None)


def get_image_by_file_id(key: str, image_file_id: str, image_file_name: str) -> str:

    file_object: dict = openai_utils.get_file_by_id(key, image_file_id)

    if file_object is None:
        print(f"ERROR: No file by id '{image_file_id}'")
        exit(1)

    else:

        if image_file_name is None:
            image_file_name = image_file_id + ".png"
        elif image_file_name[-3:] != ".png":
            image_file_name = image_file_id + ".png"
            print(f"ERROR: File name not ending in '.png'. Using default ({image_file_name})")

        io_utils.display_image_from_bytes(image_file_name, file_object)

        return image_file_name

