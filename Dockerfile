# 表示依赖 alpine 最新版
FROM alpine:latest
ENV VERSION 1.0

# 在容器根目录 创建一个 apps 目录
WORKDIR /apps

# 拷贝可执行文件
COPY words /apps/golang_app

# 拷贝前端页面与静态资源（注意：web/* 不含子目录，static/ 需单独拷贝）
COPY web/*.html /apps/web/
COPY web/static /apps/web/static

# 数据目录（本地文件存储），建议挂载宿主机目录持久化
RUN mkdir -p /apps/data
VOLUME ["/apps/data"]
ENV WORDS_DATA_DIR=/apps/data

# 设置时区为上海
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

# 设置编码
ENV LANG C.UTF-8

# 暴露端口
EXPOSE 8900

# 设置为 release 生产模式
ENV GIN_MODE=release

# 运行golang程序的命令
ENTRYPOINT ["/apps/golang_app"]
