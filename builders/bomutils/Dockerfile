FROM ubuntu

RUN apt-get update && apt-get install -y \
    build-essential \
    git && \
    git clone https://github.com/hogliux/bomutils.git && \
    cd bomutils && \
    make && \
    make install


COPY create_bom.sh /
