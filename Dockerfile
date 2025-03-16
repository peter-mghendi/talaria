FROM alpine:latest
WORKDIR /app

COPY talaria /app/talaria

EXPOSE 9999
ENTRYPOINT ["/app/talaria"]
