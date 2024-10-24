#!/bin/bash

temporal server start-dev \
    --http-port 7243 \
    --dynamic-config-value system.enableNexus=true \
    --search-attribute OrderStatus="Keyword"
