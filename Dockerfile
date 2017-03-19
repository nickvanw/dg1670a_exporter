FROM alpine:latest
COPY cmd/dg1670a_exporter/dg1670a_exporter / 
EXPOSE 9191
ENTRYPOINT ["/dg1670a_exporter"]
