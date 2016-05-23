# Folder Transfer

This is two applications that work in tandem to allow one host to quicky 
download the contents of a folder on the serving host. A connection is opened
for each file to be downloaded. Beware folders with lots of files, lots of 
connections will be made.

The main use case is transferring a folder of 1-30 videos from one server to
another.

## Why?

Because I couldn't find anything that was automatic and simple enough to 
just download the contents of a directory.

### What about rsync?

Only one file at a time. Plus I don't care about delta-changes, I'm always
downloading fresh, the files don't ever exist on the pulling host.

### What about http server and aria2?

I was using aria2, it is awesome. However, I coulnd't find a way to download
a whole directory without manual work of creating a file with a list of URLS
for each file in the directory.

## Should I use this?

Probably not, but feel free. There is no authentication at all. So if you keep
the serving process up, anyone who knows the port can download everything.

I personally use this on an ad-hoc basis. I find a folder with some data I want
on another machine, I start the server, ssh to the puller and start the get.

## Usage

### Serving host

```sh
	$ ft serve
```

### Pulling host

```sh
	$ ft get serving-host
```