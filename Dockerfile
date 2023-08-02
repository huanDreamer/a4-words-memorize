# 表示依赖 alpine 最新版
FROM alpine:latest
ENV VERSION 1.0

# 在容器根目录 创建一个 apps 目录
WORKDIR /apps

# 挂载容器目录
#VOLUME ["/apps"]

# 拷贝当前目录下 go_docker_demo1 可以执行文件
COPY words /apps/golang_app

# 拷贝配置文件到容器中
COPY web/* /apps/web/

# 设置时区为上海
#RUN cp /etc/localtime /etc/localtime
#RUN echo 'Asia/Shanghai' > /etc/timezone
RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

# 设置编码
ENV LANG C.UTF-8

# 暴露端口
#EXPOSE 8090

# 设置为 release 生产模式
ENV GIN_MODE=release

# 运行golang程序的命令
ENTRYPOINT ["/apps/golang_app"]