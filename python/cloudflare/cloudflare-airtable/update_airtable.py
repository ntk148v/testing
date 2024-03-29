"""Update environment to Airtable according CloudFlare"""
from datetime import datetime
import os

import CloudFlare
from airtable import Airtable


# Modify these varaibles if needed
# Airtable authentication
airtable_base_key = os.getenv('AIRTABLE_BASE_KEY')
airtable_table_name = os.getenv('AIRTABLE_TABLE_NAME')
airtable_api_key = os.getenv('AIRTABLE_API_KEY')
# Cloudflare authentication
cloudflare_email = os.getenv('CLOUDFLARE_EMAIL')
cloudflare_token = os.getenv('CLOUDFLARE_TOKEN')
cloudflare_zone_name = os.getenv('CLOUDFLARE_ZONE_NAME')
# Number of DNS records per page
cloudflare_dns_records_per_page = 500

print('Initilize CloudFlare API Client.')
cf = CloudFlare.CloudFlare(
    email=cloudflare_email,
    token=cloudflare_token)
# Grab the zone identifier
try:
    print('Get the zone identifier with %s' % (cloudflare_zone_name))
    zones = cf.zones.get(params={'name': cloudflare_zone_name})
except CloudFlare.exceptions.CloudFlareAPIError as err:
    exit('/zones %d %s - api call failed' % (err, err))
except Exception as err:
    exit('/zones.get - %s - api call failed' % (err))

if len(zones) == 0:
    exit('/zones.get - %s - zone not found' % (cloudflare_zone_name))

# The zone identifier should be unique
if len(zones) != 1:
    exit('/zones.get - %s - api call return more than one '
         'items' % (cloudflare_zone_name))

zone_id = zones[0]['id']
try:
    print('Get DNS records with zone id - %s' % (zone_id))
    dns_records = cf.zones.dns_records.get(
        zone_id, params={'per_page': cloudflare_dns_records_per_page,
                         'type': 'A'})
except CloudFlare.exceptions.CloudFlareAPIError as err:
    exit('/zones/dns_records/get %s - %d %s - api call failed' %
         (cloudflare_zone_name, err, err))

print('Initilize Airtable API Client')

airtable = Airtable(airtable_base_key, airtable_table_name,
                    api_key=airtable_api_key)

for dns_record in dns_records:
    dns_type = dns_record['type']
    dns_name = dns_record['name']
    dns_content = dns_record['content']
    dns_id = dns_record['id']
    dns_short_name = dns_name
    dns_modified_on = dns_record['modified_on']

    if dns_content == '35.224.177.97':
        dns_env = 'development'
    elif dns_content == '35.224.175.98':
        dns_env = 'staging'
    else:
        dns_env = 'unknown'
    # Format shortname
    dns_short_name = dns_short_name.replace('.clare.ai', '')
    dns_short_name = dns_short_name.replace('-demo', 'Demo')
    dns_short_name = dns_short_name.replace('-hkdemo', 'HKDemo')
    row = {
        'ShortName': dns_short_name,
        'Name': dns_name,
        'IP': dns_content,
        'Environment': dns_env,
        'LastUpdated': str(datetime.utcnow()),
        'ModifiedOn': dns_modified_on
    }
    # Search by name
    print('* Update environment for record %s' % dns_name)
    exist_row = airtable.search('Name', dns_name)
    if len(exist_row) == 0:
        # Insert new row
        airtable.insert(row)
    elif len(exist_row) == 1:
        last_modified = datetime.strptime(exist_row[0]['fields']['ModifiedOn'],
                                          '%Y-%m-%dT%H:%M:%S.%fZ')
        current_modified = datetime.strptime(dns_modified_on,
                                             '%Y-%m-%dT%H:%M:%S.%fZ')
        # Only update the row when dns record was modified.
        if last_modified < current_modified:
            # Replace exist row with new one
            airtable.replace(exist_row[0]['id'], row)
    else:
        exit('Airtable has multiple rows with name %s' % dns_name)

print('Updated!')
exit(0)
