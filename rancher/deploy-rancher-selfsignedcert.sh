#!/bin/bash
if [[ $DEBUG == "true" ]]; then
  set -x
fi

# Check if FQDN is given
if [ -z "$1" ]; then
    echo "Usage: $0 rancher.yourdomain.com 127.0.0.1,your-ip"
    exit 1
fi

# Set config here
export FQDN=$1
export CA_SUBJECT="My own root CA"
export CA_EXPIRE="1825" # CA expires in 5 years
export SSL_EXPIRE="365" # Certificate expires in 1 year
export SSL_SUBJECT="${FQDN}"
export SSL_DNS="${FQDN}" # Additional SANs (comma separated) can be added
export SSL_IP="$2" # Additional IPs (comma separated) can be added
export SILENT="true"

# Due to this open PR (https://github.com/paulczar/omgwtfssl/pull/10) I changed to use the edited version of the Docker image under superseb/omgwtfssl. Of course with appropriate referral in the description.
docker run -v --rm "$PWD/certs:/certs" \
  -e CA_SUBJECT \
  -e CA_EXPIRE \
  -e SSL_EXPIRE \
  -e SSL_SUBJECT \
  -e SSL_DNS \
  -e SSL_IP \
  -e SILENT \
  superseb/omgwtfssl

docker network create \
  --driver bridge \
  --subnet 192.168.0.0/24 \
  rancher_net

docker run -d --restart=unless-stopped \
  -p 80:80 -p 443:443 \
  --name rancher \
  --privileged \
  --log-driver json-file \
  --log-opt max-size=10m \
  --log-opt max-file=3 \
  --network rancher_net \
  --ip 192.168.0.10 \
  --cpus="2" \
  --memory="4g" \
  -v "$PWD/data:/var/lib/rancher" \
  -v "$PWD/certs/cert.pem:/etc/rancher/ssl/cert.pem" \
  -v "$PWD/certs/key.pem:/etc/rancher/ssl/key.pem" \
  -v "$PWD/certs/ca.pem:/etc/rancher/ssl/cacerts.pem" \
  rancher/rancher:latest

# echo "Waiting for Rancher to be started"
# while true; do
#   docker run --rm --net=host appropriate/curl -sLk "https://$FQDN/ping" && break
#   echo -n "."
#   sleep 5
# done
#
# echo ""
#
# docker run --rm --net=host superseb/rancher-check "https://${FQDN}"
