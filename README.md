# sutaba-server
Twitter bot that recognize whether the replied image is photographed in Starbucks or not.

## Configuration
You can specify server configurations by flags or config file.
Also, you must specify twitter credentials via environment variable.

### flags
```bash
Start server

Usage:
  sutaba-server start [flags]

Flags:
      --bot-id int                 bot twitter user id
      --classifier-server string   classifier server url
      --error-message string       text of tweet for error notification
  -h, --help                       help for start
      --keyword string             process only tweets which contain this value
      --owner-id int               owner twitter user id (error tweet will be send to owner if something is failed)
      --sorry-message string       text of tweet to send to user if process is failed

Global Flags:
      --config string   config file (default is $HOME/.sutaba-server.yaml)
```

### config file
See [config.yml](https://github.com/mpppk/sutaba-server/config.yml)

### Twitter credentials
You must specify below environment variables.

* TWITTER_CONSUMER_KEY
* TWITTER_CONSUMER_SECRET
* OWNER_TWITTER_ACCESS_TOKEN
* OWNER_TWITTER_ACCESS_TOKEN_SECRET
* BOT_TWITTER_ACCESS_TOKEN
* BOG_TWITTER_ACCESS_TOKEN_SECRET

## Run server

### binary

Download from GitHub Releases and put it anywhere in your executable path.

```bash
$ sutaba-server start --config ./config.yml
```

### from source

```bash
$ git clone https://github.com/mpppk/sutaba-server
$ cd sutaba-server
$ go run main.go start --config ./config.yml
```

### Docker

```bash
$ git clone https://github.com/mpppk/sutaba-server
$ cd sutaba-server
$ docker-compose up
```

## Deploy
This repository includes deploy scripts for GCP(Container Registry and Cloud Run).  
Before run scripts, update scripts by your project id.

```bash
$ ./scripts/gcp/submit_container_image.sh
$ ./scripts/gcp/deploy_to_cloud_run.sh
```

