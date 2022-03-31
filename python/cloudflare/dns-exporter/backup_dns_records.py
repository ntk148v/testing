"""Backup DNS records script using CloudFlare API v4."""
import logging
import os
import sys
import yaml

import CloudFlare

LOG = logging.getLogger(__name__)


def setup_logging(default_path='logging.yaml',
                  default_level=logging.INFO,
                  env_key='LOG_CFG'):
    """Setup logging configuration"""
    path = default_path
    value = os.getenv(env_key, None)
    if value:
        path = value
    if os.path.exists(path):
        with open(path, 'rt') as f:
            config = yaml.safe_load(f.read())
        logging.config.dictConfig(config)
    else:
        logging.basicConfig(
            format='%(asctime)s [%(levelname)s] %(message)s',
            level=default_level)


def export_dns_records():
    """CloudFlare API export DNS records"""
    try:
        LOG.info('Get required environment variables.')
        zone_name = os.environ['CF_ZONE_NAME']
        api_email = os.environ['CF_API_EMAIL']
        api_key = os.environ['CF_API_KEY']
        api_certkey = os.environ['CF_API_CERTKEY']
    except KeyError as err:
        LOG.error('Required environment variable is missing: %s' % (err))
        raise

    LOG.info('Initilize CloudFlare API Client.')
    cf = CloudFlare.CloudFlare(
        email=api_email,
        token=api_key,
        certtoken=api_certkey)
    # Grab the zone identifier
    try:
        LOG.info('Get the zone identifier with %s' % (zone_name))
        zones = cf.zones.get(params={'name': zone_name})
    except CloudFlare.exceptions.CloudFlareAPIError as err:
        LOG.error('/zones %d %s - api call failed' % (err, err))
        raise
    except Exception as err:
        LOG.error('/zones.get - %s - api call failed' % (err))
        raise

    if len(zones) == 0:
        LOG.error('/zones.get - %s - zone not found' % (zone_name))
        raise

    # The zone identifier should be unique
    if len(zones) != 1:
        LOG.error(
            '/zones.get - %s - api call return more than one items' % (zone_name))
        raise

    zone_id = zones[0]['id']
    try:
        LOG.info('Export DNS records with zone id - %s' % (zone_id))
        dns_records = cf.zones.dns_records.export.get(zone_id)
    except CloudFlare.exceptions.CloudFlareAPIError as err:
        LOG.error('/zones/dns_records/export %s - %d %s - api call failed' %
                  (zone_name, err, err))
        raise

    backupfile_path = os.environ.get('BACKUP_DNS_FILE', 'dns_records.bak')
    with open(backupfile_path, 'a') as backupfile:
        for line in dns_records.splitlines():
            if len(line) == 0 or line[0] == ';':
                # Blank line or comment line are skipped
                continue
            backupfile.write(line)


if __name__ == '__main__':
    setup_logging()
    export_dns_records()
