FROM node:alpine

WORKDIR /opt/nodedir
ENV PATH /opt/nodedir/node_modules/.bin:$PATH

RUN npm install elm elm-live

WORKDIR /opt/ui

EXPOSE 8000