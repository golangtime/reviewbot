FROM golang:1.22

COPY . .

RUN go mod download -x

RUN go build main.go -o bot

RUN cp bot /usr/bin/bot

CMD ["bot"]

