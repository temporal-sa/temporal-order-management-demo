#!/bin/bash
bundle install
source ../setcloudenv.sh
bundle exec ruby worker.rb
