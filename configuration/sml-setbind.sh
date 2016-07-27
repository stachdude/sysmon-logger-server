#!/bin/sh
chmod +x /opt/sml/sml
sudo setcap cap_net_bind_service+ep /opt/sml/sml
