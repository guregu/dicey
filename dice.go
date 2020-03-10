package dicey

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// Dice is a parsed dice formula.
type Dice struct {
	ops []operation
}

// Roll returns the sum of this dice's random rolls.
// It uses math/rand's global random source.
func (d Dice) Roll() int {
	total := 0
	for _, op := range d.ops {
		total += op.roll()
	}
	return total
}

// Max returns the sum of every dice's highest possible value.
func (d Dice) Max() int {
	total := 0
	for _, op := range d.ops {
		total += op.max()
	}
	return total
}

// Min returns the sum of every dice's lowest possible value.
func (d Dice) Min() int {
	total := 0
	for _, op := range d.ops {
		total += op.min()
	}
	return total
}

// Parse creates a new Dice with a D&D style formula.
// It can contain any number of dice (such as 3d6 for 3 six-sided die) or
// bonuses combined with + for addition or - for subtraction.
func Parse(formula string) (Dice, error) {
	d := Dice{}
	var err error
	l := lex(formula)
	add := true
loop:
	for {
		item := l.nextItem()
		switch item.Type {
		case itemError:
			err = fmt.Errorf("dicey: expression lex error: %s at position %d", item.Val, item.Pos)
			break loop
		case itemEOF:
			break loop
		case itemBonus:
			n, err := strconv.Atoi(item.Val)
			if err != nil {
				return Dice{}, err
			}
			if !add {
				n = -n
			}
			d.ops = append(d.ops, bonusOp(n))
		case itemDice:
			split := strings.Split(item.Val, "d")
			dice, err := strconv.Atoi(split[0])
			if err != nil {
				return Dice{}, err
			}
			sides, err := strconv.Atoi(split[1])
			if err != nil {
				return Dice{}, err
			}
			d.ops = append(d.ops, diceOp{
				dice:  dice,
				sides: sides,
				add:   add,
			})
		case itemAdd:
			add = true
		case itemSub:
			add = false
		}
	}
	return d, err
}

// MustParse returns a Dice for the given formula or panics. See Parse for more detail.
func MustParse(formula string) Dice {
	d, err := Parse(formula)
	if err != nil {
		panic(err)
	}
	return d
}

type operation interface {
	roll() int
	max() int
	min() int
}

type bonusOp int

func (bo bonusOp) roll() int {
	return int(bo)
}

func (bo bonusOp) max() int {
	return int(bo)
}

func (bo bonusOp) min() int {
	return int(bo)
}

type diceOp struct {
	dice  int
	sides int
	add   bool
}

func (do diceOp) roll() int {
	total := 0
	for i := 0; i < do.dice; i++ {
		total += rand.Intn(do.sides) + 1
	}
	if !do.add {
		total = -total
	}
	return total
}

func (do diceOp) max() int {
	n := do.dice * do.sides
	if !do.add {
		return -n
	}
	return n
}

func (do diceOp) min() int {
	n := do.dice
	if !do.add {
		return -n
	}
	return n
}
