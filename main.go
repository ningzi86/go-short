package main

import (
	"go-short/app"
	"go-short/env"
)

func main() {
	app := app.App{}
	app.Init(env.GetEnv())
	app.Run(":80")
}
