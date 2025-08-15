#!/usr/bin/env bash
IFS=$'\n\t'
set -euo pipefail

read -p "Enter the name of the systemd service [auth]: " service_name
read -p "Enter the working directory [/root/go/projects/auth]: " working_directory
read -p "Enter the exec start binary [/root/go/projects/auth/app]: " exec_start
read -p "Enter the user [root]: " user
read -p "Do you know create config files? [y/n]: " create_config_files

if [ "$create_config_files" == "y" ]; then
    mkdir -p $working_directory/configs
    sudo tee $working_directory/configs/config.env <<EOF
RABBIT_NAME=default
RABBIT_HOST=localhost
RABBIT_PASSWORD=default
RABBIT_PORT=5672
RABBIT_QUEUE_FORGOT_TOKEN=default
RABBIT_QUEUE_CONFIRM_EMAIL=default
LOG_PATH=./cmd/app
EMAIL_PROVIDER=default
SMTP_HOST=default
SMTP_PORT=465
SMTP_USERNAME=default
SMTP_PASSWORD=default
SMTP_FROM=default
MAIN_DOMAIN=default
EOF

    sudo tee $working_directory/configs/config.yaml <<EOF
env: "dev"
token_access_ttl: 60m
token_refresh_ttl: 720h
log_path: "/app/logs"
name_sso_service: "default"
rabbit:
    max_retries: 5
    retry_delay: 1000ms
grps: 
    port: 00000
    timeout: 10s
EOF
fi

read -p "Do you know create systemd service? [y/n]: " create_systemd_service
if [ "$create_systemd_service" != "y" ]; then
    exit 0
fi

#go build
echo "Building the project..."
if [ -d $working_directory ] && [ -f $working_directory/cmd/app/app.go ]; then
    if [ -f $exec_start ]; then
        echo "Exec start binary already exists"
    else
        go build -o $exec_start $working_directory/cmd/app/app.go
    fi
else
    echo "Not in the working directory or app.go not found"
    exit 1
fi

sudo tee /lib/systemd/system/$service_name.service <<EOF
[Unit]
Description=${service_name}
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=${user}
WorkingDirectory=${working_directory}
ExecStart=${exec_start}
Restart=on-failure
RestartSec=10s
StartLimitInterval=60s
StartLimitBurst=3
LimitNOFILE=65535
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=${service_name}

[Install]
WantedBy=multi-user.target
EOF

echo "Reloading systemd..."
sudo systemctl daemon-reload
sudo systemctl enable $service_name

echo "Checking the status of the service..."
sudo systemctl status $service_name

echo "Starting the service..."
sudo systemctl start $service_name

echo "Checking the status of the service..."
sudo systemctl status $service_name

echo "Stopping the service..."
sudo systemctl stop $service_name