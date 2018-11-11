package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func main() {
	if len(os.Args) < 1 {
		fmt.Println("missing audio url")
		os.Exit(1)
	}
	res, err := http.Get(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer res.Body.Close()

	d, err := mp3.NewDecoder(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()

	p, err := oto.NewPlayer(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer p.Close()

	if _, err := io.Copy(p, d); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
