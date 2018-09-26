FROM golang:1.11.0

VOLUME /config

RUN mkdir -p /go/src/github.com/Adi146/ScryfallSpoilerBot
WORKDIR /go/src/github.com/Adi146/ScryfallSpoilerBot
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["ScryfallSpoilerBot", "--config=/config/config.yaml", "--spoiledCards=/config/spoiledCards.json", "--logFile=/config/scryfallSpoilerBot.log"]
