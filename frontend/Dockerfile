FROM node:22-alpine

WORKDIR /app


# Copy dependency manifests and install cleanly inside image
COPY package*.json ./
RUN npm install


# Copy application source code into image
COPY . .
# Ensure files belong to node user
RUN chown -R node:node /app

# Switch to non-root user
USER node

# RUN mkdir -p /app/node_modules && chown -R node:node /app

# USER node

# COPY --chown=node:node package*.json ./

# RUN npm install

# COPY --chown=node:node . .

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host"]