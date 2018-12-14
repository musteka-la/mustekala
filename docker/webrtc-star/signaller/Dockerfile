FROM node:9
MAINTAINER dryajov

# setup app dir
RUN mkdir -p /webrtc-star/
WORKDIR /webrtc-star/

run npm init -y
# RUN npm install libp2p-webrtc-star
RUN npm install dryajov/js-libp2p-webrtc-star.git#master

# start server
CMD npx star-signal --port=9090 --host=0.0.0.0

# expose server
EXPOSE 9090
