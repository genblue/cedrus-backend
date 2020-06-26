FROM alpine:latest
USER root
ENV API_PORT 8002

RUN adduser -D sethservice && \
    mkdir -m 0755 /nix && chown sethservice /nix && \
    mkdir -p /etc/nix && echo 'sandbox = false' > /etc/nix/nix.conf && \
    apk --no-cache add ca-certificates curl git

USER sethservice
RUN curl https://nixos.org/nix/install | sh && source $HOME/.nix-profile/etc/profile.d/nix.sh && \
    echo "export PATH=$PATH:$HOME/.nix-profile/bin" >> $HOME/.profile && \
    source $HOME/.profile && \
    git clone --recursive https://github.com/dapphub/dapptools $HOME/.dapp/dapptools && \
    nix-env -f $HOME/.dapp/dapptools -iA seth

RUN source $HOME/.profile && seth help

EXPOSE $API_PORT

CMD ["/bin/sh"]
