FROM golang:alpine as builder

RUN mkdir -p /go/src/github.com/moonkeat/chainstack

WORKDIR /go/src/github.com/moonkeat/chainstack

ADD . ./

RUN apk add --no-cache curl gcc build-base

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

RUN dep ensure -v

RUN go build -o main .

WORKDIR /go/src/github.com/moonkeat/chainstack/db

RUN go build -o goose .

WORKDIR /go/src/github.com/moonkeat/chainstack/scripts/create_user

RUN go build -o createuser .

FROM alpine

RUN apk add --no-cache postgresql-client

RUN adduser -S -D -H -h /app appuser

USER appuser

COPY --from=builder /go/src/github.com/moonkeat/chainstack/main /app/

COPY --from=builder /go/src/github.com/moonkeat/chainstack/db /app/db

COPY --from=builder /go/src/github.com/moonkeat/chainstack/scripts/create_user/createuser /app/

COPY --from=builder /go/src/github.com/moonkeat/chainstack/scripts/wait-for-postgres.sh /app/

WORKDIR /app

CMD ["./wait-for-postgres.sh", "postgres"]

