package grab

import "testing"

func TestSanitize(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Screenshot 2026-04-05.png", "screenshot-2026-04-05.png"},
		{"My   Document (1).pdf", "my-document-1.pdf"},
		{"ALLCAPS.TXT", "allcaps.txt"},
		{"already-clean.md", "already-clean.md"},
		{"spaces  and---hyphens.tar.gz", "spaces-and-hyphens.tar.gz"},
		{"...hidden", "...hidden"},
		{"café résumé.doc", "caf-r-sum.doc"},
		{"file with (parens) [brackets].jpg", "file-with-parens-brackets.jpg"},
		{".dotfile", "file.dotfile"},
		{"名前.txt", "file.txt"},
		{"", "file"},
		{"----.pdf", "file.pdf"},
		{"hello world", "hello-world"},
		{"  leading trailing  .md", "leading-trailing.md"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := Sanitize(tt.input)
			if got != tt.want {
				t.Errorf("Sanitize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
