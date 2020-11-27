FROM balenalib/raspberry-pi-debian-golang:stretch-build as build

RUN install_packages \
      curl \
      ca-certificates \
      git \
      autoconf \
      automake \
      g++ \
      protobuf-compiler \
      zlib1g-dev \
      libncurses5-dev \
      libssl-dev \
      pkg-config \
      libprotobuf-dev \
      make \
      go-bindata

ENV GOPATH=/go-home
ENV BASE=$GOPATH/src/browsh/interfacer
WORKDIR $BASE
ADD interfacer $BASE

# Build Browsh
RUN $BASE/contrib/build_browsh.sh


###########################
# Actual final Docker image
###########################
FROM balenalib/raspberry-pi-debian:stretch

ENV HOME=/app
WORKDIR /app

COPY --from=build /go-home/src/browsh/interfacer/browsh /app/browsh

RUN install_packages \
      xvfb \
      libgtk-3-0 \
      curl \
      ca-certificates \
      bzip2 \
      libdbus-glib-1-2 \
      procps \
      firefox-esr

# Block ads, etc. This includes porn just because this image is also used on the
# public SSH demo: `ssh brow.sh`.
RUN curl -o /etc/hosts https://raw.githubusercontent.com/StevenBlack/hosts/master/alternates/fakenews-gambling-porn-social/hosts

# Don't use root
RUN useradd -m user --home /app
RUN chown user:user /app
USER user

# Firefox behaves quite differently to normal on its first run, so by getting
# that over and done with here when there's no user to be dissapointed means
# that all future runs will be consistent.
RUN TERM=xterm script \
  --return \
  -c "/app/browsh" \
  /dev/null \
  >/dev/null & \
  sleep 10

CMD ["balena-idle"]
