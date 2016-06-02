FROM alpine:3.3

RUN apk add --update bash ca-certificates curl
RUN mkdir -p /opt/bin && \
		curl -Lo /opt/bin/s3kms https://s3-us-west-2.amazonaws.com/opsee-releases/go/vinz-clortho/s3kms-linux-amd64 && \
    chmod 755 /opt/bin/s3kms && \
    curl -Lo /opt/bin/migrate https://s3-us-west-2.amazonaws.com/opsee-releases/go/migrate/migrate-linux-amd64 && \
    chmod 755 /opt/bin/migrate

ENV VAPE_PUBLIC_HOST=":8081"
ENV VAPE_PRIVATE_HOST=":9091"
ENV VAPE_KEYFILE="/vape.test.key"
ENV VAPE_CERT="cert.pem"
ENV VAPE_CERT_KEY="key.pem"
ENV VAPE_SLACK_DOMAIN=""
ENV VAPE_SLACK_ADMIN_TOKEN=""
ENV VAPE_LAUNCH_DARKLY_TOKEN=""
ENV POSTGRES_CONN="postgres://postgres@postgresql/vape_test?sslmode=disable"
ENV MANDRILL_API_KEY=""
ENV INTERCOM_KEY=""
ENV CLOSEIO_KEY=""
ENV SLACK_ENDPOINT=""
ENV OPSEE_HOST="staging.opsy.co"
ENV GODEBUG="netdns=cgo"
ENV VAPE_SPANX_HOST ""
ENV APPENV=""

COPY run.sh /
COPY target/linux/amd64/bin/* /
COPY vape.test.key /
COPY migrations /migrations
COPY cert.pem /
COPY key.pem /

EXPOSE 8081 9091
CMD ["/vape"]
