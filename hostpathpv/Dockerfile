FROM alpine
MAINTAINER patrick

RUN apk add --no-cache openssh-server openssh-client bash coreutils rsync && \
  ssh-keygen -t rsa -P "" -f ~/.ssh/id_rsa && \
  ssh-keygen -t rsa -f /etc/ssh/ssh_host_rsa_key && \
  cp ~/.ssh/id_rsa.pub ~/.ssh/authorized_keys && \
  sed -i "s/#.*StrictHostKeyChecking.*/    StrictHostKeyChecking no/g" /etc/ssh/ssh_config

ADD entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh



