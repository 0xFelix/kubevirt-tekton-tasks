FROM registry.access.redhat.com/ubi9/ubi-minimal AS builder
RUN microdnf install -y tar gzip && microdnf clean all
ENV TASK_NAME=modify-data-object \
    GOFLAGS="-mod=vendor"
WORKDIR /src/${TASK_NAME}
RUN curl -L https://go.dev/dl/go1.19.2.linux-amd64.tar.gz | tar -C /usr/local -xzf -
ENV PATH=$PATH:/usr/local/go/bin
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /${TASK_NAME} cmd/${TASK_NAME}/main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
ENV TASK_NAME=modify-data-object
ENV ENTRY_CMD=/usr/local/bin/${TASK_NAME}

# install task binary
COPY --from=builder /${TASK_NAME} ${ENTRY_CMD}

ENTRYPOINT ["/usr/local/bin/${TASK_NAME}"]
CMD ["--help"]