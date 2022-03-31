import os

from seatable_api import Base

try:
    server_url = os.environ['SEATABLE_SERVER_URL']
    api_token = os.environ['SEATABLE_API_TOKEN']
except KeyError as err:
    exit('Required environment variables is missing: %s' % (err))

base = Base(api_token, server_url)
base.auth(with_socket_io=False)
print(base.get_metadata())
