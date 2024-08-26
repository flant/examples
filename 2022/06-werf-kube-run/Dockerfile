FROM golang:1.18-alpine@sha256:77f25981bd57e60a510165f3be89c901aec90453fd0f1c5a45691f6cb1528807
WORKDIR /app
ADD . /app/
RUN go build -o main .
RUN chmod +x ./main
CMD ./main