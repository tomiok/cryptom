package main

import "cryptom/cli"

func main() {
	c := cli.CLI{}
	c.GetBalance("el nica")

	c.Send("el nica", "tomi", 1)
}
