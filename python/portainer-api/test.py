import json
import logging
import operator
import simplejson
import time

import requests

LOG = logging.getLogger(__name__)


class PortainerAPIClient(object):
    jwt = None
    jwt_expired_at = None
    headers = None

    def __init__(self, endpoint, username, password):
        self.endpoint = endpoint.strip('/')
        self.username = username
        self.password = password

        self.get_jwt_token()
        self.jwt_expired_at = time.time() + 8 * 3600

    def get_jwt_token(self):
        try:
            r = requests.post(
                self.endpoint + "/api/auth",
                data=json.dumps({"Username": self.username,
                                 "Password": self.password}))
            r.raise_for_status()
            jwt = r.json()
            self.jwt = jwt.get('jwt')
            self.headers = {"Authorization": "Bearer {}". format(self.jwt)}
        except Exception as e:
            LOG.exception("Unable to authenticate a user: {}" . format(e))

    class decorator(object):
        @staticmethod
        def refresh_jwt_token(decorated_func):
            def wrapper(api, *args, **kwargs):
                if time.time() > api.jwt_expired_at:
                    api.get_jwt_token()
                return decorated_func(api, *args, **kwargs)
            return wrapper

    def _process_response(self, response):
        try:
            response.raise_for_status()
        except requests.exceptions.HTTPError as exec:
            err_msg = str(exec)

            # Attempt to get Error message from response
            try:
                error_dict = response.json()
            except (json.decoder.JSONDecodeError, simplejson.errors.JSONDecodeError):
                pass
            else:
                err_msg += " [Error: {}]".format(error_dict)
            raise requests.exceptions.HTTPError(err_msg)
        else:
            return response.json()

    @decorator.refresh_jwt_token
    def get_status(self):
        r = requests.get(
            self.endpoint + "/api/status",
            headers=self.headers)
        return self._process_response(r)

    @decorator.refresh_jwt_token
    def get_endpoints(self):
        try:
            r = requests.get(
                self.endpoint + "/api/endpoints",
                headers=self.headers)
            r.raise_for_status()
            return r.json()
        except Exception as e:
            LOG.exception("Unable to get endpoints: %s" % e)

    @decorator.refresh_jwt_token
    def get_endpoint(self, id):
        try:
            r = requests.get(
                self.endpoint + "/api/endpoints/{}" . format(id),
                headers=self.headers)
            r.raise_for_status()
            return r.json()
        except Exception as e:
            LOG.exception("Unable to get a endpoint {}: {}" . format(id, e))

    @decorator.refresh_jwt_token
    def get_stacks(self):
        try:
            r = requests.get(
                self.endpoint + "/api/stacks",
                headers=self.headers)
            r.raise_for_status()
            return r.json()
        except Exception as e:
            LOG.exception("Unable to get stacks: {}" . format(e))

    @decorator.refresh_jwt_token
    def get_stack(self, id):
        try:
            r = requests.get(
                self.endpoint + "/api/stacks/{}" . format(id),
                headers=self.headers)
            r.raise_for_status()
            return r.json()
        except Exception as e:
            LOG.exception("Unable to get a stack {}: {}" . format(id, e))

    @decorator.refresh_jwt_token
    def create_stack(self, endpoint_id, type, method, body=None, files=None):
        try:
            # self.headers['Content-Type'] = 'multipart/form-data'
            r = requests.post(
                self.endpoint + "/api/stacks?endpointId={}&method={}&type={}" .format(
                    endpoint_id, method, type),
                data=body,
                files=files,
                headers=self.headers)
            # r.raise_for_status()
            print(r.request.body.decode())
            return r.json()
        except Exception as e:
            LOG.exception("Unable to create a stack: {}" . format(e))

    @decorator.refresh_jwt_token
    def update_stack(self, id, endpoint_id=None, body=None):
        url = self.endpoint + "/api/stacks/{}" . format(id)
        if endpoint_id is not None:
            url += "?endpointId={}" . format(endpoint_id)
        r = requests.put(
            url,
            json=body,
            headers=self.headers)
        return self._process_response(r)

    @decorator.refresh_jwt_token
    def delete_stack(self, id, external=None, endpoint_id=None):
        try:
            url = self.endpoint + "/api/stacks/{}" . format(id)
            if endpoint_id is not None:
                url += "?endpointId={}" . format(endpoint_id)
            r = requests.delete(
                url,
                headers=self.headers)
            r.raise_for_status()
            return r.json()
        except Exception as e:
            LOG.exception("Unable to delete a stack {}: {}" . format(id, e))

    @decorator.refresh_jwt_token
    def get_stack_file(self, id):
        try:
            r = requests.get(
                self.endpoint + "/api/stacks/{}/file" . format(id),
                headers=self.headers)
            return r.json()
        except Exception as e:
            LOG.exception("Unable to get a stack file {}: {}" . format(id, e))


if __name__ == '__main__':
    p = PortainerAPIClient('http://127.0.0.1:9000',
                           'admin', '12345678')
    # print(p.get_status())
    with open('docker-compose.yml', 'rb') as f:
        # body = {
        #     'Name': 'test2',
        #     'Env': "[{\"name\": \"tag\", \"value\": \"3\"}]",
        # }
        # xitpe = json.dumps([{'name': 'tag', 'value': '3'}])
        # print(xitpe)
        # body = {
        #     "Name": "test",
        #     "Env": [{'name': 'tag', 'value': '3'}],
        # }
        # print(body)
        # files = {'file': f}
        # print(p.create_stack(1, 2, 'file', body=json.dumps(body), files=files))
        sf = p.get_stack_file(23)
        update = sf
        # update = {}
        update['Env'] = [
            {
                'name': 'tag',
                'value': '3'
            }
        ]
        print(p.update_stack(23, 1, body=update))
