FROM golang:1.17
MAINTAINER zhijin_feng@berkeley.edu 
LABEL maintainer "zhijin_feng" 


# 在container里面创建文件夹;
WORKDIR /go/src/appstore 
# copy当前文件夹appstore的所有的文件到container上面;  
ADD . / /go/src/appstore/


RUN go get cloud.google.com/go/storage
RUN go get github.com/auth0/go-jwt-middleware
RUN go get github.com/form3tech-oss/jwt-go
RUN go get github.com/gorilla/handlers
RUN go get github.com/gorilla/mux
RUN go get github.com/olivere/elastic/v7
RUN go get github.com/stripe/stripe-go/v74
RUN go get github.com/pborman/uuid
RUN go get gopkg.in/yaml.v2

# 打开8080端口;
EXPOSE 8080
# 程序听取8080端口的requests;
CMD ["/usr/local/go/bin/go", "run", "main.go"]
