package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func main() {
	dir := "./bin"
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		info, err := f.Info()
		if err == nil {
			if info.Mode()&0111 == 0 {
				log.Printf("File ignored (non exec) : %s", f.Name())
				continue
			}
		}

		fileName := f.Name()
		fullPath := filepath.Join(dir, fileName)

		cmd := exec.Command(fullPath)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		wg.Add(1)

		go func(name string, c *exec.Cmd) {
			defer wg.Done()

			if err := c.Run(); err != nil {
				log.Printf("[%s] Error : %v", name, err)
				return
			}
			log.Printf("[%s] Finished", name)
		}(fileName, cmd)
	}

	wg.Wait()
	log.Println("Tous les exécutables sont terminés.")

}
