FROM bitnami/minideb:stretch

RUN install_packages \
  build-essential \
  ca-certificates \
  curl \
  git \
  procps \
  sudo \
  ;

RUN addgroup --gid 1000 docker \
  && adduser --uid 1000 --ingroup docker --home /home/docker --shell /bin/sh --disabled-password --gecos "" docker \
  && echo '%docker ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN curl -sSL https://github.com/boxboat/fixuid/releases/download/v0.3/fixuid-0.3-linux-amd64.tar.gz | tar -C /usr/local/bin -xz \
  && sudo chown root:root /usr/local/bin/fixuid \
  && sudo chmod 4755 /usr/local/bin/fixuid \
  && sudo mkdir -p /etc/fixuid \
  && echo "user: docker\ngroup: docker\n" > /etc/fixuid/config.yml

USER docker:docker
ENV HOME /home/docker
WORKDIR /home/docker

ARG GO_IMPORT_PATH
ENV GOPATH=/gopath
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH
WORKDIR /gopath/src/${GO_IMPORT_PATH}
RUN sudo chown -R docker:docker /gopath

COPY .go-version ./
RUN curl -sSL "https://dl.google.com/go/go$(cat .go-version).linux-amd64.tar.gz" | sudo tar -C /usr/local -xz
