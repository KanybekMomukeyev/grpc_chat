FROM resin/raspberrypi3-golang
#FROM google/golang

WORKDIR /gopath/src/github.com/KanybekMomukeyev/grpc_chat

ADD . /gopath/src/github.com/KanybekMomukeyev/
ADD client /gopath/src/github.com/KanybekMomukeyev/grpc_chat/client
ADD cmd /gopath/src/github.com/KanybekMomukeyev/grpc_chat/cmd
ADD database /gopath/src/github.com/KanybekMomukeyev/grpc_chat/database
ADD defaults /gopath/src/github.com/KanybekMomukeyev/grpc_chat/defaults
ADD proto /gopath/src/github.com/KanybekMomukeyev/grpc_chat/proto
ADD server /gopath/src/github.com/KanybekMomukeyev/grpc_chat/server
ADD test_main_samples /gopath/src/github.com/KanybekMomukeyev/grpc_chat/test_main_samples
ADD utils /gopath/src/github.com/KanybekMomukeyev/grpc_chat/utils

# go get all of the dependencies
RUN go get google.golang.org/grpc
RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx
RUN go get github.com/blevesearch/bleve
RUN go get github.com/KanybekMomukeyev/grpc_chat

ADD main.go /gopath/src/github.com/KanybekMomukeyev/grpc_chat/main.go

EXPOSE 8080

CMD ["go", "run", "/gopath/src/github.com/KanybekMomukeyev/grpc_chat/main.go serve 10000"]

#ENTRYPOINT ["/gopath/bin/testingpackages"]
#ENTRYPOINT /go/bin/streamtest