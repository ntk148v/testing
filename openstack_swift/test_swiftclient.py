from keystoneauth1.identity import v3
from keystoneauth1 import session
from swiftclient.client import Connection


auth = v3.Password(auth_url='http://192.168.122.68:5000/v3',
                   user_domain_name='default',
                   username='admin',
                   password='admin@123',
                   project_domain_name='default',
                   project_name='admin')
sess = session.Session(auth=auth)
try:
    client = Connection('2', session=sess)
    # print(client.get_account())
    print(client.put_container('fake-container'))
    # print(client.post_container('fake-container-1', headers={'x-container-meta-new': 'value'}))
    with open('/home/kiennt/Pictures/XIEC_2014_029.jpg', 'r') as local:
        print(client.put_object(
            'fake-container',
            'fake-obj',
            contents=local,
        ))
    # print(client.get_container('fake-container'))
    # print(client.post_object('fake-container', 'fake-obj', {'x-object-meta-newkey': 'newvalue'}))
    # print(client.head_object('fake-container', 'fake-obj'))
    print(client.copy_object('fake-container', 'fake-obj',
                             headers=None, destination='/fake-container/fake-obj'))
finally:
    client.close()
