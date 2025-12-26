FROM golang:1.23-bookworm AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o dpbuild .

FROM debian:trixie-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    pbuilder \
    dpkg-dev \
    ubuntu-keyring \
    debian-archive-keyring \
    sudo \
    && rm -rf /var/lib/apt/lists/*

RUN cat <<EOF > /root/.pbuilderrc
BUILD_HOME="\$BUILDDIR"
APTCACHEHARDLINK=no
EOF

COPY --from=builder /build/dpbuild /usr/local/bin/dpbuild

WORKDIR /work

ENTRYPOINT ["dpbuild"]
