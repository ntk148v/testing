import os

from calplus.client import Client
from calplus.provider import Provider


def get_env():
    cloud_type = os.getenv('CLOUD_TYPE')
    cloud_config = {
        'os_auth_url': os.getenv('OS_AUTH_URL'),
        'os_project_name': os.getenv('OS_PROJECT_NAME'),
        'os_username': os.getenv('OS_USERNAME'),
        'os_password': os.getenv('OS_PASSWORD'),
        'os_project_domain_name': os.getenv('OS_PROJECT_DOMAIN_NAME'),
        'os_user_domain_name': os.getenv('OS_USER_DOMAIN_NAME'),
        'os_identity_api_version': os.getenv('OS_IDENTITY_API_VERSION'),
        'os_auth_version': os.getenv('OS_AUTH_VERSION'),
        'os_swiftclient_version': os.getenv('OS_SWIFTCLIENT_VERSION'),
    }
    return cloud_type, cloud_config


def main():
    cloud_type, cloud_config = get_env()
    provider = Provider(cloud_type, cloud_config)
    # Client for object storage
    client = Client(version='1.0.0', resource='object_storage',
                    provider=provider)
    try:
        print(client.list_containers())
        # Add more tests bellow
    except Exception as e:
        raise e


if __name__ == '__main__':
    main()
