# Wall-Collect 🖼️

A beautiful wallpaper collection manager built with Go. Upload, organize, view, and download your wallpaper collection with ease.

![Screenshot](https://github.com/user-attachments/assets/f5b1b6e1-c5ef-4d65-b2f8-6f40cdeaca2b)

## Features ✨

- **Batch Upload**: Upload multiple wallpapers at once
- **Responsive Gallery**: Beautiful grid display of your collection
- **Easy Management**: Rename, delete, or download individual wallpapers
- **Fast Performance**: Lightweight Go backend serves images quickly
- **Simple Interface**: Clean, intuitive UI

## Installation 💻

### Prerequisites
- Docker installed
- (Optional) Go 1.21+ for local development

### Quick Start with Docker
```bash
# Pull the latest image
docker pull ghcr.io/z66n/wall-collect

# Run with persistent storage
docker run -d \
  -p 8080:8080 \
  -v ./wallpapers:/app/uploads \
  -e UPLOAD_DIR=/app/uploads \
  --name wall-collect \
  ghcr.io/z66n/wall-collect
```
Access at: `http://localhost:8080`

### Docker Compose
Create `docker-compose.yml`:
```yaml
version: '3.8'

services:
  wall-collect:
    image: ghcr.io/z66n/wall-collect
    container_name: wall-collect
    ports:
      - "8080:8080"
    volumes:
      - ./wallpapers:/app/uploads
    environment:
      UPLOAD_DIR: /app/uploads
      PORT: "8080"
	restart: unless-stopped
```

Start with:
```bash
docker-compose up -d
```

### Local Development
1. Clone the repository:
```bash
git clone https://github.com/z66n/wall-collect.git
cd wall-collect
```

2. Build and run:
```bash
go build -o wally
./wally -addr localhost:8080
```

## Configuration ⚙️

| Variable       | Default     | Description                          |
|----------------|-------------|--------------------------------------|
| `UPLOAD_DIR`   | `./uploads` | Wallpaper storage directory          |
| `PORT`         | `8080`      | Port to listen on                    |
