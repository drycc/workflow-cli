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

ADD docker/rootfs /

RUN install-packages bash-completion git vim \
  && groupadd drycc --gid ${DRYCC_GID} \
  && useradd drycc -u ${DRYCC_UID} -g ${DRYCC_GID} -s /bin/bash -m -d ${DRYCC_HOME_DIR} \
  && mkdir -p /etc/wait \
  && chown ${DRYCC_GID}:${DRYCC_UID} /etc/wait /usr/local/bin \
  && echo 'alias cd="echo \"cd: restricted\"; false"' > ${DRYCC_HOME_DIR}/.bash_aliases \
  && echo 'alias unalias="echo \"unalias: restricted\"; false"' >> ${DRYCC_HOME_DIR}/.bash_aliases \
  && ln /usr/bin/{ls,git,date,vim,cat,sed,clear,find,flock,bash,mkdir,sleep,dircolors} /usr/local/bin/ \
  && apt autoremove python3 -yq \
  && rm -rf /var/lib/apt/lists/* /sbin /usr/sbin /usr/bin

ENTRYPOINT [ "/entrypoint.sh" ] 

ENV PATH=/usr/local/bin

USER ${DRYCC_UID}
WORKDIR ${DRYCC_HOME_DIR}
