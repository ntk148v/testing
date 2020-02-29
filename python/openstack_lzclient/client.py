import os

from keystoneauth1 import loading
from keystoneauth1 import session
from oslo_utils import importutils


class OpenStackLzClient(object):
    def __init__(self):
        super(OpenStackLzClient, self).__init__()
        self.created_clients = {}
        self._sess = self._create_keystone_session()

    def _create_keystone_session(self):
        loader = loading.get_plugin_loader('password')
        try:
            auth = loader.load_from_options(
                auth_url=os.environ['OS_AUTH_URL'],
                username=os.environ['OS_USERNAME'],
                password=os.environ['OS_PASSWORD'],
                project_name=os.environ['OS_PROJECT_NAME'],
                user_domain_name=os.environ['OS_USER_DOMAIN_NAME'],
                project_domain_name=os.environ['OS_PROJECT_DOMAIN_NAME']
            )
        except KeyError as err:
            raise err
        return session.Session(auth=auth)

    def _import_client(self, module_name):
        """Import Client with a given module name. for example:
        cinderclient.client, neutronclient.v2_0.client
        """
        try:
            return importutils.import_module(module_name)
        except ModuleNotFoundError as err:
            raise err

    def create_client(self, module_name, version, client_class):
        """Create a OpenStack module client"""
        # NOTE(kiennt): Get created client rather create a new one.
        #               The key is the combination of module_name and version.
        #               because we can create multiple clients of a module with
        #               different versions.
        client = self.created_clients.get(module_name + version)
        if client:
            return client
        module_client = self._import_client(module_name)
        try:
            client = getattr(module_client, client_class)(
                version=version,
                session=self._sess)
            self.created_clients[module_name+version] = client
            return client
        except Exception as err:
            raise err

    def exec_method(self, module_name, version=None, client_class=None,
                    method_name=None, *args, **kwargs):
        """Execute method of a given project module

        :param method_name(str): the method/attr want to be executed.
                                 for example "servers.list, servers.find".
        :param module_name(str): the module name, for example
                                 for example, "novaclient.client".
        :param version(str): the version of module client.
        :param client_class(str): the name of Client class, usually be 'Client'.
        :param *args, **kwargs: the agruments that will be passed when
                                execute method call.
        """
        client_class = client_class or 'Client'
        client_version = version or 2
        _client = self.create_client(module_name, client_version,
                                     client_class)
        try:
            # NOTE(kiennt): method_name could be a combination
            #               for example 'servers.list'. Here is the
            #               workaround.
            method = getattr(_client, method_name.split('.')[0])
            for attr in method_name.split('.')[1:]:
                method = getattr(method, attr)
            return method(*args, **kwargs)
        except Exception as err:
            raise err
