from datetime import datetime
import os

from heatclient import client
from keystoneauth1 import loading
from keystoneauth1 import session


def create_heat_client(auth_url, username, password,
                       project_name, user_domain_name,
                       project_domain_name):
    # TODO(kiennt): Working with Keystone V3, should add
    #               check version here.
    loader = loading.get_plugin_loader('password')
    auth = loader.load_from_options(auth_url=auth_url,
                                    username=username,
                                    password=password,
                                    project_name=project_name,
                                    user_domain_name=user_domain_name,
                                    project_domain_name=project_domain_name)

    sess = session.Session(auth=auth)
    return client.Client('1', session=sess)


if __name__ == '__main__':
    # Don't try catch here - simple test
    heat_client = create_heat_client(
        auth_url=os.environ.get('OS_AUTH_URL'),
        username=os.environ.get('OS_USERNAME'),
        password=os.environ.get('OS_PASSWORD'),
        project_name=os.environ.get('OS_PROJECT_NAME'),
        user_domain_name=os.environ.get('OS_USER_DOMAIN_NAME'),
        project_domain_name=os.environ.get('OS_PROJECT_DOMAIN_NAME')
    )

    start = datetime.now()
    print(start)
    stacks = heat_client.stacks.list()
    for stack in stacks:
        print(stack.id)
    end = datetime.now()
    print(end)
    print(end - start)

