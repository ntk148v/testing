import os
import boto3
from dotenv import load_dotenv

def main():
    load_dotenv()
    
    params = {
        'aws_access_key_id': os.getenv('AWS_ACCESS_KEY_ID'),
        'aws_secret_access_key': os.getenv('AWS_SECRET_ACCESS_KEY'),
        'endpoint_url': os.getenv('S3_ENDPOINT_URL')
    }

    client = boto3.client('s3', **params)
    
    try:
        bucket_name = 'migration'
        print(f"Listing objects in bucket: {bucket_name}")
        response = client.list_objects_v2(Bucket=bucket_name)
        
        objects = response.get('Contents', [])
        if not objects:
            print(f"No objects found in {bucket_name}.")
            return

        # Filter out directories (keys ending with /)
        files = [obj for obj in objects if not obj['Key'].endswith('/')]
        if not files:
            print(f"No files found in {bucket_name} (only directories).")
            return

        for obj in files:
            print(f"- {obj['Key']} ({obj['Size']} bytes)")

        # Attempt to download the first actual file
        first_file = files[0]['Key']
        print(f"\nAttempting to download first file: {first_file}")
        
        # Ensure local directory exists
        local_dir = os.path.dirname(first_file)
        if local_dir:
            os.makedirs(local_dir, exist_ok=True)
            
        client.download_file(bucket_name, first_file, first_file)
        print(f"Successfully downloaded {first_file}")

    except Exception as e:
        print(f"Error: {e}")


if __name__ == "__main__":
    main()
