## Timelapse Camera

Just a quick docker container to read a RTMP camera on interval and dump a lossless PNG into a directory for later processing (like to make a timelapse video)

Config should look like this:
```
{
  "stream_url": "rtsp://[ip]:[port]/pathtovideostream",
  "interval": "* * * * *",
  "output_dir": "./output"
}
```
The interval is crontab format.

The output directory should be a volume so it can persist between container restarts.

An example docker-compose.yml would be:
```
services:
  timelapse-camera:
    build: .
    container_name: timelapse-camera
    restart: unless-stopped
    volumes:
      - ./config.json:/config/config.json
      - ./output:/output
```

Then simply `docker compose up -d`

The output directory will then contain a series of PNG files named with the timestamp.png