## Block Header Syncer

### Deploy instructions

You need:

* Docker (tested in the `18.03.0-ce` version)
* Docker compose (tested in the `1.18.0` version)

Everything else happens inside the containers!

### Get this directory

Download this repository with git, or just modify the included `deploy.sh` file
with your server data. The latter will upload all this directory contents into a
remote directory and run `prepare_box` and `run_box`.

### Prepare your box

	./prepare_box

Will take care of doing one-time downloads and some compilation.

### Run and stop this box

	./run_box

Will execute the docker-compose spec.

	./stop_box

Will stop the docker-compose spec.

### Profit

Now you have a service that synchronizes the ethereum block header into a redis DB.

### I want to learn more

Visit the README.md file inside the block header syncer service directory.