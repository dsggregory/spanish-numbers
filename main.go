package main

import (
	"context"
	"flag"
	"fmt"
	"spanish-numbers/pkg/langpractice"
	"spanish-numbers/pkg/term"
	"strconv"
	"time"
)

const prompt = "> "

func main() {
	var min, max = 100, 1000
	var autoNext bool

	flag.IntVar(&min, "min", min, "min")
	flag.IntVar(&max, "max", max, "max")
	flag.BoolVar(&autoNext, "noan", false, "no auto next number")
	flag.Parse()

	c := langpractice.NewClient(min, max)
	c.AutoNext = !autoNext

	ctx, cancel := context.WithCancel(context.Background())
	kpr, err := term.NewKeypressReader(cancel)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer kpr.Reset()

	printLegend(nil)

loop:
	for {
		var answer []byte
		var begin time.Time

		resp, err := c.RequestNumber()
		if err != nil {
			fmt.Println(err)
			return
		}

		if err = c.PlayResponse(resp); err != nil {
			fmt.Println(err)
			return
		}

		nstr := fmt.Sprintf("%d", resp.Number)

		begin = time.Now()
		fmt.Printf("%s", prompt)

	answerloop:
		for {
			select {
			case <-ctx.Done():
				break loop
			case b := <-kpr.KeyEvent():
				switch b {
				case 'r':
					if err = c.PlayResponse(resp); err != nil {
						fmt.Println(err)
						return
					}
				case 'q':
					// upstream keypress will send a cancel()
				case 'n', 'x':
					// show answer and skip to next
					fmt.Printf("!x! \"%d\"", resp.Number)
					break answerloop
				case '?':
					// hint
					fmt.Printf("!?! \"%s\"\n%s", resp.Target.Written, prompt)
				default:
					if b >= byte('0') && b <= byte('9') {
						l := len(answer)
						if l == 0 {
							begin = time.Now() // start timer after first number keypress
						}
						// incorrect digit in this spot
						if nstr[l] != b {
							c.Beep()
							continue
						}
						// correct digit in this spot
						fmt.Printf("%s", string(b))
						answer = append(answer, byte(b))
						v, err := strconv.Atoi(string(answer))
						if err == nil && v == resp.Number {
							break answerloop
						}
					} else {
						printLegend(answer)
					}
				}
			}
		}

		// answered successfully
		fmt.Println(" ==>", time.Now().Sub(begin), "-", resp.Target.Written)
		if c.AutoNext {
			time.Sleep(time.Millisecond * 1000)
		} else {
			fmt.Printf("next?")
			b := <-kpr.KeyEvent() // wait on any key press
			fmt.Println("")
			if b == 'q' {
				break loop
			}
		}
	}
}

func printLegend(curAnswer []byte) {
	fmt.Printf(`
Legend:
  r: replay
  ?: hint
  l: print this legend
  n: next number
  q: quit
  0-9: digit
`)
	if len(curAnswer) > 0 {
		fmt.Printf("%s%s", prompt, curAnswer)
	}
}
