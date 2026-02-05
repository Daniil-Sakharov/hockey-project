#!/bin/bash
# Script to obtain initial SSL certificate from Let's Encrypt
# Run this ONCE on the server before the first deploy

set -e

DOMAIN=${1:-rinkstar.ru}
EMAIL=${2:-sakharov7404@gmail.com}

echo "=== SSL Certificate Setup for $DOMAIN ==="

# Create directories
echo "Creating directories..."
mkdir -p /opt/hockey/data/certbot/{conf,www}
mkdir -p /opt/hockey/config/nginx
mkdir -p /opt/hockey/frontend

# Create temporary nginx config for ACME challenge
echo "Creating temporary nginx config..."
cat > /opt/hockey/config/nginx/default.conf << 'EOF'
server {
    listen 80;
    server_name _;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 200 'SSL setup in progress...';
        add_header Content-Type text/plain;
    }
}
EOF

# Create placeholder index.html
echo "RinkStar - Setup in progress" > /opt/hockey/frontend/index.html

# Start temporary nginx container for ACME challenge
echo "Starting temporary nginx..."
docker run -d --name temp-nginx \
    -p 80:80 \
    -v /opt/hockey/config/nginx:/etc/nginx/conf.d:ro \
    -v /opt/hockey/data/certbot/www:/var/www/certbot:ro \
    nginx:alpine

echo "Waiting for nginx to start..."
sleep 3

# Request certificate from Let's Encrypt
echo "Requesting SSL certificate..."
docker run --rm \
    -v /opt/hockey/data/certbot/conf:/etc/letsencrypt \
    -v /opt/hockey/data/certbot/www:/var/www/certbot \
    certbot/certbot certonly \
    --webroot -w /var/www/certbot \
    -d "$DOMAIN" \
    --email "$EMAIL" \
    --agree-tos \
    --no-eff-email \
    --non-interactive

# Stop and remove temporary nginx
echo "Stopping temporary nginx..."
docker stop temp-nginx && docker rm temp-nginx

# Update nginx config with SSL
echo "Updating nginx config for HTTPS..."
cat > /opt/hockey/config/nginx/default.conf << EOF
server {
    listen 80;
    server_name _;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name _;

    ssl_certificate /etc/letsencrypt/live/$DOMAIN/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$DOMAIN/privkey.pem;

    # SSL settings
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Frontend (SPA)
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files \$uri \$uri/ /index.html;
    }

    # API proxy
    location /api/ {
        proxy_pass http://api:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Static assets caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff2?|ttf|eot)\$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/json application/xml;
}
EOF

echo ""
echo "=== SSL Certificate obtained successfully! ==="
echo ""
echo "Next steps:"
echo "1. Set up GitHub Secrets (SERVER_HOST, SERVER_USER, SERVER_SSH_KEY, GHCR_TOKEN, POSTGRES_PASSWORD)"
echo "2. Create .env file at /opt/hockey/deploy/compose/production/.env"
echo "3. Push to main branch to trigger deploy"
echo ""
echo "Certificate location: /opt/hockey/data/certbot/conf/live/$DOMAIN/"
