FROM golang:latest


WORKDIR /go/src/ozon_test
COPY . /go/src/ozon_test

#RUN go mod tidy
RUN go build -o ./bin/ozon_test ./cmd/ozon_test/
#RUN go build -o app
# Для возможности запуска скрипта
RUN chmod +x /go/src/ozon_test/scripts/*


CMD ["/go/src/ozon_test/bin/ozon_test"]