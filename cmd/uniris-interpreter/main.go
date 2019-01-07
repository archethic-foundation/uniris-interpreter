package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	uniris "github.com/uniris/uniris-interpreter/pkg"
	"github.com/urfave/cli"
)

func main() {

	app := cli.NewApp()
	app.Name = "uniris-interpreter"
	app.Usage = "Interpreter for Uniris smart contract"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "file, f",
			Usage: "Interpret from a `FILE` source code",
		},
		cli.BoolFlag{
			Name:  "console",
			Usage: "Open console to interpret code instantly",
		},
	}

	app.Action = func(c *cli.Context) error {

		if c.String("file") != "" {
			code, err := ioutil.ReadFile(c.String("file"))
			if err != nil {
				return err
			}
			if err := uniris.Interpret(string(code)); err != nil {
				fmt.Printf("Error: %s\n", err)
			}
			return nil
		} else if c.Bool("console") {
			fmt.Println("Type Ctrl-C to exit the console")
			for {
				text := read()
				if err := uniris.Interpret(text); err != nil {
					fmt.Printf("Error: %s\n", err)
				}
			}
		}

		return cli.ShowAppHelp(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

func read() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	text, _ := reader.ReadString('\n')

	return strings.Trim(text, "")
}
