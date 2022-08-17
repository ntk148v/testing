import json
import os

import requests
import seatable_api

try:
    server_url = os.environ['SEATABLE_SERVER_URL']
    api_token = os.environ['SEATABLE_API_TOKEN']
    table = os.environ['SEATABLE_TABLE']
    username = os.environ['SEATABLE_USERNAME']
    password = os.environ['SEATABLE_PASSWORD']
except KeyError as err:
    exit('Required environment variables is missing: %s' % (err))

base = seatable_api.Base(api_token, server_url)
base.auth()

# Test list-rows and query
# with open('/tmp/test-query-no-convert.json', 'w', encoding='utf-8') as f:
#     json.dump(base.query(
#         f'select * from {table} limit 10000000', convert=False), f)

# with open('/tmp/test-list-rows.json', 'w', encoding='utf-8') as f:
#     json.dump(base.list_rows(table), f)

# with open('/tmp/test-query-convert.json', 'w', encoding='utf-8') as f:
#     json.dump(base.query(
#         f'select * from {table} limit 10000000', convert=True), f)

# Test account


def parse_response(response):
    if response.status_code >= 400:
        raise ConnectionError(response.status_code, response.text)
    else:
        try:
            data = json.loads(response.text)
            return data
        except:
            pass


class AdminAccount(seatable_api.Account):
    def __init__(self, login_name, password, server_url):
        super().__init__(login_name, password, server_url)

    def _list_all_users_url(self, page, per_page):
        return f'{self.server_url}/api/v2.1/admin/users/?page={page}&per_page={per_page}'

    def _search_a_user_url(self, query, page, per_page):
        return f'{self.server_url}/api/v2.1/admin/search-user/?query={query}&page={page}&per_page={per_page}'

    def list_all_users(self, page=1, per_page=25):
        response = requests.get(self._list_all_users_url(page, per_page),
                                headers=self.token_headers,
                                timeout=self.timeout)

        return parse_response(response)

    def search_a_user(self, query, page=1, per_page=25):
        response = requests.get(self._search_a_user_url(query, page, per_page),
                                headers=self.token_headers,
                                timeout=self.timeout)
        return parse_response(response)


# Test!
account = AdminAccount(username, password, server_url)
account.auth()
# print(account.list_all_users())
user = account.search_a_user("sample@auth.local")['user_list']
if user:
    print('Found: ', user)
else:
    print('Not found :(')