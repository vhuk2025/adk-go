// Copyright 2026 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package console

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ergochat/readline"
)

func TestReadConsoleInputDeletesWholeCJKCharacter(t *testing.T) {
	var output bytes.Buffer
	reader, err := readline.NewFromConfig(&readline.Config{
		Stdin:          strings.NewReader("中文\x7f\r"),
		Stdout:         &output,
		Stderr:         &output,
		FuncIsTerminal: func() bool { return true },
		FuncMakeRaw:    func() error { return nil },
		FuncExitRaw:    func() error { return nil },
		FuncGetSize:    func() (int, int) { return 80, 24 },
		FuncOnWidthChanged: func(func()) {
		},
	})
	if err != nil {
		t.Fatalf("readline.NewFromConfig() error = %v", err)
	}
	defer reader.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	inputChan := make(chan string, 1)
	readErrChan := make(chan error, 1)
	go readConsoleInput(ctx, reader, inputChan, readErrChan)

	select {
	case got := <-inputChan:
		if want := "中\n"; got != want {
			t.Errorf("readConsoleInput() = %q, want %q", got, want)
		}
	case err := <-readErrChan:
		t.Fatalf("readConsoleInput() error = %v", err)
	case <-time.After(time.Second):
		t.Fatal("readConsoleInput() timed out")
	}
}
