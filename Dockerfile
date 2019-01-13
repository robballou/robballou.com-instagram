FROM golang:1.9

# WORKDIR /go/src
# COPY paste_table.go .

# RUN go-wrapper download   # "go get -d -v ./..."
# RUN go-wrapper install    # "go install -v ./..."

ADD instagram.go /go/src/github.com/robballou/instagram-image-links-go/
RUN go get github.com/mmcdole/gofeed && \
  go get github.com/go-redis/redis && \
  go get github.com/PuerkitoBio/goquery
RUN go install github.com/robballou/instagram-image-links-go
ENTRYPOINT "/go/bin/instagram-image-links-go"
