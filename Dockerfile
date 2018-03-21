FROM golang:alpine

#RUN apt-get update && apt-get install -y ca-certificates git-core ssh
RUN apk add --update git openssh ca-certificates bzr && rm -rf /var/cache/apk/*

ADD keys/id_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
RUN git config --global http.https://gopkg.in.followRedirects true
RUN ssh-keyscan github.com >> ~/.ssh/known_hosts

COPY keys/sr.json /etc/
COPY raw_view.sql /etc/
COPY agg_view.sql /etc/



#Create folder
RUN mkdir -p $GOPATH/src/github.com/streamrail/views
COPY ./ $GOPATH/src/github.com/streamrail/views

# Get deps
RUN cd $GOPATH/src/github.com/streamrail/views && go get
RUN cd $GOPATH/src/github.com/streamrail/views && go build

CMD $GOPATH/src/github.com/streamrail/views/views