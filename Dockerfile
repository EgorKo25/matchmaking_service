FROM golang:1.23.0
LABEL authors="eko"

RUN mkdir app
COPY . app/
WORKDIR app/

RUN go mod download
RUN go build -o mm-service .

CMD ["sh", "-c", "./mm-service"]

EXPOSE 8000
