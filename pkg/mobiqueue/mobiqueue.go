package mobiqueue

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type MobiQueue struct {
	wg       *sync.WaitGroup
	mu       *sync.Mutex
	errArray []string
}

func New() MobiQueue {
	return MobiQueue{
		wg: &sync.WaitGroup{},
		mu: &sync.Mutex{},
	}
}

func (m *MobiQueue) Add(chapterUrl string) {
	m.wg.Add(1)
	go m.generateMOBI(chapterUrl)
}

func (m *MobiQueue) Wait() {
	m.wg.Wait()
}

func (m *MobiQueue) Err() error {
	if len(m.errArray) > 0 {
		return errors.New(strings.Join(m.errArray, "\n"))
	}
	return nil
}

func (m *MobiQueue) generateMOBI(chapterDir string) {
	wg := m.wg
	mu := m.mu

	mu.Lock()

	cmdOutput := &bytes.Buffer{}
	cmd := exec.Command("kcc-c2e", "-m", "-f", "MOBI", "-g", "1.0", chapterDir)
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		template := "Error\n\tDir: %s\n\tErr: %s\n\tCmd: %s\n\tContent: \n%s\n=====\n"
		errMessage := fmt.Sprintf(template, chapterDir, err.Error(), fmt.Sprint("kcc-c2e", "-f", "MOBI", chapterDir), string(cmdOutput.Bytes()))
		m.errArray = append(m.errArray, errMessage)
	}
	mu.Unlock()
	wg.Done()
}
