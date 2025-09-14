
FROM node:18-alpine


WORKDIR /app

COPY package.json ./
COPY pnpm-lock.yaml ./
COPY postcss.config.mjs ./
COPY tailwind.config.js ./
COPY .postcssrc.json ./
COPY tsconfig*.json ./
COPY vite.config.ts ./


COPY package.json  ./

RUN npm install 


COPY . .


EXPOSE 5173

ENV DEBUG=tailwindcss:*


CMD ["npm", "run", "dev", "--", "--host"]