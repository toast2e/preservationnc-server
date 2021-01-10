FROM alpine

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY build/linux/preservationnc-server /app/preservationnc-server
EXPOSE 8080

CMD /app/preservationnc-server
