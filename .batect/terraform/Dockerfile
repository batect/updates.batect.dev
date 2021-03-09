FROM alpine:3.12.0

RUN apk add --no-cache \
    bash \
    ca-certificates \
    curl \
    jq \
    ncurses \
    python3

ENV TERRAFORM_VERSION=0.14.7
ADD https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip /tmp/terraform.zip
RUN unzip /tmp/terraform.zip -d /usr/local/bin && rm /tmp/terraform.zip

ARG CLOUD_SDK_VERSION
ADD https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-${CLOUD_SDK_VERSION}-linux-x86_64.tar.gz /tmp/google-cloud-sdk.tar.gz
RUN tar xzf /tmp/google-cloud-sdk.tar.gz && rm /tmp/google-cloud-sdk.tar.gz

ENV PATH /google-cloud-sdk/bin:$PATH

RUN gcloud components install beta
