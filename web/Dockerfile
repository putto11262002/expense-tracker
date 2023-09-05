
FROM node:lts AS build

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

ARG API_BASE_URL
ENV VITE_API_BASE_URL=$API_BASE_URL

RUN npm run build

FROM nginx:alpine


COPY nginx.conf /etc/nginx/nginx.conf

COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]
