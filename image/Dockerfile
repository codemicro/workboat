FROM ubuntu:22.04

RUN mkdir /workboat
WORKDIR /workboat

ADD --chown=root:root setup.sh setup.sh
ADD --chown=root:root bootstrap.sh bootstrap.sh

RUN bash ./setup.sh

ENV PATH "$PATH:/usr/local/go/bin:/root/go/bin"

ENTRYPOINT ["bash", "bootstrap.sh"]