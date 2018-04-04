[![Build Status](https://travis-ci.org/kcmerrill/tracer.svg?branch=master)](https://travis-ci.org/kcmerrill/tracer) 

![tracer](tracer.jpg "tracer")

> A simple Cron job `||` ETL Stage Monitor. 

# Tracer

A simple webserver designed to alert you when cron jobs fail, or when long running multi phased ETL processes fail to complete.

For more information [RTFM](TFM.md "additional documentation").

## Quick Start

```sh
$ tracer --bind 0.0.0.0:8080 --panic "echo {{ .Name }} was not seen in the expected duration of {{ .Duration }}"
```

## Sample Usage

```sh
# crontab
0 5,17 * * * http http://localhost:8080/my.script/12h; /scripts/script.sh && http http://localhost:8080/my.script

* * * * *  http http://localhost:8080/job.every.minute/1m < /panic/job.every.minute; /scripts/script.sh && http http://localhost:8080/job.every.minute
```