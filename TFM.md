# TFM

Use tracer to check cron jobs, or tasks that you expect to run at specific intervals through a simple HTTP request.

## Creating a new check

For the sake of demonstration lets assume that tracer is running on `localhost` port `80`. Creating a hook is as simple as a `GET`, `POST` to `http://localhost/<checkname>/<duration>` where `<checkname>` is a unique identifer that you can provide, where `<duration>` is a string that can be parsed via [golang duration](https://golang.org/pkg/time/#Duration). So you can do every hour(`1h`), every minute(`1m`), thirty seconds(`30s`), an hour and thirty minutes(`1h30m`) as an example.

By creating a check, if not cancelled or ok'd within the given duration `panic` can occur.

You can define `panic` in two ways. First, by default, you can pass in `--panic` when starting tracer. The next is `POST` to the same endpoint, where the body is the template command to be run. 

Note that the command is a golang template(using [sprig](https://github.com/Masterminds/sprig)) with a few values to customize the check.

1. Name - The name of the check
1. Duration - The given duration of the check
1. Panic - The panic template given
1. Created - When the check was created

## Cancelling or Ok'ing a check

Simply call `GET` on the checkname. Lets say you created a check by going to `http://localhost/my.cool.check/1h`, you can ok the check by going to `http://localhost/my.cool.check`.

Creating a check while a check is already in progress will cancel any current checks and recreate the check, so if this is not intended, be specific with your check names!

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

## Options

1. **--binding** determines which port to bind tracer to. By default, binds to `0.0.0.0:<port>`, however, you can pass in anything here. Example: `--binding localhost:12345`
1. **--token** if passed in, default no token, will enable HTTP Authentication as the `<username>`. If you're using [httpie](https://github.com/jakubroztocil/httpie#basic-auth) is here is a quick example: ```sh
http -a <token>: https://localhost:8080/<checkname>/<duration>
```
1. **--panic** is used to determine what to do, by default, if no panic is given on check creation. Use this to send emails, slack messages, bash scripts or anything that you can execute via a shell command. A quick example if you're using [alfred](https://github.com/kcmerrill/alfred). `--panic 'alfred /slack:msg "#alerts" "{{ .Name }} not found in {{ .Duration }}"'`

