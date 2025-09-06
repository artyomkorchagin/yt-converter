include config.env

build:
	docker build -t yt-converter .

up:
	docker run -d -p 3000:3000 --name yt-converter-app yt-converter

down:
	docker stop yt-converter-app || true
	docker rm yt-converter-app || true

restart: down up
	
clean:
	docker stop yt-converter-app || true
	docker rm yt-converter-app || true
	docker rmi yt-converter || true