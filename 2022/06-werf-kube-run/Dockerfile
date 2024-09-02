FROM golang:1.23-alpine@sha256:d0b31558e6b3e4cc59f6011d79905835108c919143ebecc58f35965bf79948f4
WORKDIR /app
ADD . /app/
RUN go build -o main .
RUN chmod +x ./main
CMD ./main