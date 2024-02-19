from os.path import splitext
from time import sleep

def get_file_type(file_path: str) -> str:
    return splitext(file_path)[1][1:]


def get_file_without_extension(file_path: str) -> str:
    return splitext(file_path)[0]


def update_progress(progress, total):
    print(f"Progress:[{progress}/{total}]", end="\r") 


def display_progress_dots(num_dots: int):
    print("."* num_dots, end="\r")


def cycle_progress_dots(num_dots: int, delay_s: int):
    while True:
        for i in range(num_dots):
            display_progress_dots(i + 1)
            sleep(delay_s / num_dots)


