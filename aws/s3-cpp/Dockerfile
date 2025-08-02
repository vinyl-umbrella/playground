FROM ubuntu:22.04

RUN ln -sf /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    curl \
    libcurl4-openssl-dev \
    libssl-dev \
    uuid-dev \
    zlib1g-dev \
    libpulse-dev \
    awscli \
    && rm -rf /var/lib/apt/lists/*

RUN git clone --recurse-submodules --depth 1 https://github.com/aws/aws-sdk-cpp.git /tmp/aws-sdk-cpp && \
    cd /tmp/aws-sdk-cpp && \
    mkdir build && cd build && \
    cmake .. -DCMAKE_BUILD_TYPE=Release \
    -DBUILD_ONLY="s3;transfer" \
    -DCMAKE_PREFIX_PATH=/usr/local/ \
    -DCMAKE_INSTALL_PREFIX=/usr/local/ && \
    make -j$(nproc) && \
    make install && \
    rm -rf /tmp/aws-sdk-cpp

RUN useradd -m -u 1000 appuser
USER appuser
WORKDIR /app

COPY src/ ./src/
COPY CMakeLists.txt ./

RUN mkdir build && cd build && \
    cmake .. && \
    make

CMD ["./build/s3_stream_processor"]
