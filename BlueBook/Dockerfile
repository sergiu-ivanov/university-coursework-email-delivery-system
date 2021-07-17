FROM golang:latest

RUN mkdir /app
RUN go get "github.com/gorilla/mux"
ADD . /app/
WORKDIR /app
EXPOSE 8008
CMD [ "go", "run", "/app/blueBook.go" ]
