# Development image: runs the Vite dev server. Source is mounted by docker compose.
FROM node:24-alpine
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
EXPOSE 5174
CMD ["npm", "run", "dev"]
