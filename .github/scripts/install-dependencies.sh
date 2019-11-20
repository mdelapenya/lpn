#!/bin/bash

sudo apt update -y
sudo apt install realpath -y
go get gotest.tools/gotestsum
gem install bundler
bundle install