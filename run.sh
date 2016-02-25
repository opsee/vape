#!/bin/bash
set -e

APPENV=${APPENV:-vapenv}

/opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$APPENV > /$APPENV

source /$APPENV && \
  /opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/vape.key > /vape.key && \
  /opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$VAPE_CERT > /$VAPE_CERT && \
  /opt/bin/s3kms -r us-west-1 get -b opsee-keys -o dev/$VAPE_CERT_KEY > /$VAPE_CERT_KEY && \
  chmod 600 /$VAPE_CERT_KEY && \
	/opt/bin/migrate -url "$POSTGRES_CONN" -path /migrations up && \
	/vape
