	"os"
		"unterminatedQuotes": {
			Input: `"dir/file.txt`, Err: true,
		},

func TestParseGitHeaderData(t *testing.T) {
	tests := map[string]struct {
		InputFile   *File
		Line        string
		DefaultName string

		OutputFile *File
		End        bool
		Err        bool
	}{
		"fragementEndsParsing": {
			Line: "@@ -12,3 +12,2 @@\n",
			End:  true,
		},
		"unknownEndsParsing": {
			Line: "GIT binary file\n",
			End:  true,
		},
		"oldFileName": {
			Line: "--- a/dir/file.txt\n",
			OutputFile: &File{
				OldName: "dir/file.txt",
			},
		},
		"newFileName": {
			Line: "+++ b/dir/file.txt\n",
			OutputFile: &File{
				NewName: "dir/file.txt",
			},
		},
		"oldMode": {
			Line: "old mode 100644\n",
			OutputFile: &File{
				OldMode: os.FileMode(0100644),
			},
		},
		"invalidOldMode": {
			Line: "old mode rw\n",
			Err:  true,
		},
		"newMode": {
			Line: "new mode 100755\n",
			OutputFile: &File{
				NewMode: os.FileMode(0100755),
			},
		},
		"invalidNewMode": {
			Line: "new mode rwx\n",
			Err:  true,
		},
		"deletedFileMode": {
			Line:        "deleted file mode 100644\n",
			DefaultName: "dir/file.txt",
			OutputFile: &File{
				OldName:  "dir/file.txt",
				OldMode:  os.FileMode(0100644),
				IsDelete: true,
			},
		},
		"newFileMode": {
			Line:        "new file mode 100755\n",
			DefaultName: "dir/file.txt",
			OutputFile: &File{
				NewName: "dir/file.txt",
				NewMode: os.FileMode(0100755),
				IsNew:   true,
			},
		},
		"copyFrom": {
			Line: "copy from dir/file.txt\n",
			OutputFile: &File{
				OldName: "dir/file.txt",
				IsCopy:  true,
			},
		},
		"copyTo": {
			Line: "copy to dir/file.txt\n",
			OutputFile: &File{
				NewName: "dir/file.txt",
				IsCopy:  true,
			},
		},
		"renameFrom": {
			Line: "rename from dir/file.txt\n",
			OutputFile: &File{
				OldName:  "dir/file.txt",
				IsRename: true,
			},
		},
		"renameTo": {
			Line: "rename to dir/file.txt\n",
			OutputFile: &File{
				NewName:  "dir/file.txt",
				IsRename: true,
			},
		},
		"similarityIndex": {
			Line: "similarity index 88\n",
			OutputFile: &File{
				Score: 88,
			},
		},
		"similarityIndexTooBig": {
			Line: "similarity index 9001\n",
			OutputFile: &File{
				Score: 0,
			},
		},
		"indexFullSHA1AndMode": {
			Line: "index 79c6d7f7b7e76c75b3d238f12fb1323f2333ba14..04fab916d8f938173cbb8b93469855f0e838f098 100644\n",
			OutputFile: &File{
				OldOID:  "79c6d7f7b7e76c75b3d238f12fb1323f2333ba14",
				NewOID:  "04fab916d8f938173cbb8b93469855f0e838f098",
				OldMode: os.FileMode(0100644),
			},
		},
		"indexFullSHA1NoMode": {
			Line: "index 79c6d7f7b7e76c75b3d238f12fb1323f2333ba14..04fab916d8f938173cbb8b93469855f0e838f098\n",
			OutputFile: &File{
				OldOID: "79c6d7f7b7e76c75b3d238f12fb1323f2333ba14",
				NewOID: "04fab916d8f938173cbb8b93469855f0e838f098",
			},
		},
		"indexAbbrevSHA1AndMode": {
			Line: "index 79c6d7..04fab9 100644\n",
			OutputFile: &File{
				OldOID:  "79c6d7",
				NewOID:  "04fab9",
				OldMode: os.FileMode(0100644),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var f File
			if test.InputFile != nil {
				f = *test.InputFile
			}

			end, err := parseGitHeaderData(&f, test.Line, test.DefaultName)
			if test.Err {
				if err == nil {
					t.Fatalf("expected error parsing header data, but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error parsing header data: %v", err)
			}

			if test.OutputFile != nil && !reflect.DeepEqual(test.OutputFile, &f) {
				t.Errorf("incorrect output:\nexpected: %+v\nactual: %+v", test.OutputFile, &f)
			}
			if end != test.End {
				t.Errorf("incorrect end state, expected %t, actual %t", test.End, end)
			}
		})
	}
}