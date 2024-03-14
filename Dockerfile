FROM redhat/ubi9-minimal:9.1 AS build

ARG GO_VERSION

# This copy assumes `go mod tidy` has already been run in the repo.
# In this way the docker file does not have to deal with private repositories.
COPY --chown=1001:0 . /build

# golang
RUN microdnf update -y --nodocs \
    && microdnf install -y --nodocs --setopt=install_weak_deps=0 --disableplugin=subscription-manager \
    make git gcc tar jq which findutils\
    && microdnf clean all --disableplugin=subscription-manager

RUN curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz -o go${GO_VERSION}.linux-amd64.tar.gz
RUN echo ${GO_VERSION}
RUN tar -C /usr/local -xf go${GO_VERSION}.linux-amd64.tar.gz \
    && ln -s /usr/local/go/bin/go /usr/bin/go

# This code is fix the docker failed issue of dubious ownership of build folder
RUN git config --global --add safe.directory /build

# Have to change the npm prefix, otherwise npm install inside of make build fails with permissions errors
RUN cd /build && \
    export GOPATH=$(go env GOPATH) && \
    export PATH=$GOPATH/bin:$PATH && \
    go mod tidy && \
    go mod vendor && \
    make build-light    

#----------------------------------------------------------------
FROM  redhat/ubi9-minimal:9.1 AS runtime

RUN microdnf update -y --nodocs && microdnf clean all \
    && microdnf install --nodocs -y shadow-utils openssh-clients \
    && groupadd -g 1000 memeuser && adduser -u 1000 -g memeuser memeuser && chmod 755 /home/memeuser \
    && microdnf install --nodocs -y openssl ca-certificates gettext \
    && microdnf clean all \
    && mkdir -p /licenses

USER memeuser
WORKDIR /home/memeuser

COPY --from=build --chown=memeuser:memeuser /build/meme /home/memeuser/

# Commented out, because the port will be parameterized to work with podman shared networks
# EXPOSE 7120/tcp

ENTRYPOINT ["/home/memeuser/meme"]
