FROM centos:latest

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD epel-7.repo /etc/yum.repos.d


#安装YUM源
RUN yum -y update && yum -y install epel-release && yum -y install redis && yum -y install golang

#修改绑定IP地址
RUN sed -i -e 's@bind 127.0.0.1@bind 0.0.0.0@g' /etc/redis.conf
#关闭保护模式
RUN sed -i -e 's@protected-mode yes@protected-mode no@g' /etc/redis.conf

#ENV GOPROXY https://goproxy.io
ENV GO111MODULE on

WORKDIR /app
ADD . /app

#RUN go mod download
#RUN go build .
RUN go build -mod=vendor

EXPOSE 80 443

RUN chmod +x ./run.sh
# ENTRYPOINT  ["./go-short"]
CMD ["./run.sh"]
