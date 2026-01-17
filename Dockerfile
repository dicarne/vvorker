FROM node:22.15 AS dependencies

USER root

LABEL maintainer=dicarne@zhishudali.ink

WORKDIR /app

RUN npm config set registry https://registry.npmmirror.com/ && \
    npm install -g pnpm

COPY admin/package.json /app/admin/package.json
COPY admin/pnpm-lock.yaml /app/admin/pnpm-lock.yaml
RUN cd /app/admin && pnpm i

COPY ext/ai/package.json /app/ext/ai/package.json
COPY ext/ai/pnpm-lock.yaml /app/ext/ai/pnpm-lock.yaml
RUN cd /app/ext/ai && pnpm i

COPY ext/kv/package.json /app/ext/kv/package.json
COPY ext/kv/pnpm-lock.yaml /app/ext/kv/pnpm-lock.yaml
RUN cd /app/ext/kv && pnpm i

COPY ext/oss/package.json /app/ext/oss/package.json
COPY ext/oss/pnpm-lock.yaml /app/ext/oss/pnpm-lock.yaml
RUN cd /app/ext/oss && pnpm i

COPY ext/pgsql/package.json /app/ext/pgsql/package.json
COPY ext/pgsql/pnpm-lock.yaml /app/ext/pgsql/pnpm-lock.yaml
RUN cd /app/ext/pgsql && pnpm i

COPY ext/assets/package.json /app/ext/assets/package.json
COPY ext/assets/pnpm-lock.yaml /app/ext/assets/pnpm-lock.yaml
RUN cd /app/ext/assets && pnpm i

COPY ext/task/package.json /app/ext/task/package.json
COPY ext/task/pnpm-lock.yaml /app/ext/task/pnpm-lock.yaml
RUN cd /app/ext/task && pnpm i

COPY ext/control/package.json /app/ext/control/package.json
COPY ext/control/pnpm-lock.yaml /app/ext/control/pnpm-lock.yaml
RUN cd /app/ext/control && pnpm i

COPY ext/mysql/package.json /app/ext/mysql/package.json
COPY ext/mysql/pnpm-lock.yaml /app/ext/mysql/pnpm-lock.yaml
RUN cd /app/ext/mysql && pnpm i


FROM dependencies AS node-builder

RUN npm i workerd@v1.20250619.0 -g

WORKDIR /app

COPY admin /app/admin
COPY ext /app/ext

RUN cd /app/admin && pnpm run build
RUN cd /app/ext/ai && pnpm run build
RUN cd /app/ext/kv && pnpm run build
RUN cd /app/ext/oss && pnpm run build
RUN cd /app/ext/pgsql && pnpm run build
RUN cd /app/ext/assets && pnpm run build
RUN cd /app/ext/task && pnpm run build
RUN cd /app/ext/control && pnpm run build
RUN cd /app/ext/mysql && pnpm run build

######################################################################################
FROM golang:1.25-alpine AS go-builder

# COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct
ENV PATH=/usr/local/go/bin:$PATH
ENV GOROOT=/usr/local/go

# RUN go install github.com/cweill/gotests/gotests@latest 		&& \
# 	go install github.com/fatih/gomodifytags@latest     		&& \
# 	go install github.com/josharian/impl@latest             	&& \
# 	go install github.com/haya14busa/goplay/cmd/goplay@latest 	&& \
# 	go install github.com/go-delve/delve/cmd/dlv@latest     	&& \
# 	go install honnef.co/go/tools/cmd/staticcheck@latest    	&& \
# 	go install golang.org/x/tools/gopls@latest

COPY . /app

COPY --from=node-builder /app/admin/dist /app/admin/dist
COPY --from=node-builder /app/ext/ai/dist /app/ext/ai/dist
COPY --from=node-builder /app/ext/kv/dist /app/ext/kv/dist
COPY --from=node-builder /app/ext/oss/dist /app/ext/oss/dist
COPY --from=node-builder /app/ext/pgsql/dist /app/ext/pgsql/dist
COPY --from=node-builder /app/ext/assets/dist /app/ext/assets/dist
COPY --from=node-builder /app/ext/task/dist /app/ext/task/dist
COPY --from=node-builder /app/ext/control/dist /app/ext/control/dist
COPY --from=node-builder /app/ext/mysql/dist /app/ext/mysql/dist
# 执行 go build 命令，-o 指定输出的二进制文件名称
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=cache,target=/go/pkg,sharing=locked \
    go mod tidy
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=cache,target=/go/pkg,sharing=locked \
    go build -o vvorker .

#######################################################################################

FROM ubuntu:24.04

RUN sed -i s@/archive.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list && \
	sed -i s@/security.ubuntu.com/@/mirrors.aliyun.com/@g /etc/apt/sources.list && \
	sed -i 's/ports.ubuntu.com/mirrors.aliyun.com/g' /etc/apt/sources.list && \
	sed -i 's#http://archive.ubuntu.com/#http://mirrors.tuna.tsinghua.edu.cn/#' /etc/apt/sources.list;

# RUN rm /var/lib/dpkg/info/libc-bin.*
RUN apt-get clean && apt-get update && DEBIAN_FRONTEND="noninteractive" apt-get install -y \
	apt-transport-https \
	ca-certificates \
	fuse3 \
	sqlite3 \
	curl \
	vim

# COPY litefs.yml /etc/litefs.yml
# COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs

# RUN apt-get update && \
# 	DEBIAN_FRONTEND=noninteractive apt-get install -qy libc++1 ca-certificates

# 从 builder 阶段拷贝 workerd 到最终镜像的 /bin 目录
COPY --from=node-builder /usr/local/lib/node_modules/workerd/bin/workerd /bin/workerd

# 从 builder 阶段拷贝 go build 产物到最终镜像的 /bin 目录
COPY --from=go-builder /app/vvorker /app/vvorker

RUN chmod +x /bin/*

WORKDIR /app

# COPY .env.sample /app/.env

# 控制台
EXPOSE 8888
# 纯服务
EXPOSE 8080
# frp
EXPOSE 18080

CMD [ "/app/vvorker" ]
