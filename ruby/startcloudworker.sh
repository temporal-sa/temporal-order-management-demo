#!/bin/bash
source ../setcloudenv.sh
bundle install
bundle exec ruby worker.rb
