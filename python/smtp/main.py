from email import header
from email.mime import text
import logging
import os
import smtplib
import sys

LOG = logging.getLogger(__name__)


def setup_logging():
    stdout_handler = logging.StreamHandler(sys.stdout)
    logging.basicConfig(
        level=logging.INFO,
        format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        handlers=[stdout_handler])


if __name__ == '__main__':
    setup_logging()
    try:
        smtp_server = os.environ.get('SMTP_SERVER')
        smtp_username = os.environ.get('SMTP_USERNAME')
        smtp_password = os.environ.get('SMTP_PASSWORD')
        smtp_receivers = os.environ.get('SMTP_RECEIVERS').split(',')
    except KeyError as e:
        LOG.error('Missing environment variables! %s' % str(e))
    s = smtplib.SMTP(smtp_server)
    message = text.MIMEText('Test', _charset='utf-8')
    message['Subject'] = header.Header('Test with StartTLS', 'utf-8')
    message['From'] = smtp_username
    s.ehlo()
    s.starttls()
    s.ehlo()
    s.login(smtp_username, smtp_password)
    s.sendmail(from_addr=smtp_username, to_addrs=smtp_receivers,
               msg=message.as_string())
    message = text.MIMEText('Test', _charset='utf-8')
    message['Subject'] = header.Header('Test without StartTLS', 'utf-8')
    message['From'] = smtp_username
    s.ehlo()
    s.login(smtp_username, smtp_password)
    s.sendmail(from_addr=smtp_username, to_addrs=smtp_receivers,
               msg=message.as_string())
