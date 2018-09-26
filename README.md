# ScryfallSpoilerBot [![Build Status](https://travis-ci.com/Adi146/ScryfallSpoilerBot.svg?branch=master)](https://travis-ci.com/Adi146/ScryfallSpoilerBot)

This bot posts recent magic card spoilers using [Scryfall](https://scryfall.com/) to [Pushbullet](https://www.pushbullet.com/)

## Basic Setup
* Go to [Pushbullet Account Settings](https://www.pushbullet.com/#settings/account) and create an Access Token
* Copy the Token to your config.yaml
```yaml
messengers:
    pushbullet:
        - api_token: ******
```
* Run ScryfallSpoilerBot
```
ScryfallSpoilerBot --config=/path/to/your/config.yaml
```

## Run With Docker
* Build a docker image
```
docker build -t scryfall-spoiler-bot .
```
* Run docker container
```
docker run -d --name="scryfall-spoiler-bot" -v /path/to/your/config/directory:/config scryfall-spoiler-bot
```
