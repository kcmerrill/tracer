package main

import (
	"flag"

	"github.com/kcmerrill/tracer/pkg/tracer"
)

func main() {
	token := flag.String("token", "", "HTTP Basic Authentication token(username)")
	bind := flag.String("bind", "80", "Bind web server, Example <0.0.0.0:8080>")
	panic := flag.String("panic", "touch /tmp/tracer.{{ .Name }}", "Command to execute on panic")
	flag.Parse()

	// giddy up
	tracer.Start(*token, *bind, *panic)
}
