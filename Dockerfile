FROM golang:1.16 as build

COPY ./ /go/src/github.com/meyskens/irail-csv

WORKDIR /go/src/github.com/meyskens/irail-csv

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo .

FROM alpine

RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Europe/Brussels
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

COPY --from=build /go/src/github.com/meyskens/irail-csv/irail-csv /usr/local/bin
RUN chmod +x /usr/local/bin/irail-csv

CMD irail-csv
