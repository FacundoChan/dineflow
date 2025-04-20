FROM node:23-alpine3.20 as builder

WORKDIR /app
COPY . .
RUN npm install && npm run build

FROM nginx:1.27.5-alpine
COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
