package internal

import (
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"runtime"
	"testing"
)

func makeFunctionRunOnRootFolder() {
	_, filename, _, _ := runtime.Caller(0)
	// The ".." may change depending on you folder structure
	dir := path.Join(path.Dir(filename), "..")
	err := os.Chdir(dir)
	if err != nil {
		panic(err)
	}
}

func TestSliceIntersect(t *testing.T) {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error("Error loading .env file")
	}
	makeFunctionRunOnRootFolder()
	s := NewScreen()

	testCases := []struct {
		name     string
		imgPath  string
		expected string
	}{
		{
			name:     "1",
			imgPath:  "\\static\\screens\\tmp.png",
			expected: "7900 +7?\r\n",
		},
		{
			name:     "2",
			imgPath:  "\\static\\screens\\test11.png",
			expected: "350 +8?\r\n",
		},
	}
	for _, tc := range testCases {
		out, err := s.FindText(tc.imgPath)
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, out)
	}
}
