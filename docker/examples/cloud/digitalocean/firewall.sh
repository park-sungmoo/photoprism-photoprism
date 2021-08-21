#!/usr/bin/env bash

# Install ufw package if needed:
apt-get update && apt-get install --no-install-recommends ufw && apt-get autoclean && apt-get autoremove

# Basic ufw firewall setup allowing ssh, http, and https:
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow http
ufw allow https
ufw logging off
rm -f /var/log/ufw.log
ufw --force enable
