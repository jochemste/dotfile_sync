# dotfile_sync
Configuration file synchroniser. An easy way to keep all your configuration files synchronised between devices, without having a git repo in your home directory.

## Features
* Synchronise files between devices, using a public or private Github repository
* Automatically updates files in the repository and/or local file system, based on last modification date/time.
* Allows a device to only upload configuration files, while never overwriting local files with the remote ones. This means you can have a single device that will update the repository, while the changes in the repository will never change the local files on that machine.
* Allow a device to only pull from the repository, without pushing to it. For example when experimeting with configuration, you can start off from your usual config, and safely change something, without it being pushed to the repository.

## Problems (Are on the short list to be sorted)
* Does not support synchronising multiple files with the same filename (not path, filename. This means that synchronising /root/.bashrc and /home/whatever/.bashrc can lead to problems and loss of one of the two)
* Does not contain enough unit tests at the moment. Unexpected behaviour may still be hidden in the code somewhere.
* Merging is not done. The latest version of the file is uploaded and the old one is overwritten. For local files, a back-up will be made. This back-up is overwritten every time a new backup is made. The back-up is stored in the same location and uses the same filename, followed by ".bak".

For other issues, please create a issue [here](https://github.com/jochemste/dotfile_sync/issues).

## Quick Start

### Step 1: Install

Make sure you have Go installed ([Get Go](https://go.dev/doc/install)), and run the following command:

```
go install github.com/jochemste/dotfile_sync@latest
```

### Step 2: Initialise synchronisation file

When installed, you should be able to run dotfile_sync (assuming it is placed in your path). On the first run, use the following command to create a configuration file (as root or sudo):

```
dotfile_sync init
```

The configuration will be stored in /etc/dotfile_sync/config.toml

Open it as root (or sudo) and adjust the "UserSettings" parameters to your liking. Do not change the "DoNotChange" parameters. They have that name for a reason.

The Origin and AccessToken will require you to create a repository on Github and an access token that will provide access to that repository and pull, commit and push to the repo. Checkout this link to create a Personal Access Token on Github: [How to generate a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)

Once you have entered the repository URL at Origin and the access token at AccessToken, you can set the Mode. In most cases you would leave it at 'all', but other possibilities are 'push' and 'pull'. See [Features](#features) for more info on this.

Lastly, enter the files you would like to sync at Files, as a array of strings. For example:

```
Files = ['/path/to/file1', '/path/to/another/file', '/root/file2']
```

Avoid having files with the same name, since this can overwrite one of them. So do not do:

```
Files = ['/root/.bashrc', '/home/harry/.bashrc']
```

This could cause issues. I'm sorting this out soon.

### Step 3: Synchronise

Once the configuration is done, you can see if it works, by synchronising. To do this, run the following command as root (or sudo):

```
dotfile_sync sync
```

And do the same thing on your other devices. I would recommend having one "master" device that will only push to the repo, so that there will be no issues with merging and such. Merging is not supported at the moment.
