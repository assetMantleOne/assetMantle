# syntax=docker/dockerfile:latest
FROM golang:1.14-buster as build

# Set up dependencies
ENV PACKAGES curl make git
ENV PATH=/root/.cargo/bin:$PATH

# Set working directory for the build
WORKDIR /usr/local/app

# Install minimum necessary dependencies
RUN apt update && apt install -y $PACKAGES

# Install Rust and wasm32 dependencies
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

# Add source files
RUN --mount=type=bind,source=.,rw \
  go mod download

# Build
RUN --mount=type=bind,source=.,rw \
  --mount=type=cache,target=/root/.cache \
  make install

FROM ubuntu
COPY --from=build '/go/pkg/mod/github.com/!cosm!wasm/go-cosmwasm@v0.10.0/api/libgo_cosmwasm.so' /lib/x86_64-linux-gnu/libgo_cosmwasm.so
COPY --from=build /go/bin/mantleNode /usr/bin/mantleNode
COPY --from=build /go/bin/mantleNode /usr/bin/mantleNode
ENTRYPOINT ["mantleNode"]
