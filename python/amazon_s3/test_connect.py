import boto3


params = {
    'aws_secret_access_key': '5f53d8040aba43dfbdc6a286908a8adf',
    'aws_access_key_id': 'a0d6844dd463439d98bdf0f0dcf13233',
    'endpoint_url': 'http://192.168.122.115:8080'
}

client = boto3.client('s3', **params)
# print(client.delete_object(Bucket='fake-container', Key='fake-obj'))
# print(client.create_bucket(Bucket='fake-container-2'))
# print(client.list_buckets())
# print(client.list_objects(Bucket='fake-container'))
# print(client.head_object(Bucket='fake-container', Key='fake-obj'))
# print(client.copy_object(Bucket='fake-container', Key='fake-obj',
#                          Metadata={'orig-filename': 'fake-obj', 'x-amz-test3': 'another-testing'},
#                          MetadataDirective='COPY',
#                          CopySource={'Bucket': 'fake-container', 'Key': 'fake-obj'}))
print(client.list_buckets())
