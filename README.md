# Title: go ftp server and data transformer

This project was to solve a real business problem, and have a go at learning go. A previously working netsuite script is broken due to a file exceeding the 100mb limit. This project aims to transform that file, and makes it avaliable to the script. 

There were some things I took for granted- the pkg/sftp package is not a sftp server. it is a SFTP protocol handler, not a full blown sftp server. (i was expecting the package to have the ability to restrict connected users to certain dir )

this is a ftp proxy to solve a problem where a netsuite script cannot handle a file larger than 100mb


# Plan
1. get things working locally on a docker image
2. get things working locally on a docker swarm
3. move this repo to the monolith

# deploy to existing swarm plan
- update firewall
- deploy local secrets to existing docker swarm