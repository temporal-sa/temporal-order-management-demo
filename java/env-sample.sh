#!/bin/bash
#
# Script used to set various environment variables that are needed to configure the order management demo.
# Please edit to provide the correct values for your envionnment
#
# Once set source this file prior to running the app locally.
#

# *************************************************************************************************
# Env vars for main order demo namespaac
# Notes:
# - The namespace for self hosted is hard wired in application.yaml for default.
# - Temporal address for self hosted defaults to localhost:7233
# - Some of the code defaults to picking the private key from the PCCS8 env variable.
# *************************************************************************************************

export TEMPORAL_NAMESPACE=<Your Namespace>.<Account ID>
export TEMPORAL_TASK_QUEUE=orders
export TEMPORAL_ADDRESS=<Namespace>.<account id>.tmprl.cloud:7233
export TEMPORAL_CERT_PATH=<Path>/<To>/<Client public certificate>.pem
export TEMPORAL_KEY_PATH=<Path>/<To>/<Client private key file>
export TEMPORAL_KEY_PKCS8_PATH=${TEMPORAL_KEY_PATH}

# *************************************************************************************************
# Additional env vars for using nexus to demonstrate the shipping fulfillment in another namespace.
# Notes:
# - For self hosted this will default to using the default namespace
# - Nexus only allows unique endpoints in an account so you may need to change this to a 
#   Unique value to match the endpoint you created.  All code references are from this env var.
# *************************************************************************************************
export TEMPORAL_SHIPPING_NAMESPACE=<Your SHIPPING Namespace>.<Account ID>
export TEMPORAL_SHIPPING_TASK_QUEUE=shipping
export TEMPORAL_SHIPPING_ADDRESS=<Namespace>.<account id>.tmprl.cloud:7233
export TEMPORAL_SHIPPING_CERT_PATH=<Path>/<To>/<Client public certificate for Shipping NS>.pem
export TEMPORAL_SHIPPING_KEY_PATH=<Path>/<To>/<Client private key file for Shipping NS>
export TEMPORAL_SHIPPING_CERT_RELOAD_PERIOD=30   # Refresh period for certificaPath>/<To>/<Client ptes (in minutes)
export TEMPORAL_SHIPPING_NEXUS_ENDPOINT=shipping-endpoint

