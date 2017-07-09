from keystoneauth1.identity import v3
from keystoneauth2 import session
from glanceclient.client import Client

auth = v3.Password(auth_url='http://25.28.148.127:5000/v3',
                  user_domain_name='default',
                  username='admin',
                  password='bkcloud',
                  project_domain_name='default',
                  project_name='admin')

sess = session.Session(auth=auth)
client = Client('2', session=sess)
print(client)
