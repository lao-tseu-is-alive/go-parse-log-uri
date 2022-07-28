#!/bin/bash 
goaccess /var/log/nginx/yourserver_yourdomain_com_ssl_access.log  -o /var/www/yourserver.yourdomain.com/html/stats/index.html --log-format=COMBINED --db-path /root/web_analysis_goaccess_data_store --restore  --persist  --keep-last=31
