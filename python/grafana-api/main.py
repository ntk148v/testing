import csv
import os

from grafana_api.grafana_face import GrafanaFace
from grafana_api.grafana_api import GrafanaServerError

gf_host = os.environ['GF_HOST']
gf_username = os.environ['GF_USERNAME']
gf_password = os.environ['GF_PASSWORD']

grafana_api = GrafanaFace(auth=(gf_username, gf_password), host=gf_host)
# open csv file
with open('users.csv') as csv_file:
    csv_reader = csv.reader(csv_file, delimiter=',')
    for row in csv_reader:
        if len(row) < 2:
            print('Invalid row format')
            continue
        name = row[0]
        email = row[1]
        username = email.strip('@gmail.com')
        user = {
            "name": name,
            "email": email,
            "login": username,
            "password": "defaultpassword",
            "OrgId": 1
        }
        try:
            result = grafana_api.admin.create_user(user)
            print(result)
        except GrafanaServerError:
            pass
        except Exception as e:
            raise e
