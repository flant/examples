#!/bin/bash

while true;
do
  if [ -f /etc/rsyslog.d/remote_v2.conf ]; then
    rm -r /etc/rsyslog.d/remote_v2.conf
  fi
  if [ -f /etc/rsyslog.d/remote_v3.conf ]; then
    rm -r /etc/rsyslog.d/remote_v3.conf
  fi
  if [ -f /etc/rsyslog.d/remote_v4.conf ]; then
    rm -r /etc/rsyslog.d/remote_v4.conf
  fi
  if [ ! -f /etc/rsyslog.d/remote_v5.conf ]; then
    echo 'auth,authpriv.* @@${SYSLOG_SERVER}' > /etc/rsyslog.d/remote_v5.conf
    systemctl restart syslog
    echo "rsyslog configured for remote syslog host"
  fi
  if [ ! -f /etc/rsyslog.d/30-audisp.conf ]; then
    echo "if \$programname == 'audispd' then @@${SYSLOG_SERVER}" > /etc/rsyslog.d/30-audisp.conf
    systemctl restart syslog
  fi
  if grep -q "LogLevel INFO" /etc/ssh/sshd_config; then
    sed -i 's/LogLevel INFO/LogLevel VERBOSE/' /etc/ssh/sshd_config
    systemctl restart sshd
  fi
  if grep -q "ForwardToSyslog=no" /etc/systemd/journald.conf; then
    sed -i 's/ForwardToSyslog\=no/ForwardToSyslog\=yes/' /etc/systemd/journald.conf
    systemctl restart systemd-journald
  fi
sleep 600
done
