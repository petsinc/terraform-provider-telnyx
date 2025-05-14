# Use a minimal Go image
FROM golang:1.22-alpine

ARG PRODUCT=terraform
ARG VERSION=1.7.3
ARG ARCH=arm64

RUN apk add --update --virtual .deps --no-cache gnupg git bash curl just make && \
    cd /tmp && \
    wget https://releases.hashicorp.com/${PRODUCT}/${VERSION}/${PRODUCT}_${VERSION}_linux_${ARCH}.zip && \
    wget https://releases.hashicorp.com/${PRODUCT}/${VERSION}/${PRODUCT}_${VERSION}_SHA256SUMS && \
    wget https://releases.hashicorp.com/${PRODUCT}/${VERSION}/${PRODUCT}_${VERSION}_SHA256SUMS.sig && \
    wget -qO- https://www.hashicorp.com/.well-known/pgp-key.txt | gpg --import && \
    gpg --verify ${PRODUCT}_${VERSION}_SHA256SUMS.sig ${PRODUCT}_${VERSION}_SHA256SUMS && \
    grep ${PRODUCT}_${VERSION}_linux_${ARCH}.zip ${PRODUCT}_${VERSION}_SHA256SUMS | sha256sum -c && \
    unzip /tmp/${PRODUCT}_${VERSION}_linux_${ARCH}.zip -d /tmp && \
    mv /tmp/${PRODUCT} /usr/local/bin/${PRODUCT} && \
    rm -f /tmp/${PRODUCT}_${VERSION}_linux_${ARCH}.zip ${PRODUCT}_${VERSION}_SHA256SUMS ${VERSION}/${PRODUCT}_${VERSION}_SHA256SUMS.sig
# Set up environment variables
ENV GOCMD=go

ENV TF_ACC_TERRAFORM_PATH=/usr/local/bin/${PRODUCT}

RUN go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
# ADD the plugin docs to the PATH
ENV PATH=$PATH:/go/bin

# Create working directory
WORKDIR /workspace

# Default entrypoint, can be overridden
ENTRYPOINT [ "sh", "-c" ]
