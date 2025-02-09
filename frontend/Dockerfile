FROM node:16-alpine as build-stage
WORKDIR /app
COPY package*.json ./
RUN npm config set legacy-peer-deps=true --location=project
RUN npm install
COPY ./ .
# ARG VUE_APP_API_URL
# ENV VUE_APP_API_URL $VUE_APP_API_URL
RUN npm run build

FROM nginx as production-stage
RUN rm -rf /usr/share/nginx/html/*
COPY --from=build-stage /app/nginx.conf /etc/nginx/conf.d/default.conf:ro
COPY --from=build-stage /app/dist /usr/share/nginx/html
