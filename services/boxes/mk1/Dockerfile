FROM node:7
MAINTAINER kumavis

# setup app dir
RUN mkdir -p /www/
WORKDIR /www/

# install dependencies
COPY ./package.json /www/package.json
RUN npm install

# copy over app
COPY ./index.js /www/

# start server
CMD node ./index.js
