FROM node:18-alpine

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# PROXY
#ENV HTTP_PROXY=http://172.16.0.250:4343
#ENV HTTPS_PROXY=http://172.16.0.250:4343

# Install dependencies
RUN npm install

# Copy test files
COPY . .

# PROXY
#ENV HTTP_PROXY=
#ENV HTTPS_PROXY=

# Set default command to run tests
ENTRYPOINT ["node", "run-tests.js"]

# Default arguments - can be overridden at runtime
CMD ["--environment", "docker"]
#CMD ["--environment", "local"]
