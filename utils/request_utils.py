import requests


def get_image_from_url(url: str):
    response = requests.get(url)

    if response.status_code == 200:
        return response.content

    else:
        error_message: str = str(response)

        print("ERROR: " + error_message)
        return None

