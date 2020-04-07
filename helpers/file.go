package helpers

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func ExpandPath(str string, expansions ...string) string {
	r := strings.NewReplacer(expansions...)
	str = r.Replace(str)

	return str
}

func Exists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		// exists
		return true
	} else if os.IsNotExist(err) {
		// does not exist
		return false
	} else {
		// something else happened
		return false
	}
}

func Copy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func Download(url string, destination string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
