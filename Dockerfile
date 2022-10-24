FROM gcr.io/distroless/static:nonroot
WORKDIR /

EXPOSE 8932

COPY manager.linux manager
USER 65532:65532


ENTRYPOINT ["/manager"]
