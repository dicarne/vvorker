FROM node:22.15 AS node-builder

USER root

LABEL maintainer dicarne@zhishudali.ink

WORKDIR /app

RUN npm config set registry https://registry.npmmirror.com/
RUN npm install -g pnpm
RUN npm i workerd -g

COPY . /app
WORKDIR /app


RUN cd /app/www && pnpm i && pnpm run prepareDev && pnpm run build && pnpm run export
RUN cd /app/ext/ai && pnpm i && pnpm run build
RUN cd /app/ext/kv && pnpm i && pnpm run build
RUN cd /app/ext/oss && pnpm i && pnpm run build
RUN cd /app/ext/pgsql && pnpm i && pnpm run build
RUN cd /app/ext/assets && pnpm i && pnpm run build
RUN cd /app/ext/task && pnpm i && pnpm run build
RUN cd /app/ext/control && pnpm i && pnpm run build


######################################################################################
FROM golang:1.24-alpine AS go-builder

COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs
WORKDIR /app

ENV GOPROXY https://goproxy.cn,direct
ENV PATH /usr/local/go/bin:$PATH
ENV GOROOT /usr/local/go

RUN go install github.com/cweill/gotests/gotests@latest 		&& \
	go install github.com/fatih/gomodifytags@latest     		&& \
	go install github.com/josharian/impl@latest             	&& \
	go install github.com/haya14busa/goplay/cmd/goplay@latest 	&& \
	go install github.com/go-delve/delve/cmd/dlv@latest     	&& \
	go install honnef.co/go/tools/cmd/staticcheck@latest    	&& \
	go install golang.org/x/tools/gopls@latest

COPY --from=node-builder /app /app
# 执行 go build 命令，-o 指定输出的二进制文件名称
RUN go mod tidy
RUN go build -o vvorker .

#######################################################################################

FROM ubuntu:22.04

# RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list && \
# 	sed -i s@/security.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list && \
# 	sed -i 's/ports.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list

# RUN sed -i 's#http://archive.ubuntu.com/#http://mirrors.tuna.tsinghua.edu.cn/#' /etc/apt/sources.list;

RUN rm /var/lib/dpkg/info/libc-bin.*
RUN apt-get clean
RUN apt-get update
RUN apt-get install libc-bin

RUN DEBIAN_FRONTEND="noninteractive" apt-get install -y \
	apt-transport-https \
	ca-certificates \
	fuse3 \
	sqlite3

COPY litefs.yml /etc/litefs.yml
COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

# 从 builder 阶段拷贝 workerd 到最终镜像的 /bin 目录
COPY --from=node-builder /usr/local/lib/node_modules/workerd/bin/workerd /bin/workerd

# 从 builder 阶段拷贝 go build 产物到最终镜像的 /bin 目录
COPY --from=go-builder /app/vvorker /app/vvorker

RUN apt-get update && \
	DEBIAN_FRONTEND=noninteractive apt-get install -qy libc++1 ca-certificates

RUN chmod +x /bin/*

WORKDIR /app

COPY .env.sample /app/.env

EXPOSE 8888
EXPOSE 8080

CMD [ "/app/vvorker" ]
