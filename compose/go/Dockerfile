FROM golang:1.19.1-alpine

WORKDIR /go/src/app
COPY . .

RUN apk add build-base &&\
    go get -d -v ./... &&\
    go install -v ./...

# RUN addgroup --system todd &&\
#     adduser --system --ingroup todd todd &&\
#     chown todd:todd /go/src/app

# USER todd

CMD ["main"]