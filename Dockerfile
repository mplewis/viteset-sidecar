FROM golang:1.16 AS build
WORKDIR /build

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd
COPY . .
RUN CGO_ENABLED=0 go build -o app

#######################################

FROM scratch AS prod

COPY --from=build /etc_passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER nobody
COPY --from=build /build/app /app

EXPOSE 80
CMD [ "/app" ]
