FROM ubuntu:22.04

USER root

LABEL maintainer me@vaala.cat

RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list && \
	sed -i s@/security.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list && \
	sed -i 's/ports.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list

RUN apt-get update && \
	DEBIAN_FRONTEND=noninteractive apt-get install -qy libc++1

RUN apt-get update && DEBIAN_FRONTEND="noninteractive" apt-get install -y\
	apt-transport-https \
	ca-certificates \
	curl \
	gnupg \
	zsh \
	fish \
	lsb-release \
	wget \
	tmux git \
	build-essential \
	sudo \
	rsync \
	ssh \
	vim \
	unzip \
	p7zip-full \
	bash \
	inetutils-ping \
	net-tools \
	pgcli \
	htop \
	locales \
	man \
	python3 \
	python3-pip \
	software-properties-common \
	systemd \
	systemd-sysv \
	fuse3 \
	sqlite3 \
	--no-install-recommends 

RUN wget http://s3.cloud.zhishudali.ink/public/golang/go1.24.3.linux-amd64.tar.gz && \
	rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.3.linux-amd64.tar.gz

ENV NODE_VERSION=22.15.0
RUN curl -o- https://gh-proxy.com/raw.githubusercontent.com/nvm-sh/nvm/v0.39.0/install.sh | bash
ENV NVM_DIR=/root/.nvm
RUN . "$NVM_DIR/nvm.sh" && nvm install ${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm use v${NODE_VERSION}
RUN . "$NVM_DIR/nvm.sh" && nvm alias default v${NODE_VERSION}
ENV PATH="/root/.nvm/versions/node/v${NODE_VERSION}/bin/:${PATH}"

WORKDIR /app

RUN pip config set global.index-url http://pypi.douban.com/simple/ && \
	pip config set install.trusted-host pypi.douban.com
	

COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

ENV GOPROXY https://proxy.golang.com.cn,direct
ENV PATH /usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go

RUN go install github.com/cweill/gotests/gotests@latest 		&& \
	go install github.com/fatih/gomodifytags@latest     		&& \
	go install github.com/josharian/impl@latest             	&& \
	go install github.com/haya14busa/goplay/cmd/goplay@latest 	&& \
	go install github.com/go-delve/delve/cmd/dlv@latest     	&& \
	go install honnef.co/go/tools/cmd/staticcheck@latest    	&& \
	go install golang.org/x/tools/gopls@latest



RUN npm config set registry https://registry.npmmirror.com/
RUN npm install -g pnpm
RUN npm i workerd -g
# RUN cp /root/.nvm/versions/node/v22.15.0/bin/workerd /workspaces/vorker/bin/workerd

EXPOSE 8888
EXPOSE 8080

CMD [ "/bin/bash" ]