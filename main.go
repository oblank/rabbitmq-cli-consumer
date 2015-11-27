package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/oBlank/rabbitmq-cli-consumer/command"
	"github.com/oBlank/rabbitmq-cli-consumer/config"
	"github.com/oBlank/rabbitmq-cli-consumer/consumer"
	"io"
	"log"
	"os"
	"path/filepath"
)

var default_concurrency = 5

func main() {
	app := cli.NewApp()
	app.Name = "rabbitmq-cli-consumer"
	app.Usage = "Consume RabbitMQ easily to any cli program"
	app.Author = "Richard van den Brand"
	app.Email = "richard@vandenbrand.org"
	app.Version = "1.2.0"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "concurrency, n",
			Usage: "Number of Concurrency, default is 5",
		},
		cli.StringFlag{
			Name:  "executable, e",
			Usage: "Location of executable",
		},
		cli.StringFlag{
			Name:  "configuration, c",
			Usage: "Location of configuration file",
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "Enable verbose mode (logs to stdout and stderr)",
		},
	}
	app.Action = func(c *cli.Context) {
		concurrency := c.Int("concurrency")

		if c.String("configuration") == "" && c.String("executable") == "" {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		verbose := c.Bool("verbose")

		logger := log.New(os.Stderr, "", log.Ldate|log.Ltime)
		cfg, err := config.LoadAndParse(c.String("configuration"))
		if concurrency > 0 {
			cfg.Concurrency.Max = concurrency
		}
		if cfg.Concurrency.Max <= 0 {
			cfg.Concurrency.Max = default_concurrency
		}

		if err != nil {
			logger.Fatalf("Failed parsing configuration: %s\n", err)
		}

		errLogger, err := createLogger(cfg.Logs.Error, verbose, os.Stderr)
		if err != nil {
			logger.Fatalf("Failed creating error log: %s", err)
		}

		infLogger, err := createLogger(cfg.Logs.Info, verbose, os.Stdout)
		if err != nil {
			logger.Fatalf("Failed creating info log: %s", err)
		}

		//todo
		exec_path := c.String("executable")
		if !filepath.IsAbs(exec_path) {
			localtion, err := filepath.Abs(exec_path)
			if err != nil {
				logger.Fatalf("Failed executable path log: %s", err)
			}
			exec_path = localtion
		}
		fmt.Println(exec_path)
		factory := command.Factory(exec_path)

		client, err := consumer.New(cfg, factory, errLogger, infLogger)
		if err != nil {
			errLogger.Fatalf("Failed creating consumer: %s", err)
		}

		client.Consume()
	}

	app.Run(os.Args)
}

func createLogger(filename string, verbose bool, out io.Writer) (*log.Logger, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)

	if err != nil {
		return nil, err
	}

	var writers = []io.Writer{
		file,
	}

	if verbose {
		writers = append(writers, out)
	}

	return log.New(io.MultiWriter(writers...), "", log.Ldate|log.Ltime), nil
}
