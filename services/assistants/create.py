from utils import openai_utils
from utils.logger import logger


def add_service(parser):
    create_parser = parser.add_parser('create', help="Create an assistant. See `oait asst create --help`")

    create_parser.description = "OpenAI Assistant Create Tools." 
    create_parser.usage = "oait asst create [-n NAME] [-d DESC] [-m MODEL] [-t TOOLS]"

    create_parser.add_argument('--name', '-n', type=str, help="Name of created assistant.", default=None)
    create_parser.add_argument('--instructions', '-i', type=str, help="Instructions for created assistant", default=None)
    create_parser.add_argument('--description', '-d', type=str, help="Description of created assistant", default=None)
    create_parser.add_argument('--model', '-m', type=str, help="Model used by created assistant (Default 'gpt-3.5-turbo-0125')", default="gpt-3.5-turbo-0125")
    create_parser.add_argument('--tools', '-t', type=str, help="Tools to provide to created assistant (code, retrieval, functions)", default=None)


def run_service(key, args):
    name: str = args.name
    instructions: str = args.instructions
    description: str = args.description
    model: str =  args.model
    tools: str = args.tools

    if not name:
        logger.fatal(f"Assistant name must be passed for creation.", method=handle_assistant_create)
        exit(1)
    if not model:
        logger.fatal(f"Assistant model must be passed for creation.", method=handle_assistant_create)
        exit(1)
    if not instructions:
        logger.fatal(f"Assistant instructions must be passed for creation.", method=handle_assistant_create)
        exit(1)

    assistant_config: dict = { 'name': name, 'model': model, 'instructions': instructions }

    if description:
        assistant_config['description'] = description
    if tools:
        assistant_config['tools'] = tools

    logger.info(f"Creating assistant with config: '{assistant_config}'...")

    asst_id: str = create_assistant_from_config(key, assistant_config)

    logger.info(f"Created assistant (assistant_id: '{asst_id}').")


def create_assistant_from_config(key: str, assistant_config: dict):
    logger.info("Creating assistant...")

    response = openai_utils.create_assistant(key, assistant_config)

    if not response:
        logger.fatal("No data returned from OpenAI.", method=create_assistant_from_config)

    asst_id: str | None = response.get('id')

    if not asst_id:
        logger.fatal("No assistant id returned from OpenAI", method=create_assistant_from_config)

    logger.info(f"OpenAI assistant id created: '{asst_id}'")

    return asst_id
