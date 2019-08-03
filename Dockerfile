FROM golang:1.12 as build

ENV TZ=Europe/Brussels
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY ./ /go/src/github.com/meyskens/irail-csv

WORKDIR /go/src/github.com/meyskens/irail-csv

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo .

FROM alpine

RUN apk add --no-cache ca-certificates

COPY --from=build /go/src/github.com/meyskens/irail-csv/irail-csv /usr/local/bin
RUN chmod +x /usr/local/bin/irail-csv

CMD irail-csv
