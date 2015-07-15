# munkiimport server
A proof of concept server to handle remote munkiimport commands.

Instead of using munkiimport directly, we can set up a remote server to accept file uploads over http and shell out to munkiimport.
A few potential uses:

* create a WebUI for the CLI averse users.
* automate git commits to your munki repo.

# Requirements

A OS X Machine with munkitools and configured munkiimport.

# Example Usage

HTTP PUT binary file to `http://munkiimport-server:8080/import`
curl example:

`curl -v http://localhost:8080/import/ --upload-file "~/Downloads/Firefox 39.0.dmg"`

# Web Upload

![img](http://i.imgur.com/3oR5oTd.png)
