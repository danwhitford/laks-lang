package laks

import (
	"bufio"
	"bytes"
	"embed"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed programs
var f embed.FS

func TestPrograms(t *testing.T) {
	programs, err := f.ReadDir("programs")
	if err != nil {
		t.Fatalf("could not read programs: %v", err)
	}

	for _, program := range programs {
		if program.IsDir() {
			continue
		}

		t.Run(program.Name(), func(tt *testing.T) {
			b, err := f.ReadFile("programs/" + program.Name())
			if err != nil {
				tt.Fatalf("could not read file %s: %v", program.Name(), err)
			}

			var expectedBuf strings.Builder
			scanner := bufio.NewScanner(bytes.NewReader(b))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "# ") {
					expected := strings.TrimPrefix(line, "# ")
					expectedBuf.WriteString(expected)
					expectedBuf.WriteByte('\n')
				}
			}

			outputBuf := &bytes.Buffer{}
			err = RunBytes(b, outputBuf)
			if err != nil {
				tt.Fatalf("could not run program %s: %v", program.Name(), err)
			}

			if cmp.Diff(expectedBuf.String(), outputBuf.String()) != "" {
				tt.Errorf("output mismatch for program %s:\n%s", program.Name(), cmp.Diff(expectedBuf.String(), outputBuf.String()))
			}
		})
	}
}
