package inspector_test

import (
	"os"
	"reflect"
	"testing"

	"degit/internal/template/inspector"
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

func TestInspector_Parse(t *testing.T) {
	tests := []struct {
		name           string
		templateString string
		want           []string
		wantErr        bool
	}{
		{
			name:           "empty template",
			templateString: "",
			want:           []string{},
			wantErr:        false,
		},
		{
			name:           "single variable",
			templateString: "Hello {{.Name}}",
			want:           []string{"Name"},
			wantErr:        false,
		},
		{
			name:           "multiple variables",
			templateString: "Hello {{.FirstName}} {{.LastName}}",
			want:           []string{"FirstName", "LastName"},
			wantErr:        false,
		},
		{
			name:           "invalid template",
			templateString: "Hello {{.Name",
			want:           nil,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := inspector.New()
			got, err := i.Parse(tt.templateString)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Parse() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("Parse() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInspector_ParseFiles(t *testing.T) {
	tests := []struct {
		name    string
		files   map[string]string
		want    []string
		wantErr bool
	}{
		{
			name: "single file with one variable",
			files: map[string]string{
				"file1.tmpl": "Hello {{.Name}}",
			},
			want:    []string{"Name"},
			wantErr: false,
		},
		{
			name: "multiple files with different variables",
			files: map[string]string{
				"file1.tmpl": "Hello {{.FirstName}}",
				"file2.tmpl": "Welcome {{.LastName}}",
			},
			want:    []string{"FirstName", "LastName"},
			wantErr: false,
		},
		{
			name: "file with complex template",
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
			want:    []string{"ShowHeader", "Title", "Items", "Name"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePaths []string
			// Create temporary files
			for _, content := range tt.files {
				path := createTempFile(t, content)
				filePaths = append(filePaths, path)
				// Cleanup after test
				defer os.Remove(path)
			}

			i := inspector.New()
			got, err := i.ParseFiles(filePaths)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseFiles() error = nil, wantErr %v", tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseFiles() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseFiles() got = %v, want %v", got, tt.want)
			}
		})
	}
}
