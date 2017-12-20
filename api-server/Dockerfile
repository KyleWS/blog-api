FROM alpine
RUN apk add --no-cache ca-certificates
ADD api-server /api-server
EXPOSE 443
ENTRYPOINT ["/api-server"]
