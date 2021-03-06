FROM golang:1.16-alpine

RUN mkdir /app

WORKDIR /app 

COPY . /app

RUN go mod download

RUN go get -u github.com/cosmtrek/air

EXPOSE 8080

ENTRYPOINT ["air","-c",".air.toml"]