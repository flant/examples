FROM golang:1.18-alpine
WORKDIR /app
ADD . /app/
RUN go build -o main .
RUN chmod +x ./main
CMD ./main