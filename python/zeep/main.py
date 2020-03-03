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
        kwargs['gnoc_url'] = os.environ.get('GNOC_URL')
    except KeyError as e:
        LOG.error('Missing environment variables! %s' % str(e))
    soap_client = Client(kwargs['gnoc_url'])
    srdto_type = soap_client.get_type('ns0:srdto')
    srdto_obj = srdto_type(country='281',
                           offset=100,
                           srId='26347',)
    LOG.info(soap_client.service.getListSR(rowStart=0, maxRow=100,
                                           srDTO=srdto_obj))
