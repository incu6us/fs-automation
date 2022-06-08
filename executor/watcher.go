package executor

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher *fsnotify.Watcher
	// matchedRules   map[string]*Rule
	// unmatchedRules map[string]*Rule
	ruleStorage *RuleStorage
}

func NewWatcher() (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		watcher: watcher,
		// matchedRules:   make(map[string]*Rule),
		// unmatchedRules: make(map[string]*Rule),
		ruleStorage: NewRuleStorage(),
	}

	return w, nil
}

func (w *Watcher) Add(rule *Rule) error {
	w.ruleStorage.Add(rule.Path, rule.Operation, rule)

	err := w.watcher.Add(rule.Path)
	if err != nil {
		var pathErr *fs.PathError
		if errors.As(err, &pathErr) {
			path := findExistedPath(rule.Path)
			err := w.watcher.Add(path)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}

func (w *Watcher) Watch() error {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return nil
			}

			rule := w.ruleStorage.Contains(event.Name, operations[event.Op])
			if rule == nil {
				continue
			}
			executeRule(rule, event.Name, event.Op)
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return nil
			}
			fmt.Printf("--------\n")
			fmt.Printf("ERR: %s\n", err)
		}
	}
}

func (w *Watcher) Close() error {
	return w.watcher.Close()
}

var operations = map[fsnotify.Op]Operation{
	fsnotify.Create: CreateOperation,
	fsnotify.Write:  WriteOperation,
	fsnotify.Remove: RemoveOperation,
	fsnotify.Rename: RenameOperation,
	fsnotify.Chmod:  ChmodOperation,
}

func executeRule(rule *Rule, fullPath string, action fsnotify.Op) {
	stdOut, stdErr, err := executeCmd(replaceCmdVariables(rule, fullPath, action))
	if err != nil {
		panic(err)
	}

	if stdOut != "" {
		fmt.Printf("STDOUT: %s", stdOut)
	}

	if stdErr != "" {
		fmt.Printf("STDERR: %s", stdErr)
	}
}

func replaceCmdVariables(rule *Rule, fullPath string, action fsnotify.Op) string {
	s := strings.ReplaceAll(rule.Cmd, "{{path}}", rule.Path)
	s = strings.ReplaceAll(s, "{{fullpath}}", fullPath)
	s = strings.ReplaceAll(s, "{{operation}}", rule.Operation.String())
	s = strings.ReplaceAll(s, "{{action}}", action.String())
	return s
}

// executeCmd returns STDOUT, STDERR and execution error
func executeCmd(command string) (string, string, error) {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return stdOut.String(), stdErr.String(), err
	}

	return stdOut.String(), stdErr.String(), err
}

func findExistedPath(path string) string {
	path = filepath.Dir(path)
	_, err := os.Stat(path)
	if err != nil {
		return findExistedPath(path)
	}

	return path
}
