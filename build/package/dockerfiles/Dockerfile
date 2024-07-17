FROM alpine
RUN apk update && apk add --update --no-cache ca-certificates bash
COPY n2x-node /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/n2x-node"]
