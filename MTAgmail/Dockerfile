FROM golang:latest

RUN mkdir /app
RUN go get "github.com/gorilla/mux"
RUN go get "github.com/kjk/betterguid"
ADD . /app/
WORKDIR /app
EXPOSE 8002
CMD [ "go", "run", "/app/mtaGmail.go" ]
