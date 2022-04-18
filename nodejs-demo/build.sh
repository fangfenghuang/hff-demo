#!/bin/bash
set -x
node -v
yarn install --registry=https://registry.npmjs.org
yarn run build
