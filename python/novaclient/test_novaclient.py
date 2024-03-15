import os

from keystoneauth1 import loading
from keystoneauth1 import session
from novaclient import client

def create_nova_client(auth_url , username, password, project_name,
                       user_domain_name, project_domain_name):
    # TODO(kiennt): Working with Keystone v3
    loader = loading.get_plugin_loader('password')
    auth = loader.load_from_options(auth_url=auth_url,
                                    username=username,
                                    password=password,
                                    project_name=project_name,
                                    user_domain_name=user_domain_name,
                                    project_domain_name=project_domain_name)
    sess = session.Session(auth=auth)
    return client.Client('2.26', session=sess)

if __name__ == '__main__':
    nova_client = create_nova_client(
        auth_url=os.environ.get('OS_AUTH_URL'),
        username=os.environ.get('OS_USERNAME'),
        password=os.environ.get('OS_PASSWORD'),
        project_name=os.environ.get('OS_PROJECT_NAME'),
        user_domain_name=os.environ.get('OS_USER_DOMAIN_NAME'),
        project_domain_name=os.environ.get('OS_PROJECT_DOMAIN_NAME')
    )

    a = nova_client.servers.list(search_opts={"tags": "test"})
    print(a)
