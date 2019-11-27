FROM golang

COPY . .

RUN go build -o patchbay-simple-server

ENTRYPOINT ["./patchbay-simple-server"]
