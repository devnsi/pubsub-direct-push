FROM gcr.io/distroless/base-debian12
COPY /build/bridge /bridge
CMD ["/bridge"]
