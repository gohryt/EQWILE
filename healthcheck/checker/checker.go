package checker

import (
	"context"
	"log"
	"time"

	"github.com/valyala/fasthttp"
)

type (
	Configuration struct {
		Interval time.Duration `toml:"interval"`

		List []URL
	}

	URL struct {
		URL        string
		CheckList  []string `toml:"checkList"`
		CheckCount int      `toml:"checkCount"`
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

func (checker *Checker) Run(ctx context.Context) (err error) {
	list := checker.configuration.List

	ticker := time.Tick(checker.configuration.Interval * time.Second)

root:
	for {
		for i := range list {
			err = checker.Check(list[i])
			if err != nil {
				return
			}
		}

		select {
		case <-ticker:
		case <-ctx.Done():
			break root
		}
	}

	return
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
