ARG VERSION
ARG CODENAME

FROM registry.drycc.cc/drycc/go-dev:latest AS build
ADD . /workspace
RUN export GO111MODULE=on \
  && cd /workspace \
  && sed -i "s#canary#$VERSION#g" version/version.go \
  && CGO_ENABLED=0 init-stack go build -o /usr/local/bin/drycc drycc.go

FROM registry.drycc.cc/drycc/base:$CODENAME

ENV DRYCC_UID=1001 \
  DRYCC_GID=1001 \
  DRYCC_HOME_DIR=/workspace

COPY --from=build /usr/local/bin/drycc /usr/local/bin/drycc

ADD docker/entrypoint.sh /entrypoint.sh

RUN install-packages bash-completion \
  && groupadd drycc --gid ${DRYCC_GID} \
  && useradd drycc -u ${DRYCC_UID} -g ${DRYCC_GID} -s /bin/bash -m -d ${DRYCC_HOME_DIR} \
  && mkdir -p /etc/bash_completion.d \
  && drycc completion bash > /etc/bash_completion.d/drycc \
  && chown ${DRYCC_GID}:${DRYCC_UID} /usr/local/bin \
  && echo 'alias cd="echo \"cd: restricted\"; false"' > ${DRYCC_HOME_DIR}/.bash_aliases \
  && echo 'alias unalias="echo \"unalias: restricted\"; false"' >> ${DRYCC_HOME_DIR}/.bash_aliases \
  && ln -s /usr/bin/ls /usr/local/bin/ls \
  && ln -s /usr/bin/cat /usr/local/bin/cat \
  && ln -s /usr/bin/sed /usr/local/bin/sed \
  && ln -s /usr/bin/find /usr/local/bin/find \
  && ln -s /usr/bin/bash /usr/local/bin/bash \
  && ln -s /usr/bin/mkdir /usr/local/bin/mkdir \
  && ln -s /usr/bin/sleep /usr/local/bin/sleep \
  && ln -s /usr/bin/dircolors /usr/local/bin/dircolors \
  && ln -s /usr/bin/tini /usr/local/bin/tini \
  && ln -s /usr/bin/init-stack /usr/local/bin/init-stack

ENTRYPOINT [ "/entrypoint.sh" ] 

ENV PATH=/usr/local/bin

USER ${DRYCC_UID}
WORKDIR ${DRYCC_HOME_DIR}
