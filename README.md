# YouTube Video Downloader

A self-hosted YouTube video downloader built with Go, Gin, and FFmpeg. Download high-quality videos with audio by simply providing a YouTube URL.

## Prerequisites

- Docker
- Go 1.23+ and FFMpeg (if running locally)

## Built with
- Gin
- Youtube by kkdai
- Viper
- Zap
- FFmpeg
  
## Installation

### Using Makefile (Recommended)

```bash
# Clone the repository
git clone https://github.com/artyomkorchagin/yt-converter.git
cd yt-converter

# Build using Makefile
make build

# Run the contaier using Makefile
make up
```
### Using Docker
```bash
# Clone the repository
git clone https://github.com/artyomkorchagin/yt-converter.git
cd yt-converter

# Build the Docker image
docker build -t yt-converter .

# Run the container
docker run -d -p 3000:3000 --name yt-downloader yt-converter
```
### Local developement
```bash
# Install dependencies
go mod download

# Install FFmpeg
# Ubuntu/Debian: sudo apt install ffmpeg
# macOS: brew install ffmpeg
# Windows: Download from https://ffmpeg.org

# Run the application
go run cmd/main.go
```
## Usage 
- Open your browser and navigate to http://localhost:3000
- Paste a YouTube URL in the input field
- Click "Download"
- The video will be processed and downloaded automatically

## API Endpoint
You can also use the API directly:
```bash
# Download a video
curl http://localhost:3000/api/v1/video/VIDEO_ID -o video.mp4
```

### Configuration
Create a `config.env` file with your settings:
```bash
# Server configuration
LOG_LEVEL # DEV for develepoment, anything else for production
PORT=3000
HOST=localhost
```
