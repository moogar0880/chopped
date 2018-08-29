package chopped

import (
	"fmt"
	"os"
)

func ls(cwd string) (filenames []string, dirs []string) {
	d, _ := os.Open(cwd)

	filenames = make([]string, 0)
	dirs = make([]string, 0)
	for {
		fileInfo, err := d.Readdir(1)
		if err != nil {
			break
		}

		if fileInfo[0].IsDir() {
			dirs = append(dirs, cwd+"/"+fileInfo[0].Name())
		} else {
			filenames = append(filenames, cwd+"/"+fileInfo[0].Name())
		}
	}
	return filenames, dirs
}

func walk(files chan *string) {
	cwd, _ := os.Getwd()

	dirsToWalk := []string{cwd}

	for {
		cwd = dirsToWalk[0]
		fmt.Println(cwd)

		filenames, dirs := ls(cwd)

		for _, f := range filenames {
			files <- &f
		}

		fmt.Println(filenames)
		fmt.Println(dirs)

		dirsToWalk = append(dirsToWalk[1:], dirs...)

		if len(dirs) == 0 {
			break
		}

		fmt.Println("")
	}

	files <- nil
}
