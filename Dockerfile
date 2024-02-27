FROM golang:1.22.0
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o . ./cmd/main.go
EXPOSE 8080
CMD [ "/app/main" ]