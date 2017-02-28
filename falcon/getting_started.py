import falcon


class ThingsResource(object):

    def on_get(self, req, resp):
        resp.status = falcon.HTTP_200
        resp.body = ('Hello Word')

# Falcon.API instances are callable WSGI apps
app = falcon.API()

# Resources are represented by long-lived class instances
things = ThingsResource()

# things will handle all requests to the '/thing' URL path
app.add_route('/things', things)
