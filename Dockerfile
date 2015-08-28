FROM gliderlabs/alpine:3.2

ENV VAPE_PUBLIC_HOST=":8081"
ENV VAPE_PRIVATE_HOST=":9091"
ENV VAPE_KEYFILE="/vape.dev.key"
ENV POSTGRES_CONN="host=postgresql user=postgres dbname=vape_test sslmode=disable"
ENV AWS_ACCESS_KEY_ID=""
ENV AWS_SECRET_ACCESS_KEY=""

RUN apk add --update bash

COPY target/linux/amd64/bin/* /
COPY vape.dev.key /

EXPOSE 8081 9091
CMD ["/vape"]
