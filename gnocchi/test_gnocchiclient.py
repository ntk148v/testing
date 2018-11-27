from datetime import datetime
import os

from gnocchiclient import client
from gnocchiclient import exceptions
from keystoneauth1 import loading
from keystoneauth1 import session


def create_gnocchi_client(auth_url, username, password,
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
    gnocchi_client = create_gnocchi_client(
        auth_url=os.environ.get('OS_AUTH_URL'),
        username=os.environ.get('OS_USERNAME'),
        password=os.environ.get('OS_PASSWORD'),
        project_name=os.environ.get('OS_PROJECT_NAME'),
        user_domain_name=os.environ.get('OS_USER_DOMAIN_NAME'),
        project_domain_name=os.environ.get('OS_PROJECT_DOMAIN_NAME')
    )
    # existed_resources = gnocchi_client.resource.list(resource_type='instance')
    # for resource in existed_resources:
    #     resource_id = resource['id']
    #     print('******************************************************')
    #     if resource['server_group'] == 'vsmart-configuration-service' and resource['metrics'].get('cpu') is not None:
    #         metric_cpu_id = resource['metrics']['cpu']
    #         try:
    #             print('Deleting metric cpu {}' . format(metric_cpu_id))
    #             gnocchi_client.metric.delete(metric_cpu_id, resource_id)
    #             print('Deleted metric cpu {}' . format(metric_cpu_id))
    #         except exceptions.MetricNotFound as e:
    #             print('Ignoring error {}' . format(str(e)))
    #             pass
    #     print('Deleting resource {}' . format(resource_id))
    #     gnocchi_client.resource.delete(resource_id)
    #     print('Deleted resource {}' . format(resource_id))
    #     print('******************************************************')

    # query = {
    #     'and': [{'=': {'server_group': 'vsmart-configuration-service'}},
    #             {'=': {'ended_at': None}},
    #             {'=': {'project_id': 'c5a8b5960ac04cc68f18a541a7a9c51e'}}]}
    query = {
        '=': {
            'server_group': 'vsmart-wirez-service'
        }
    }
    # result1 = gnocchi_client.aggregates.fetch(
    #     # operations='(/ (aggregate mean (metric cpu rate:mean)) 10e9)',
    #     operations='(metric cpu rate:mean)',
    #     granularity=300,
    #     needed_overlap=0,
    #     search=query,
    #     resource_type='instance'
    # )
    # print(result1['measures'].keys())
    # print(len(result1['measures']['aggregated']))
    # print(result1['measures']['aggregated'])
    # result2 = gnocchi_client.metric.aggregation(
    #     metrics='cpu_util',
    #     granularity=300,
    #     needed_overlap=0,
    #     aggregation='mean',
    #     query=query,
    #     resource_type='instance',
    # )
    # print(result2)
    # print(len(result2))
    # resources = gnocchi_client.resource.update(
    #     'instance',
    #     '469ddc02-d47e-4f3b-9fe8-1ec0697e23c2',
    #     {'server_group': 'vsmart-configuration-service'}
    # )
    resources = gnocchi_client.resource.search(
        resource_type='instance',
        query=query,
        details=True
    )
    print(resources)