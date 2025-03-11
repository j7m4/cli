package repl

import (
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

type replCtx struct {
	systemsFilter      []string
	deploymentsFilter  []string
	environmentsFilter []string
}

func newContext() *replCtx {
	return &replCtx{
		systemsFilter:      make([]string, 0),
		deploymentsFilter:  make([]string, 0),
		environmentsFilter: make([]string, 0),
	}
}

type Ctx interface {
	clear()
	buildPrompt() string
	listSystemSlugs() func(string) []string
	listDeploymentSlugs() func(string) []string
	listEnvironmentSlugs() func(string) []string
}

func (rc *replCtx) clear() {
	rc.systemsFilter = make([]string, 0)
	rc.deploymentsFilter = make([]string, 0)
	rc.environmentsFilter = make([]string, 0)
}

func (rc *replCtx) show() {
	fmt.Printf("Systems: %v\n", rc.systemsFilter)
	fmt.Printf("Deployments: %v\n", rc.deploymentsFilter)
	fmt.Printf("Environments: %v\n", rc.environmentsFilter)
}

func (rc *replCtx) buildPrompt() string {
	return fmt.Sprintf(
		"s[%d] d[%d] e[%d] \033[31m»\033[0m ",
		len(rc.systemsFilter),
		len(rc.deploymentsFilter),
		len(rc.environmentsFilter),
	)
}

func (rc *replCtx) listSystemSlugs() func(string) []string {
	return func(line string) []string {
		names := []string{
			"cluster-security",
			"cluster-tooling",
			"ephemeral-data-services",
		}
		return names
	}
}

func (rc *replCtx) listDeploymentSlugs() func(string) []string {
	return func(line string) []string {
		names := []string{
			"teleport-kube-agent",
			"datadog",
			"kubernetes-replicator",
			"cert-manager",
			"minio",
			"mysql",
			"kafka",
			"google-emulators",
			"clickhouse",
			"redis",
		}
		return names
	}
}

func (rc *replCtx) listEnvironmentSlugs() func(string) []string {
	return func(line string) []string {
		names := []string{
			"ephemeral-instances",
			"ci",
			"testing-clusters",
		}
		return names
	}
}

func buildCompleter(ctx Ctx) *readline.PrefixCompleter {
	return readline.NewPrefixCompleter(
		readline.PcItem("filter",
			readline.PcItem("system",
				readline.PcItemDynamic(ctx.listSystemSlugs())),
			readline.PcItem("deployment",
				readline.PcItemDynamic(ctx.listDeploymentSlugs())),
			readline.PcItem("environment",
				readline.PcItemDynamic(ctx.listEnvironmentSlugs())),
			readline.PcItem("clear"),
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

func StartLoop() error {
	ctx := newContext()

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
		l.SetPrompt(ctx.buildPrompt())
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
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
				ctx.clear()
			case "show":
				ctx.show()
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
