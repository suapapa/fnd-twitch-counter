# fnd-twitch-counter

Twitch follower counter of mine(suapapa) in FND display module(tm1638) in RaspberryPi

![image](photo/fnd-twitch-counter.jpg)

## build and rund

Build in host:

    GOOS=linux GOARCH=arm go build

Run in target (raspberrypi):

    ./fnd-twitch-counter -t 50 -l 30 # update every 30 secs and target followers are 50

## install systemd service

In target fix `fnd-twitch-counter.service` to your client-id and client-secret and;

    $ sudo cp fnd-twitch-counter.service /lib/systemd/system/fnd-twitch-counter.service
    $ sudo systemctl enable fnd-twitch-counter
