package checker

import (
	"context"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

type (
	Configuration struct {
		Interval int `yaml:"interval"`

		List []URL `yaml:"List"`
	}

	URL struct {
		URL        string   `yaml:"URL"`
		CheckList  []string `yaml:"checkList"`
		CheckCount int      `yaml:"checkCount"`
	}

	Checker struct {
		configuration Configuration

		checkMap map[string]Check
		client   *fasthttp.Client
	}

	Check func(response *fasthttp.Response) any
)

func Constructor(configuration *Configuration) Checker {
	return Checker{
		configuration: *configuration,

		checkMap: make(map[string]Check),
		client: &fasthttp.Client{
			DialDualStack: true,
		},
	}
}

func (checker *Checker) Register(name string, check Check) {
	checker.checkMap[name] = check
}

func (checker *Checker) Run(ctx context.Context) {
	list := checker.configuration.List

	ticker := time.Tick(time.Duration(checker.configuration.Interval) * time.Second)

root:
	for {
		for i := range list {
			err := checker.Check(list[i])
			if err != nil {
				log.Println(err)
			}
		}

		select {
		case <-ticker:
		case <-ctx.Done():
			break root
		}
	}
}

func (checker *Checker) Check(URL URL) (err error) {
	request := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(request)

	request.Header.SetMethod(fasthttp.MethodGet)
	request.URI().Update(URL.URL)

	response := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(response)

	err = checker.client.Do(request, response)
	if err != nil {
		return
	}

	for i := range URL.CheckList {
		check := checker.checkMap[URL.CheckList[i]]

		if check == nil {
			continue
		}

		log.Println(URL.URL, URL.CheckList[i], check(response))
	}

	return
}
