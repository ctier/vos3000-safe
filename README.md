# vos3000-safe
vos-web安全认证，致力于解决web安全

## 修改密码 /etc/init.d/safed
```shell
#ettings user
export VOS_SAFE_USERNAME=vos
export VOS_SAFE_PASSWORD=vos
```
## 在系统centos 6.x 64位 2.1.6.0测试通过，安装步骤如下
```shell
#safe
curltps://oss.1nth.com/vospag/iptables -o /etc/sysconfig/iptables
#增加白名单 ims     iptables  -I INPUT -s 10.0.0.0/8 -j ACCEPT  && service iptables save
echo -e "01 01 * * * /etc/init.d/iptables restart" >> /var/spool/cron/root
HOST=$(ip addr|grep inet|grep brd|grep -v "lo:"| awk  '{print $2}'|awk -F"/" '{print $1}'  | grep -Ev "^10|^172|^100|^192")
[  -z  "$HOST" ] && HOST=$(ip addr|grep inet|grep brd|grep -v "lo:"| awk  '{print $2}'|awk -F"/" '{print $1}')
echo $HOST
openssl rand -writerand .rnd
openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout /etc/ssl/key.pem -out /etc/ssl/cert.pem -subj "/O=${HOST}/CN=${HOST}"
curl http://release -o /etc/init.d/safed
curl http://release -o /usr/local/bin/safe
chmod +x /etc/init.d/safed
chmod +x /usr/local/bin/safe
#dos2unix /etc/init.d/safed   dos2unix -k -o filename
#sed -i -e 's/\r$//' /etc/init.d/safed
/etc/init.d/safed restart
chkconfig safed on
chkconfig  iptables on

yum install -y fail2ban
/etc/init.d/fail2ban restart
chkconfig fail2ban on
egrep -v "^$|^#|^;" /etc/fail2ban/jail.conf
echo '[DEFAULT]
ignoreip = 127.0.0.1
bantime  = 36000
findtime  = 600
maxretry = 3
[ssh-iptables]
enabled  = true
filter   = sshd
action   = iptables[name=SSH, port=10022, protocol=tcp]
           sendmail-whois[name=SSH, dest=root, sender=fail2ban@example.com]
#logpath  = /var/log/secure
maxretry = 5' >/etc/fail2ban/jail.local 
#关闭fail2ban日志
sed -ri '/^[^#]*\/var\/log\/messages/s@^@#@' /etc/fail2ban/jail.conf
SIP_PORT=5060,6060
sed -i 's/SIP_PORT=5060,6060/SIP_PORT=1980,6060/g' /home/kunshi/mbx3000/etc/softswitch.conf
sed -i 's/port="8080"/port="8888"/g' /home/kunshiweb/base/apache-tomcat/conf/server.xml
mysqladmin -u root password "xiaofan@1"
#sed -i '/^Listen/cListen 10000'  /etc/httpd/conf/httpd.conf
sed -i '/^#Port/cPort 10022' /etc/ssh/sshd_config
sed -i "s/dport 88/dport 8000/g" /etc/sysconfig/iptables
sed -i "s/dport 22/dport 10022/g" /etc/sysconfig/iptables
sed -i "s/dport 2080/dport 1980/g" /etc/sysconfig/iptables
chattr +i /etc/sysconfig/iptables
chattr +i /etc/rc.d/rc.local
#update uuid
sed -i '/^ACCESS_UUID=/cACCESS_UUID=vos30002160' /home/kunshi/vos3000/etc/server.conf
#update port vos-client 6.0测试修改就会超时 查看别人8.0貌似可以用
#sed -i '/^GUI_SERVER_PORT=/cGUI_SERVER_PORT=2020' /home/kunshi/vos3000/etc/server.conf

rm -rf /etc/yum.repos.d/*
#2020-05-21发现tomcat漏洞 于是禁止服务器下载命令
for i in {yum,wget,curl}; do mv /usr/bin/$i /usr/bin/myki$i; done

echo "alias curl='echo >> ~/vos.log'
alias wget='echo >> ~/vos.log'
alias yum='echo  >> ~/vos.log'">> /etc/bashrc
```