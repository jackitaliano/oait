
from os.path import splitext

def get_file_type(file_path: str) -> str:
    return splitext(file_path)[1][1:]


def get_file_without_extension(file_path: str) -> str:
    return splitext(file_path)[0]

def update_progress(progress, total):
    print(f"Progress:[{progress}/{total}]", end="\r") 

