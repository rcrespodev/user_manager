#!/bin/sh
#!/bin/bash
#!/usr/bin/perl
#!/usr/bin/tcl
#!/bin/sed -f
#!/bin/awk -f
docker_ps=$(docker ps -f name=test_app_mysql -f name=test_app_redis)

mysql_upp=0
redis_upp=0
stop_test_services=0

for i in $docker_ps; do
  if [ "$i" = "test_app_mysql" ]; then
    mysql_upp=1
  fi
  if [ "$i" = "test_app_redis" ]; then
    redis_upp=1
  fi
  if [ "$redis_upp" -eq 1 ] && [ "$mysql_upp" -eq 1 ]; then
      stop_test_services=1
      break
  fi
done

if [ $stop_test_services -eq 1 ]; then
  make stop_test_services
fi

export PATH=$PATH:/usr/local/go/bin && make go_test
