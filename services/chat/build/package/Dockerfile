FROM alpine
ADD ./service /srv/
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
ENTRYPOINT ["/srv/service", "serve"]
EXPOSE 8080
