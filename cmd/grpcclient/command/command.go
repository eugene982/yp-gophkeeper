// Package command команды пользователя
package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
)

type CommandFunc func(map[string]string) error

type Command struct {
	cmd       CommandFunc
	prompts   []string
	userInput map[string]string
	completer prompt.Completer
}

func New(cmd CommandFunc, args []string, prompt ...string) *Command {
	n := len(args)
	if len(prompt) < n {
		cmd = func(map[string]string) error {
			return fmt.Errorf("ошибка аргументов '%s'", strings.Join(args, " "))
		}
		return &Command{
			cmd: cmd,
		}
	}

	c := &Command{
		cmd:       cmd,
		prompts:   prompt,
		userInput: make(map[string]string, n),
	}
	// установим переданные аргументы как введённые пользователем комманды
	for i := 0; i < n; i++ {
		c.userInput[prompt[i]] = args[i]
	}
	return c
}

func (p *Command) Exec() {
	for _, v := range p.prompts {
		if _, ok := p.userInput[v]; !ok {
			p.userInput[v] = prompt.Input(v+": ", p.Complete)
		}
	}
	if err := p.cmd(p.userInput); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (p *Command) Complete(d prompt.Document) []prompt.Suggest {
	if p.completer == nil {
		return nil
	}
	return p.completer(d)
}

func (p *Command) SetCompleter(c prompt.Completer) {
	p.completer = c
}
