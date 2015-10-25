#!/bin/bash
set -e

APPENV=${APPENV:-vapenv}

# relying on set -e to catch errors?
/opt/bin/ec2-env > /ec2env
eval "$(< /ec2env)"
/opt/bin/s3kms get -b opsee-keys -o dev/$APPENV > /$APPENV
/opt/bin/s3kms get -b opsee-keys -o dev/vape.key > /vape.key

# these will have to wait
# TODO: tls from load-balancer -> vape
# /opt/bin/s3kms get -b opsee-keys -o dev/vape-cert.pem > /vape-cert.pem
# /opt/bin/s3kms get -b opsee-keys -o dev/vape-key.pem > /vape-key.pem

source /$APPENV && \
	/opt/bin/migrate -url "$POSTGRES_CONN" -path /migrations up && \
	/vape
