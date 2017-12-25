FROM ubuntu

RUN apt-get update && apt-get install -y \
    build-essential \
    git \
    autoconf \
    libxml2-dev \
    libcurl4-openssl-dev \
    python2.7 \
    libbz2-dev \
    liblzma-dev \
    libssl-dev && \
    git clone https://github.com/mackyle/xar.git && \
    cd xar/xar && \
    ./autogen.sh && \
    ./configure --with-bzip2 --with-lzma=/usr && \
    make && \
    make install

COPY create_xar.sh /
