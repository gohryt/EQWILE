package checker

import (
	"context"
	"time"
)

type (
	Configuration struct {
		Interval time.Duration `toml:"interval"`
	}

	Checker struct {
		configuration Configuration
	}
)

func Constructor(configuration *Configuration) Checker {
	return Checker{
		configuration: *configuration,
	}
}

func (checker *Checker) Run(ctx context.Context) (err error) {
	ticker := time.Tick(checker.configuration.Interval * time.Second)

root:
	for {
		err = Run()
		if err != nil {
			return
		}

		select {
		case <-ticker:
		case <-ctx.Done():
			break root
		}
	}

	return
}

func Run() (err error) {
	return
}
