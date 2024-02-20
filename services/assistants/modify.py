from utils.logger import logger

from . import openai_utils


def add_service(parser):
    modify_parser = parser.add_parser('modify', help="Modify an assistant. See `oait asst modify --help`")

    modify_parser.description = "OpenAI Assistant Modify Tools. Update assistant details." 
    modify_parser.usage = "oait asst modify <assistant_id> [-n NAME] [-d DESC] [-m MODEL] [-t TOOLS]"

    modify_parser.add_argument('asst', type=str, help="Assistant id to modify.")
    modify_parser.add_argument('--name', '-n', type=str, help="Name of updated assistant.", default=None)
    modify_parser.add_argument('--instructions', '-i', type=str, help="Instructions for updated assistant ", default=None)
    modify_parser.add_argument('--description', '-d', type=str, help="Description of updated assistant")
    modify_parser.add_argument('--tools', '-t', type=str, help="Tools for updated assistant (code, retrieval, functions)")   
    modify_parser.add_argument('--model', '-m', type=str, help="Model to use for updated assistant", default=None)


def run_service(key, args):
    asst: str = args.asst
    name: str = args.name
    instructions: str = args.instructions
    description: str = args.description
    model: str =  args.model
    tools: str = args.tools

    if not asst:
        logger.fatal("Must provide assistant id to modify.", method=handle_assistant_modify)
        exit(1)

    assistant_config = {}
    if name:
        assistant_config['name'] = name
    if instructions:
        assistant_config['instructions'] = instructions
    if description:
        assistant_config['description'] = description
    if model:
        assistant_config['model'] = model
    if tools:
        assistant_config['tools'] = tools

    if len(assistant_config) <= 0:
        logger.fatal("Must pass a field to update assistant. See `oait asst mod -h`")
        exit(0)

    logger.info(f"Creating assistant with config: '{assistant_config}'...")

    asst_id: str = modify_assistant_from_config(key, asst, assistant_config)

    logger.info(f"Created assistant (assistant_id: '{asst_id}').")


def modify_assistant_from_config(key: str, asst_id: str, assistant_config: dict):
    logger.info("Updating assistant...")

    response = openai_utils.modify_assistant(key, asst_id, assistant_config)

    if not response:
        logger.fatal("No data returned from OpenAI.", method=modify_assistant_from_config)

    res_asst_id: str | None = response.get('id')

    if not res_asst_id:
        logger.fatal("No assistant id returned from OpenAI", method=modify_assistant_from_config)

    logger.info(f"OpenAI assistant id updated: '{res_asst_id}'")

    return res_asst_id
