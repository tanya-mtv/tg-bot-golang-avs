
######################
#      Builder       #
######################

FROM golang:1.20 as builder

# Never ask for user input
ARG DEBIAN_FRONTEND=noninteractive

# ACCEPT_EULA=Y is required to install Microsoft ODBC Driver
ARG ACCEPT_EULA=Y

RUN apt-get update && \
    apt-get upgrade --yes && \
    apt-get install --yes \
        ca-certificates \
        openssl \
        curl \
        gnupg2 && \
    curl https://packages.microsoft.com/keys/microsoft.asc | tee /etc/apt/trusted.gpg.d/microsoft.asc && \
    gpg --dearmor < /etc/apt/trusted.gpg.d/microsoft.asc > /usr/share/keyrings/microsoft-prod.gpg && \
    curl https://packages.microsoft.com/config/debian/12/prod.list | tee /etc/apt/sources.list.d/mssql-release.list && \
    apt-get update && \
    apt-get install --yes \
        unixodbc \
        unixodbc-dev \
        # install Microsoft ODBC Driver for SQL Server
        msodbcsql18 \
        # optional: for bcp and sqlcmd
        mssql-tools18

# Set destination for COPY
WORKDIR /build

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
COPY cmd ./cmd/
COPY internal ./internal/

# RUN ls -alh /build/

# Build
RUN GOARCH=amd64 GOOS=linux go build -o /tg-avs-bot cmd/main.go

######################
#       Runner       #
######################

FROM debian:stable-slim as runner

ARG USERNAME=bot
ARG WORK_DIR=/opt/bot

# Never ask for user input
ARG DEBIAN_FRONTEND=noninteractive

# ACCEPT_EULA=Y is required to install Microsoft ODBC Driver
ARG ACCEPT_EULA=Y

ENV CONFIG_TYPE yaml
ENV CONFIG_PATH /usr/local/etc/bot.yaml

COPY --from=builder /tg-avs-bot /usr/local/bin/tg-avs-bot

RUN apt-get update && \
    apt-get upgrade --yes && \
    apt-get install --yes \
        ca-certificates \
        openssl \
        curl \
        gnupg2 \
        libglib2.0-0 \
        libnss3 \
        libxcb1 \
        libdbus-1-3 \
        libatk1.0-0 \
        libatk-bridge2.0-0 \
        libcups2 \
        libdrm2 \
        libxkbcommon0 \
        libxcomposite1 \
        libxdamage1 \
        libxfixes3 \
        libxrandr2 \
        libgbm1 \
        libpango-1.0-0 \
        libcairo2 \
        libasound2 \
        unzip \
        xvfb \
        libxi6 \
        libgconf-2-4 \
        fonts-liberation \
        fonts-liberation2 && \
    curl https://packages.microsoft.com/keys/microsoft.asc | tee /etc/apt/trusted.gpg.d/microsoft.asc && \
    gpg --dearmor < /etc/apt/trusted.gpg.d/microsoft.asc > /usr/share/keyrings/microsoft-prod.gpg && \
    curl https://packages.microsoft.com/config/debian/12/prod.list | tee /etc/apt/sources.list.d/mssql-release.list && \
    apt-get update && \
    apt-get install --yes \
        unixodbc \
        # install Microsoft ODBC Driver for SQL Server
        msodbcsql18 && \
    mkdir -p ${WORK_DIR} && \
    adduser --home ${WORK_DIR} --disabled-login --shell /bin/nologin ${USERNAME} && \
    chown ${USERNAME}:${USERNAME} ${WORK_DIR} && \
    cd ${WORK_DIR}/ && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR ${WORK_DIR}/

USER ${USERNAME}
CMD ["/usr/local/bin/tg-avs-bot"]
