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
				case 'n':
					// show answer and skip to next
					fmt.Printf("!x! \"%d\"", resp.Number)
					break answerloop
				default:
					if b >= byte('0') && b <= byte('9') {
						// incorrect digit in this spot
						if nstr[len(answer)] != b {
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
						fmt.Printf("\nLegend:\n  r: replay\n  q: quit\n  n: next number\n  0-9: digit\n%s", answer)
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
