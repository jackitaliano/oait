from utils import openai_utils, cli_utils, io_utils, request_utils
from utils.logger import logger, Timestamp


def add_service(parser):
    ret_parser = parser.add_parser('ret', help="Retrieve an image based on file id or from URL. See `oait image ret --help`")

    ret_parser.description = "OpenAI Image Retrieval Tools" 
    ret_parser.usage = "oait image ret [-f FILE] [-u URL] [-o OUTPUT]"

    ret_parser.add_argument('--file', '-f', type=str, help="Retrive image by file id.")
    ret_parser.add_argument('--url', '-u', type=process_url, help="Retrieve image by url.")
    ret_parser.add_argument('--output', '-o', nargs='?', type=str, help="Output file ('png' or 'jsonl'). Pass only -o for default: [file-id].png", const="default", default=None)


def process_url(arg):
    return arg.replace('\\', '')


def run_service(key, args):
    image_url: str = args.url
    image_file_id: str = args.fileid
    output: str = args.output

    image_data = None

    if image_file_id:
        logger.info(f"Retrieving image by file_id ('{image_file_id}')...")
        image_data = get_image_by_file_id(key, image_file_id)

    elif image_url:
        logger.info(f"Retrieving image by url ('{image_url}')...")
        image_data = get_image_by_url(image_url)

    else:
        logger.fatal("Image retrieval method unknown. Must supply file-id or url.", method=run_service)
        exit(1)


    if image_data is None:
        logger.fatal("Image data not found.")
        exit(1)

    logger.debug("Outputing image data.", method=run_service)
    handle_image_output(image_data, output)


def get_image_by_file_id(key: str, image_file_id: str):
    logger.info("Fetching image from OpenAI by file_id...")

    file_object: dict = openai_utils.get_file_by_id(key, image_file_id)

    if file_object is None:
        logger.fatal("No image returned by file id.", method=get_image_by_file_id)
        exit(1)

    logger.info("File object returned from api by file id.", method=get_image_by_file_id)
    return file_object


def get_image_by_url(image_url: str):
    logger.info("Fetching image from url...")

    image_data = request_utils.get_image_from_url(image_url)

    if not image_data:
        logger.fatal("No data return from url.", method=get_image_by_url)
        exit(1)

    logger.info("Image data returned from url.", method=get_image_by_url)
    return image_data


def handle_image_output(image_data, output_fp):
    image_type = get_instance_type_name(image_data)

    if image_type == "dict":
        logger.debug(f"Outputing image of type dict.", method=handle_image_output)
        handle_image_dict_output(image_data, output_fp)
        return

    elif image_type == "bytes":
        logger.debug(f"Outputing image of type bytes.", method=handle_image_output)
        handle_image_bytes_output(image_data, output_fp)
        return

    elif image_type == "list":
        logger.debug(f"Outputing image of type list.", method=handle_image_output)
        for item in image_data:
            handle_image_output(item, output_fp)

        return

    else:
        logger.fatal(f"Image type '{image_type}' not supported.", method=handle_image_output)
        exit(1)


def handle_image_dict_output(image_dict, output_fp):
    if output_fp is None:
        date_time = logger._format_date_time(Timestamp.DATE_AND_TIME)
        output_fp = f"image-{date_time}.json"
        logger.warning(f"No image output file provided. Defaulting to: {output_fp}", method=handle_image_dict_output)

    file_type: str = cli_utils.get_file_type(output_fp)

    if file_type == "json":
        logger.info(f"Outputing image json to fp: '{output_fp}'")
        io_utils.output_to_json(output_fp, image_dict)
        return

    else:
        logger.fatal(f"Dicts cannot be written to files of type: '{file_type}'. Accepted file types: ('json')", method=handle_image_dict_output)
        exit(1)


def handle_image_bytes_output(image_bytes, output_fp):
    if output_fp is None:
        date_time = logger._format_date_time(Timestamp.DATE_AND_TIME)
        output_fp = f"image-{date_time}.png"
        logger.warning(f"No image output file provided. Defaulting to: {output_fp}.", method=handle_image_bytes_output)

    file_type: str = cli_utils.get_file_type(output_fp)

    if file_type == "png":
        logger.info(f"Outputing image bytes to fp: '{output_fp}'.")
        io_utils.output_image_to_file(output_fp, image_bytes)
        return

    elif file_type == 'jsonl':
        logger.info(f"Outputing image jsonl to fp: '{output_fp}'.")
        io_utils.output_image_to_jsonl(output_fp, image_bytes)
        return

    else:
        logger.fatal(f"Images cannot be written to files of type: '{file_type}'. Accepted file types: ('png', 'jsonl')", method=handle_image_bytes_output)
        exit(1)


def get_instance_type_name(item):
    item_class_type: type = type(item)

    item_type: str = item_class_type.__name__

    return item_type
