import os

from seatable_api import Base
from seatable_api.constants import UPDATE_DTABLE

try:
    server_url = os.environ['SEATABLE_SERVER_URL']
    api_token = os.environ['SEATABLE_API_TOKEN']
except KeyError as err:
    exit('Required environment variables is missing: %s' % (err))

base = Base(api_token, server_url)
base.auth(with_socket_io=True)

def on_update_seatable(data, index, *args):
    print(data)

base.socketIO.on(UPDATE_DTABLE, on_update_seatable)
base.socketIO.wait()  # forever
