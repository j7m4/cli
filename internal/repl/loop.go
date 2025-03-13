package repl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/ctrlplanedev/cli/internal/api"
	"github.com/spf13/viper"
)

type replCtx struct {
	systemsFilter      []string
	deploymentsFilter  []string
	environmentsFilter []string
	systemSlugs        []string
	deploymentSlugs    []string
	environmentSlugs   []string
}

type pageableResponse struct {
	data []map[string]interface{}
}

func newReplContext(ctx context.Context) (*replCtx, error) {
	var err error
	var envSlugs []string
	var sysSlugs []string
	var deplSlugs []string

	envSlugs, err = listEnvSlugs(ctx)
	if err != nil {
		return nil, err
	}

	sysSlugs, err = listSysSlugs(ctx)
	if err != nil {
		return nil, err
	}

	deplSlugs, err = listDeplSlugs(ctx)
	if err != nil {
		return nil, err
	}

	return &replCtx{
		systemsFilter:      make([]string, 0),
		deploymentsFilter:  make([]string, 0),
		environmentsFilter: make([]string, 0),
		systemSlugs:        sysSlugs,
		deploymentSlugs:    deplSlugs,
		environmentSlugs:   envSlugs,
	}, nil
}

func listEnvSlugs(ctx context.Context) ([]string, error) {
	var err error
	var response *http.Response

	client, err := buildApiClient()
	if err != nil {
		return nil, err
	}
	environmentSlugs := make([]string, 0)
	response, err = client.ListEnvironments(ctx)
	if err != nil {
		return environmentSlugs, err
	}
	if response.Body == nil {
		return environmentSlugs, fmt.Errorf("no environmentList response body")
	}
	envList := pageableResponse{}
	err = json.NewDecoder(response.Body).Decode(&envList)
	if err != nil {
		return environmentSlugs, err
	}
	for _, env := range envList.data {
		environmentSlugs = append(environmentSlugs, env["slug"].(string))
	}
	return environmentSlugs, nil
}

func listSysSlugs(ctx context.Context) ([]string, error) {
	var err error
	var response *http.Response

	client, err := buildApiClient()
	if err != nil {
		return nil, err
	}
	systemSlugs := make([]string, 0)
	response, err = client.ListSystems(ctx)
	if err != nil {
		return systemSlugs, err
	}
	if response.Body == nil {
		return systemSlugs, fmt.Errorf("no environmentList response body")
	}
	envList := pageableResponse{}
	err = json.NewDecoder(response.Body).Decode(&envList)
	if err != nil {
		return systemSlugs, err
	}
	for _, env := range envList.data {
		systemSlugs = append(systemSlugs, env["slug"].(string))
	}
	return systemSlugs, nil
}

func listDeplSlugs(ctx context.Context) ([]string, error) {
	var err error
	var response *http.Response

	client, err := buildApiClient()
	if err != nil {
		return nil, err
	}
	deploymentSlugs := make([]string, 0)
	response, err = client.ListSystems(ctx)
	if err != nil {
		return deploymentSlugs, err
	}
	if response.Body == nil {
		return deploymentSlugs, fmt.Errorf("no environmentList response body")
	}
	envList := pageableResponse{}
	err = json.NewDecoder(response.Body).Decode(&envList)
	if err != nil {
		return deploymentSlugs, err
	}
	for _, env := range envList.data {
		deploymentSlugs = append(deploymentSlugs, env["slug"].(string))
	}
	return deploymentSlugs, nil
}

type ReplCtx interface {
	clearFilters()
	showFilters()
	buildPrompt() string
	listSystemSlugs() func(string) []string
	listDeploymentSlugs() func(string) []string
	listEnvironmentSlugs() func(string) []string
}

func (rc *replCtx) clearFilters() {
	rc.systemsFilter = make([]string, 0)
	rc.deploymentsFilter = make([]string, 0)
	rc.environmentsFilter = make([]string, 0)
}

func (rc *replCtx) showFilters() {
	var display string
	if len(rc.systemsFilter) > 0 {
		display += fmt.Sprintf("systems: %v\n", strings.Join(rc.systemsFilter, ", "))
	}
	if len(rc.deploymentsFilter) > 0 {
		display += fmt.Sprintf("deployments: %v\n", strings.Join(rc.deploymentsFilter, ", "))
	}
	if len(rc.environmentsFilter) > 0 {
		display += fmt.Sprintf("environments: %v\n", strings.Join(rc.environmentsFilter, ", "))
	}
	fmt.Print(display)
}

func (rc *replCtx) buildPrompt() string {
	var filterSummary = ""
	if len(rc.systemsFilter) > 0 {
		filterSummary += fmt.Sprintf("s{%d}", len(rc.systemsFilter))
	}
	if len(rc.deploymentsFilter) > 0 {
		filterSummary += fmt.Sprintf("d{%d}", len(rc.deploymentsFilter))
	}
	if len(rc.environmentsFilter) > 0 {
		filterSummary += fmt.Sprintf("e{%d}", len(rc.environmentsFilter))
	}
	if filterSummary != "" {
		filterSummary = fmt.Sprintf("F[%s]", filterSummary)
	}
	return fmt.Sprintf("\033[31m%s»\033[0m", filterSummary)
}

func (rc *replCtx) listSystemSlugs() func(string) []string {
	return func(line string) []string {
		return rc.systemsFilter
	}
}

func (rc *replCtx) listDeploymentSlugs() func(string) []string {
	return func(line string) []string {
		return rc.deploymentSlugs
	}
}

func (rc *replCtx) listEnvironmentSlugs() func(string) []string {
	return func(line string) []string {
		return rc.environmentSlugs
	}
}

func buildCompleter(ctx ReplCtx) *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("filter",
			readline.PcItem("system",
				readline.PcItemDynamic(ctx.listSystemSlugs())),
			readline.PcItem("deployment",
				readline.PcItemDynamic(ctx.listDeploymentSlugs())),
			readline.PcItem("environment",
				readline.PcItemDynamic(ctx.listEnvironmentSlugs())),
			readline.PcItem("clear"),
			readline.PcItem("show"),
		),
		readline.PcItem("exit"),
	)
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func StartLoop(cmd *cobra.Command) error {
	var err error

	ctx, err := newReplContext(cmd.Context())
	if err != nil {
		return err
	}

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[31m»\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    buildCompleter(ctx),
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	if err != nil {
		return err
	}
	defer l.Close()
	l.CaptureExitSignal()

	log.SetOutput(l.Stderr())
	for {
		//ctx.showFilters()
		l.SetPrompt(ctx.buildPrompt())
		line, err := l.Readline()
		if errors.Is(err, readline.ErrInterrupt) {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		words := make([]string, 0)
		for _, w := range strings.Split(line, " ") {
			if len(w) > 0 {
				words = append(words, w)
			}
		}
		if len(words) == 0 {
			continue
		}
		switch {
		case words[0] == "filter":
			if len(words) < 2 {
				log.Println("filter <system|deployment|environment> <name>")
				break
			}
			switch words[1] {
			case "system":
				if len(words) < 3 {
					log.Println("filter system <name>")
					break
				}
				ctx.systemsFilter = append(ctx.systemsFilter, words[2])
			case "deployment":
				if len(words) < 3 {
					log.Println("filter deployment <name>")
					break
				}
				ctx.deploymentsFilter = append(ctx.deploymentsFilter, words[2])
			case "environment":
				if len(words) < 3 {
					log.Println("filter environment <name>")
					break
				}
				ctx.environmentsFilter = append(ctx.environmentsFilter, words[2])
			case "clear":
				ctx.clearFilters()
			case "show":
				ctx.showFilters()
			default:
				log.Println("filter <system|deployment|environment> <name>")
			}
		case line == "exit":
			return nil
		case line == "":
		default:
			log.Println("you said:", strconv.Quote(line))
		}
	}
	return nil
}

func buildApiClient() (*api.ClientWithResponses, error) {
	apiURL := viper.GetString("url")
	apiKey := viper.GetString("api-key")
	return api.NewAPIKeyClientWithResponses(apiURL, apiKey)
}
