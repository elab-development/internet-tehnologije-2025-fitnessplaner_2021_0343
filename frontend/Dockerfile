FROM node:20-alpine

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm install

COPY . .

# Vite dev server port
EXPOSE 5173

# Default API URL inside Docker network (can be overridden)
ENV VITE_API_URL=http://backend:8080

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "5173"]

FROM node:20-alpine AS build

WORKDIR /app

# Install dependencies (including devDependencies for build)
COPY package.json package-lock.json ./
RUN npm install && npm list vite

# Copy the rest of the app and build
COPY . .
RUN npm run build

# Production image
FROM nginx:1.27-alpine

# Copy built assets from the build stage
COPY --from=build /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose default nginx port
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]


