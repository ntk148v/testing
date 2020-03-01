import logging
import os
import sys

from zeep import Client

LOG = logging.getLogger(__name__)


def setup_logging():
    stdout_handler = logging.StreamHandler(sys.stdout)
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        handlers=[stdout_handler])


if __name__ == '__main__':
    setup_logging()
    kwargs = {}
    try:
        kwargs['gnoc_url'] = os.environ.get('gnoc_url')
    except KeyError as e:
        LOG.error('Missing environment variables! %s' % str(e))
    soap_client = Client(kwargs['gnoc_url'])
    LOG.info(soap_client.service.getListSR(rowStart=0, maxRow=2))
