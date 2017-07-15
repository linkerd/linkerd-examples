FROM jenkins:2.60.1

USER root
COPY kubectl /usr/bin/
RUN curl -O https://storage.googleapis.com/golang/go1.6.linux-amd64.tar.gz
RUN tar -xvf go1.6.linux-amd64.tar.gz
RUN mv go /usr/local

ENV GOPATH /go
RUN mkdir $GOPATH
ENV PATH /usr/local/go/bin:$GOPATH/bin:$PATH
RUN go get github.com/linkerd/namerctl && go install github.com/linkerd/namerctl
ENV NAMERCTL_BASE_URL http://namerd.default.svc.cluster.local:4180

USER jenkins
RUN /usr/local/bin/install-plugins.sh workflow-aggregator:2.4
RUN /usr/local/bin/install-plugins.sh github-oauth:0.24
COPY jobs /usr/share/jenkins/ref/jobs

ENTRYPOINT ["/bin/tini", "--", "/usr/local/bin/jenkins.sh"]
