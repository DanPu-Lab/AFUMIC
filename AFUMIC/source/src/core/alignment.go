package core

import (
	"fmt"
	"strconv"
	"strings"
)

type Alignment struct {
	QName string
	Flag  int
	RName string
	Pos   int
	Mapq  int
	Cigar string
	RNext string
	PNext int
	TLen  int
	Seq   string
	Qual  string
	Tags  map[string]int
}

func NewAlignment(fields []string) (*Alignment, error) {
	if len(fields) < 11 {
		return nil, fmt.Errorf("not enough fields: %d", len(fields))
	}
	a := &Alignment{
		QName: fields[0],
		RName: fields[2],
		Cigar: fields[5],
		RNext: fields[6],
		Seq:   fields[9],
		Qual:  fields[10],
	}
	a.Flag, _ = strconv.Atoi(fields[1])
	a.Pos, _ = strconv.Atoi(fields[3])
	a.Mapq, _ = strconv.Atoi(fields[4])
	a.PNext, _ = strconv.Atoi(fields[7])
	a.TLen, _ = strconv.Atoi(fields[8])
	a.Tags = make(map[string]int)
	for i := 11; i < len(fields); i++ {
		tagFields := strings.Split(fields[i], ":")
		tagName := tagFields[0]
		tagValue, _ := strconv.Atoi(tagFields[2])
		a.Tags[tagName] = tagValue
	}
	return a, nil
}

func (this *Alignment) Unmapped() bool {
	unmappedFlags := 4
	return this.Flag&unmappedFlags != 0
}

type AlignmentData struct {
	QName    int
	RName    int
	Reversed bool
}
