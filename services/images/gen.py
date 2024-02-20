from utils import openai_utils, cli_utils, io_utils
from utils.logger import logger, Timestamp


def add_service(parser):
    gen_parser = parser.add_parser('gen', help="Generate an image based on prompt. See `oait image gen --help`")

    gen_parser.description = "OpenAI Image Generation Tools" 
    gen_parser.usage = "oait image gen [-n NAME] [-d DESC] [-m MODEL] [-t TOOLS]"

    gen_parser.add_argument('--prompt', '-p', type=str, help="Prompt for image generation.")
    gen_parser.add_argument('--output', '-o', nargs='?', type=str, help="Output file ('png', 'json', or 'jsonl') (Pass only -o for default: 'image.json')", const="default", default=None)


def run_service(key, args):
    image_prompt: str = args.prompt
    output: str = args.output

    if image_prompt:
        logger.info(f"Generating image with prompt: '{image_prompt}'")
        image_data = generate_image_from_prompt(key, image_prompt)

    else:
        logger.fatal("No prompt received for image generation.")
        exit(1)

    logger.debug("Outputing image data.")
    handle_image_output(image_data, output)


def generate_image_from_prompt(key: str, prompt: str):
    logger.info("Generating image...")
    image_res = openai_utils.generate_image(key, prompt)

    if not image_res:
        logger.fatal("No content returned from api.", method=generate_image_from_prompt)
        exit(1)

    image_data = image_res.get('data')
    if not image_data:
        logger.fatal("No data returned from api", method=generate_image_from_prompt)

    logger.info(f"Image data returned from image generation.")
    logger.debug(f"Image data: {image_data}")
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
