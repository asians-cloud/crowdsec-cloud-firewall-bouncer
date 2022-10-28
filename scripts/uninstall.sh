#!/usr/bin/env bash

SYSTEMD_PATH_FILE="/etc/systemd/system/crowdsec-cloud-firewall-bouncer.service"
LOG_FILE="/var/log/crowdsec-cloud-firewall-bouncer.log"
CONFIG_DIR="/etc/crowdsec/crowdsec-cloud-firewall-bouncer/"
BIN_PATH_INSTALLED="/usr/local/bin/crowdsec-cloud-firewall-bouncer"

uninstall() {
	systemctl stop crowdsec-cloud-firewall-bouncer
	rm -rf "${CONFIG_DIR}"
	rm -f "${SYSTEMD_PATH_FILE}"
	rm -f "${BIN_PATH_INSTALLED}"
	rm -f "${LOG_FILE}"
}

uninstall

echo "crowdsec-cloud-firewall-bouncer uninstall successfully"