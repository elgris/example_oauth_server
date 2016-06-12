FROM alpine:latest
RUN apk --update add tzdata ca-certificates make
COPY ./example_oauth_server /
EXPOSE 8000
ENTRYPOINT ["/example_oauth_server"]