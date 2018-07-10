FROM node:9
MAINTAINER dryajov

# setup app dir
RUN mkdir -p /webrtc-star/
WORKDIR /webrtc-star/

run npm init -y
RUN npm install libp2p-webrtc-star

# start server
CMD npx star-signal --port=9090 --host=127.0.0.1

# expose server
EXPOSE 9090
