# Stage 1: Build for Production
FROM node:23 as builder
WORKDIR /app
COPY ./frontend/package*.json ./
RUN npm install
COPY ./frontend ./
RUN npm run build

# Stage 2: Serve Production Build
FROM nginx:stable-alpine as production
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

# Stage 3: Development Environment
FROM node:23 as development
WORKDIR /app
COPY ./frontend/package*.json ./
RUN npm install
COPY ./frontend ./
EXPOSE 3000
CMD ["npm", "run", "dev"]