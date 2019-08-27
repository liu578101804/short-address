package main

import "github.com/liu578101804/short-address/app"

func main() {
	a := app.App{}
	a.Initialize(app.GetEnv())
	a.Run(":8000")
}