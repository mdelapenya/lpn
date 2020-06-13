#!/bin/bash

sudo apt update -y
sudo apt install realpath -y
go get gotest.tools/gotestsum
gem update --system 3.0.8
gem install bundler
bundle install