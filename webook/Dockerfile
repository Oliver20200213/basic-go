#基础镜像
FROM ubuntu:20.04
#复制编译完成的包放到/app目录下  /app/webook app是目录 webook是文件名
COPY webook /app/webook
WORKDIR /app
#CMD是执行命令
#ENTRYPOINT最佳
ENTRYPOINT ["/app/webook"]