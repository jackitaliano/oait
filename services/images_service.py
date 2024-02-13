from argparse import Namespace

from utils import openai_utils, io_utils, request_utils, cli_utils

def run_image_service(key: str, args: Namespace):

    image_file_id: str = args.file_id
    output: str = args.output
    prompt: str = args.prompt
    image_url: str = args.url

    print()
    print(args)
    print()

    if image_file_id:
        get_image_by_file_id(key, image_file_id, output)

    elif image_url:
        get_image_by_url(image_url, output)

    elif prompt:
         generate_image(key, prompt, output)

    else:
        print("ERROR (fatal): Must choose method of image service") 
        exit(1)


def add_image_service(subparsers):
    images_parser = subparsers.add_parser('images', help='Options for images service. See `oait images --help`')

    images_parser.add_argument('--file_id', '-f', type=str, help="Retrive image by file id.")
    images_parser.add_argument('--prompt', '-p', type=str, help="Generate image by prompt.")
    images_parser.add_argument('--url', '-u', type=process_url, help="Retrieve image by url.")
    images_parser.add_argument('--output', '-o', nargs='?', type=str, help="Output file ('png' or 'json'). 'png' for retrieving image by file id or url. Use either for generating image. (Pass only -o for default: [file-id].png / [image.json])", const="default", default=None)


def process_url(arg):
    return arg.replace('\\', '')


def get_image_by_file_id(key: str, image_file_id: str, image_file_name: str) -> str:

    file_object: dict = openai_utils.get_file_by_id(key, image_file_id)

    if file_object is None:
        exit(1)

    else:

        if image_file_name is None:
            image_file_name = image_file_id + ".png"
        elif image_file_name[-4:] != ".png":
            image_file_name = image_file_id + ".png"
            print(f"ERROR: File name not ending in '.png'. Using default ({image_file_name})")

        io_utils.display_image_from_bytes(image_file_name, file_object)


def get_image_by_url(image_url: str, image_file_name: str) -> None:

    image_data = request_utils.get_image_from_url(image_url)

    if image_file_name is None:
        image_file_name = "image.png"

    
    file_type: str = cli_utils.get_file_type(image_file_name)
    if file_type != "png":
        image_file_name = "image.png"
        print(f"ERROR: Only supporting output to '.png' files. Got: {file_type}. Defaulting to: {image_file_name}")

    if image_data:
        io_utils.output_image_to_file(image_file_name, image_data)


def generate_image(key: str, prompt: str, image_file_name: str):
    image_res = openai_utils.generate_image(key, prompt)

    data = image_res.get('data')
    
    if not data:
        print(f"ERROR: no data return from api.")
        return

    images = data

    if image_file_name is None:
        image_file_name = "image.json"

    if image_res:
        print(f"Image generated. Saving to file {image_file_name} ...")

    file_type: str = cli_utils.get_file_type(image_file_name)
    if file_type != "png" and file_type != "json":
        image_file_name = "image.json"
        (f"ERROR: Only supporting output to '.png' and '.json' files. Got: {file_type}. Defaulting to: {image_file_name}")

    if file_type == "png":
        urls: list[str] = [img.get('url') for img in images]

        image_data = [request_utils.get_image_from_url(url) for url in urls ]

        if image_data:
            io_utils.output_image_to_file("test.png", image_data[0])

    elif file_type == "json":
        io_utils.output_to_json(image_file_name, images)

