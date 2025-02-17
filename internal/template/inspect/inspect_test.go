package inspect_test

import (
	"os"
	"reflect"
	"testing"

	"degit/internal/command"
	"degit/internal/template/inspect"
)

func createTempFile(t *testing.T, content string) string {
	t.Helper()
	tmpfile, err := os.CreateTemp("", "template-*.tmpl")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	return tmpfile.Name()
}

func TestInspectCommand_Execute(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		contents []string
		want     []string
		wantErr  bool
	}{
		{
			name:     "empty command",
			files:    nil,
			contents: nil,
			want:     []string{},
			wantErr:  false,
		},
		{
			name: "single file inspection",
			files: map[string]string{
				"file1.tmpl": "Hello {{.Name}}",
			},
			contents: nil,
			want:     []string{"Name"},
			wantErr:  false,
		},
		{
			name: "multiple files inspection",
			files: map[string]string{
				"file1.tmpl": "Hello {{.FirstName}}",
				"file2.tmpl": "Welcome {{.LastName}}",
			},
			contents: nil,
			want:     []string{"FirstName", "LastName"},
			wantErr:  false,
		},
		{
			name:  "single content inspection",
			files: nil,
			contents: []string{
				"Hello {{.Name}}",
			},
			want:    []string{"Name"},
			wantErr: false,
		},
		{
			name: "combined files and contents",
			files: map[string]string{
				"file1.tmpl": "Hello {{.FirstName}}",
			},
			contents: []string{
				"Welcome {{.LastName}}",
			},
			want:    []string{"FirstName", "LastName"},
			wantErr: false,
		},
		{
			name: "complex template",
			files: map[string]string{
				"file1.tmpl": `
					{{if .ShowHeader}}
						Header: {{.Title}}
					{{end}}
					{{range .Items}}
						- {{.Name}}
					{{end}}
				`,
			},
			contents: nil,
			want:     []string{"ShowHeader", "Title", "Items", "Name"},
			wantErr:  false,
		},
		{
			name: "invalid template file",
			files: map[string]string{
				"file1.tmpl": "Hello {{.Name",
			},
			contents: nil,
			want:     nil,
			wantErr:  true,
		},
		{
			name:  "invalid template content",
			files: nil,
			contents: []string{
				"Hello {{.Name",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePaths []string
			// Create temporary files
			for _, content := range tt.files {
				path := createTempFile(t, content)
				filePaths = append(filePaths, path)
				defer os.Remove(path)
			}

			cmd := inspect.New(filePaths, tt.contents)
			ctx := command.NewContext()
			err := cmd.Execute(ctx)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Execute() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Execute() unexpected error = %v", err)
				return
			}

			got, ok := ctx.Get(inspect.ContextKeyExecuteResult).([]string)
			if !ok {
				t.Errorf("Execute() result not found in context or wrong type")
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Execute() got = %v, want %v", got, tt.want)
			}
		})
	}
}
