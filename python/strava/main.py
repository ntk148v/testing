import datetime
import os
import subprocess
import time

from dotenv import load_dotenv
import requests
from stravalib.client import Client

load_dotenv()  # take environment variables

client_id = os.getenv("CLIENT_ID")
client_secret = os.getenv("CLIENT_SECRET")
client = Client()
# * Run once *
authorize_url = client.authorization_url(
    client_id=client_id, redirect_uri="http://127.0.0.1:5000/authorization",
    scope=[
        'read', 'read_all', 'profile:read_all',
        'profile:write', 'activity:read', 'activity:read_all',
        'activity:write'
    ])
print(authorize_url)
# Open authorize_url in browser to get code
# Extract the code from your webapp response
# For example: http://127.0.0.1:5000/authorization?state=&code=b4ae6cf3de141a669b517e3afa2bb48554a826bf&scope=read,activity:read
subprocess.run("wl-copy", text=True, input=authorize_url)
code = input(
    "The above URL has been copied to clipboard. Go to it, and copy and paste the code: ")
token_response = client.exchange_code_for_token(
    client_id=client_id, client_secret=client_secret, code=code
)
access_token = token_response["access_token"]
refresh_token = token_response["refresh_token"]
expires_at = token_response["expires_at"]
# Now store that short-lived access token somewhere (a database?)
client.access_token = access_token
# You must also store the refresh token to be used later on to obtain another valid access token
# in case the current is already expired
client.refresh_token = refresh_token

# An access_token is only valid for 6 hours, store expires_at somewhere and
# check it before making an API call.
client.token_expires_at = expires_at

if time.time() > client.token_expires_at:
    refresh_response = client.refresh_access_token(
        client_id=client_id, client_secret=client_secret, refresh_token=client.refresh_token
    )
    access_token = refresh_response["access_token"]
    refresh_token = refresh_response["refresh_token"]
    expires_at = refresh_response["expires_at"]

athlete = client.get_athlete()
print(
    "For {id} ({firstname} {lastname}), I now have an access token {token}".format(
        id=athlete.id, firstname=athlete.firstname, lastname=athlete.lastname,
        token=access_token
    )
)

activities = client.get_activities(before=datetime.datetime(2024, 1, 1))
for activity in activities:
    print(activity.name)
