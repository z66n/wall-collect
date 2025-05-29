# Wall-Collect 🖼️

A beautiful wallpaper collection manager built with Go. Upload, organize, view, and download your wallpaper collection with ease.

![Screenshot](https://github.com/user-attachments/assets/f5b1b6e1-c5ef-4d65-b2f8-6f40cdeaca2b)

## Features ✨

- **Batch Upload**: Upload multiple wallpapers at once
- **Responsive Gallery**: Beautiful grid display of your collection
- **Easy Management**: Rename, delete, or download individual wallpapers
- **Fast Performance**: Lightweight Go backend serves images quickly
- **Simple Interface**: Clean, intuitive UI
- **Basic Authentication**: Enabled by default (admin/password)

## Installation 💻

### Prerequisites
- Docker installed
- (Optional) Go 1.21+ for local development
- Proper file permissions for your wallpaper directory

### Important Notes
1. **Permissions**: Ensure the user UID/GID (typically 1000:1000) has read/write access to:
   - Your host's wallpaper directory (`/path/to/wallpapers`)

2. **Security**: The default credentials are:
   ```
   Username: admin
   Password: password
   ```
   **Change these immediately in production!**

### Quick Start with Docker
Create a `.env` file:
```
AUTH_USERNAME=admin
AUTH_PASSWORD=secure_password_here
UPLOAD_DIR=/app/uploads
PORT=8080
```

Then run with:
```bash
docker run -d \
  -u 1000:1000 \
  -p 8080:8080 \
  -v /path/to/wallpapers:/app/uploads \
  --env-file .env \
  --name wall-collect \
  --restart unless-stopped \
  ghcr.io/z66n/wall-collect
```
Access at: `http://localhost:8080`

### Docker Compose
1. Create `docker-compose.yml`:
```yaml
services:
  wall-collect:
    image: ghcr.io/z66n/wall-collect
    container_name: wall-collect
    user: "1000:1000"
    ports:
        - 8080:8080
    volumes:
        - /path/to/wallpapers:/app/uploads
    env_file:
        - .env
    restart: unless-stopped
```

2. Start the service:
```bash
docker-compose up -d
```

## Configuration ⚙️

| Variable         | Default     | Description                          | Required |
|------------------|-------------|--------------------------------------|----------|
| `UPLOAD_DIR`     | `/app/uploads` | Wallpaper storage directory       | No       |
| `PORT`           | `8080`      | Port to listen on                    | No       |
| `AUTH_USERNAME`  | `admin`     | Basic authentication username        | No       |
| `AUTH_PASSWORD`  | `password`  | Basic authentication password        | No       |

## Development
```bash
# Clone and run with defaults
git clone https://github.com/z66n/wall-collect.git
cd wall-collect
go run wally.go -addr localhost:8080

# Or build and run:
go build -o wally
./wally -addr localhost:8080
```
