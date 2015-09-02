FROM quay.io/opsee/vinz:latest

ENV VAPE_PUBLIC_HOST=":8081"
ENV VAPE_PRIVATE_HOST=":9091"
ENV VAPE_KEYFILE="/vape.test.key"
ENV VAPE_CERT=""
ENV VAPE_CERT_KEY=""
ENV POSTGRES_CONN="host=postgresql user=postgres dbname=vape_test sslmode=disable"
ENV AWS_ACCESS_KEY_ID=""
ENV AWS_SECRET_ACCESS_KEY=""
ENV AWS_DEFAULT_REGION=""
ENV AWS_INSTANCE_ID=""
ENV AWS_SESSION_TOKEN=""

RUN apk add --update bash ca-certificates curl
RUN curl -Lo /opt/bin/migrate https://s3-us-west-2.amazonaws.com/opsee-releases/go/migrate/migrate-linux-amd64 && \
    chmod 755 /opt/bin/migrate
RUN curl -Lo /opt/bin/ec2-env https://s3-us-west-2.amazonaws.com/opsee-releases/go/ec2-env/ec2-env && \
    chmod 755 /opt/bin/ec2-env

COPY run.sh /
COPY target/linux/amd64/bin/* /
COPY vape.test.key /

EXPOSE 8081 9091
CMD ["/vape"]
