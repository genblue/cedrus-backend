FROM ubuntu:18.04 as buildBinaries

RUN adduser --disabled-password --gecos "" sethservice && \
    apt-get update -y && \
    apt-get install -y ca-certificates sed nodejs curl git bash make perl bc coreutils cmake libcurl4-openssl-dev libelf-dev libdw-dev binutils-dev libjansson-dev

# Build ethabi
RUN curl https://sh.rustup.rs -sSf | sh -s -- -y && \
    export PATH=$PATH:$HOME/.cargo/bin && \
    git clone --branch v7.0.0 https://github.com/paritytech/ethabi /home/sethservice/ethabi && \
    cd /home/sethservice/ethabi && \
    cargo install ethabi-cli && \
    cp $HOME/.cargo/bin/* /usr/local/bin/

# Build ethsign
RUN git clone --branch seth/0.8.3 --recursive https://github.com/dapphub/dapptools /home/sethservice/dapptools && \
    chmod -R +rwx /home/sethservice/dapptools/ && \
    cd /home/sethservice && \
    curl -O https://dl.google.com/go/go1.13.3.linux-amd64.tar.gz && \
    tar -xvf go1.13.3.linux-amd64.tar.gz && \
    mv go /usr/local
COPY ./vendor/ethsign/go.mod /home/sethservice/dapptools/src/ethsign/go.mod
RUN cd /home/sethservice/dapptools/src/ethsign && \
    export PATH=$PATH:/usr/local/go/bin && \
    go get && \
    go build  -o /usr/local/bin/ethsign ethsign.go

RUN git clone --branch 20131105 https://github.com/keenerd/jshon /home/sethservice/jshon && \
    cd /home/sethservice/jshon && \
    make && \
    cp jshon /usr/local/bin/jshon

# Give sethservice rights
RUN chown -R sethservice /usr/local/bin/

USER sethservice
RUN chmod +x /usr/local/bin/ethsign

FROM ubuntu:18.04

ENV API_PORT 8000
ENV DB_URI mongodb://localhost:27017

RUN adduser --disabled-password --gecos "" sethservice && \
    apt-get update -y && \
    apt-get install -y ca-certificates sed nodejs curl bash perl bc coreutils make git libjansson-dev && \
    mkdir -p /home/sethservice/secrets/ethereum && \
    chown -R sethservice /home/sethservice/secrets/ethereum && \
    git clone --branch seth/0.8.3 --recursive https://github.com/dapphub/dapptools /home/sethservice/dapptools && \
    chmod -R +rwx /home/sethservice/dapptools/ && \
    chown -R sethservice /home/sethservice/dapptools/src/seth/ && \
    cd /home/sethservice/dapptools/src/seth && \
    cp Makefile MakefileTemp && \
    sed 's/-n//g' MakefileTemp > Makefile && \
    make install

# Copy binaries
COPY --from=buildBinaries /usr/local/bin/ethsign /usr/local/bin/ethsign
COPY --from=buildBinaries /usr/local/bin/ethabi /usr/local/bin/ethabi
COPY --from=buildBinaries /usr/local/bin/jshon /usr/local/bin/jshon

# Copy (test or prod) keystore
COPY ./test/keystore_wallet_test.json /home/sethservice/secrets/ethereum/keystore.json
COPY ./test/password /home/sethservice/secrets/ethereum/password

RUN echo "/usr/local/lib" > /etc/ld.so.conf.d/usr_local.conf && /sbin/ldconfig

RUN chown -R sethservice /usr/local/bin/ && \
    chown -R sethservice /home/sethservice/secrets/ethereum

USER sethservice
RUN chmod +x /usr/local/bin/* && seth ls

WORKDIR /home/sethservice

COPY ./bin/sethservice .

EXPOSE $API_PORT

CMD ["./sethservice"]
