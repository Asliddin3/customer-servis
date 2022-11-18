FROM golang:1.18-alpine3.15
RUN mkdir customer
COPY . /customer-servis
WORKDIR /customer-servis
RUN go mod tidy
RUN export AWS_ACCESS_KEY_ID=AKIAQJAX63K7GHK6EVZP \
    export AWS_SECRET_ACCESS_KEY=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

RUN go build -o main cmd/main.go
CMD ./main
EXPOSE 8810