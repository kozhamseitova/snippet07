FROM golang:1.15.6
ADD . /go/src/
WORKDIR /app
COPY . /app
RUN go build -o app cmd/web/*
EXPOSE 4000
ENTRYPOINT /app/app