TeslaBiller
=================================================

Quick start
---------------------------

- You must have TeslaMate up and running with the postgres database port exposed
- You can edit the docker-compose.yml file for your setup (change ip, port, username and password environment variables)
- You don't need to build the server, there is a docker image available on the dockerhub
- For the client, you can download the released apk of the latest version

How it work
---------------------------

It connects to the teslamate database to retreive the charges done by your car at your home (or in the future any other locations) and calculate the cost of your car.
For now it depends on an external counter to calculate the extra costs not included in Teslamate, in the future I will make it optional.
It's based on my electricity provider : I have 3 diferent prices during the day for the electricity and the price change each months ... Their is plans to make it configurable ...

For now it suits my needs, if you want any changes for your needs feel free to open a ticket or, better make a pull request :)
