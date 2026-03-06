FROM gcr.io/distroless/base-debian12
COPY /build/bin /pubsub-direct-push
CMD ["/pubsub-direct-push"]
