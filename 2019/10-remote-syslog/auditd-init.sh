#/bin/bash

NODE_PATH='/mnt/node'

mkdir $NODE_PATH/tmp/auditd


/usr/sbin/chroot /mnt/node systemctl stop auditd
rm -rf $NODE_PATH/etc/audit/rules.d
/usr/sbin/chroot /mnt/node apt purge auditd -y

wget -P $NODE_PATH/tmp/auditd https://launchpad.net/ubuntu/+archive/primary/+files/auditd_2.8.2-1ubuntu1_amd64.deb
wget -P $NODE_PATH/tmp/auditd https://launchpad.net/ubuntu/+archive/primary/+files/libaudit1_2.8.2-1ubuntu1_amd64.deb
wget -P $NODE_PATH/tmp/auditd https://launchpad.net/ubuntu/+archive/primary/+files/libauparse0_2.8.2-1ubuntu1_amd64.deb
wget -P $NODE_PATH/tmp/auditd https://launchpad.net/ubuntu/+archive/primary/+files/libaudit-common_2.8.2-1ubuntu1_all.deb

for i in $(ls $NODE_PATH/tmp/auditd |grep '.deb'); do
    /usr/sbin/chroot /mnt/node dpkg -i /tmp/auditd/$i
    rm $NODE_PATH/tmp/auditd/$i
done

DEBIAN_FRONTEND=noninteractive /usr/sbin/chroot /mnt/node  apt install -y --no-install-recommends -o Dpkg::Options::="--force-confdef" audispd-plugins

mkdir -p /var/log/audit
chown -R root:adm /var/log/audit

sed -i 's/active = no/active = yes/g' $NODE_PATH/etc/audisp/plugins.d/syslog.conf

cat > $NODE_PATH/etc/audit/audit.rules << EOF
-D
-b 320
-a exit,always -F arch=b64 -F auid>=1000 -F auid!=4294967295 -S execve
-a exit,always -F arch=b32 -F auid>=1000 -F auid!=4294967295 -S execve
EOF

cat > $NODE_PATH/etc/audit/rules.d/audit.rules << EOF
-D
-b 320
-a exit,always -F arch=b64 -F auid>=1000 -F auid!=4294967295 -S execve
-a exit,always -F arch=b32 -F auid>=1000 -F auid!=4294967295 -S execve
EOF

rm -f $NODE_PATH/etc/rsyslog.d/30-audisp.conf

cat > $NODE_PATH/etc/audit/auditd.conf << EOF
local_events = yes
write_logs = yes
log_file = /var/log/audit/audit.log
log_group = adm
log_format = ENRICHED
flush = INCREMENTAL_ASYNC
freq = 50
max_log_file = 8
num_logs = 5
priority_boost = 4
disp_qos = lossy
dispatcher = /sbin/audispd
name_format = NONE
##name = mydomain
max_log_file_action = ROTATE
space_left = 75
space_left_action = SYSLOG
verify_email = yes
action_mail_acct = root
admin_space_left = 50
admin_space_left_action = SUSPEND
disk_full_action = SUSPEND
disk_error_action = SUSPEND
use_libwrap = yes
##tcp_listen_port = 60
tcp_listen_queue = 5
tcp_max_per_addr = 1
##tcp_client_ports = 1024-65535
tcp_client_max_idle = 0
enable_krb5 = no
krb5_principal = auditd
##krb5_key_file = /etc/audit/audit.key
distribute_network = no
EOF

/usr/sbin/chroot /mnt/node systemctl restart auditd
