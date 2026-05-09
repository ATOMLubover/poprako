package svc

import (
	"testing"
)

func TestParseLabelPlus(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
		wantErr bool
	}{
		{
			name: "real LabelPlus export without blank line after comment",
			content: "1,0\n" +
				"-\n" +
				"框内\n" +
				"框外\n" +
				"-\n" +
				"可使用 LabelPlus 脚本导入 psd 中\n" +
				">>>>>>>>[01.jpg]<<<<<<<<\n" +
				"----------------[1]----------------[0.854,0.117,2]\n" +
				"陪长着辟谷毛的\n" +
				"----------------[2]----------------[0.731,0.416,2]\n" +
				"公司前辈\n" +
				">>>>>>>>[02.jpg]<<<<<<<<\n",
			want: 2,
		},
		{
			name: "poprako self-export with blank line after comment",
			content: "1,0\n" +
				"-\n" +
				"框内\n" +
				"框外\n" +
				"-\n" +
				"Exported by PopRaKo Web\n" +
				"\n" +
				">>>>>>>>[page_1.jpg]<<<<<<<<\n" +
				"----------------[1]----------------[0.500,0.500,1]\n" +
				"hello\n",
			want: 1,
		},
		{
			name: "LabelPlus with 3 groups",
			content: "1,0\n" +
				"-\n" +
				"框内普通\n" +
				"框内心理\n" +
				"框外\n" +
				"-\n" +
				"my note\n" +
				">>>>>>>>[01.png]<<<<<<<<\n" +
				"----------------[1]----------------[0.5,0.5,1]\n" +
				"text\n",
			want: 1,
		},
		{
			name: "version 1 only (no comma)",
			content: "1\n" +
				"-\n" +
				"框内\n" +
				"-\n" +
				"note\n" +
				">>>>>>>>[01.png]<<<<<<<<\n" +
				"----------------[1]----------------[0.5,0.5,1]\n" +
				"text\n",
			want: 1,
		},
		{
			name:    "empty file",
			content: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := ChapterImportSvc{}
			pages, err := svc.ParseLabelPlus(tt.content)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.wantErr && len(pages) != tt.want {
				t.Fatalf("got %d pages, want %d", len(pages), tt.want)
			}
		})
	}
}
