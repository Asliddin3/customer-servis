FROM golang:1.18-alpine3.15
RUN mkdir customer
COPY . /customer-servis
WORKDIR /customer-servis
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8810