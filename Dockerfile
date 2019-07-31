FROM golang:1.12 as build

COPY ./ /go/src/github.com/meyskens/irail-csv

WORKDIR /go/src/github.com/meyskens/irail-csv

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo .

FROM alpine

COPY --from=build /go/src/github.com/meyskens/irail-csv/irail-csv /usr/local/bin
RUN chmod +x /usr/local/bin/irail-csv

CMD irail-csv
