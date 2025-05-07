package core

import (
	"fmt"
	"testing"
)

func TestMakeSSCS(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(i, "----------------------------------------")
		seq1 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGAAGAAAAAAAAGGAAAAAGGAG-AGAAGTAAGAAAGGCGCGGGACGGGCAGCGAGAAGGCGTGGAGTGGGTGAACAGCACCTTGCTCGAAGGA-GACGACACAGTAGAGAAAAAAAAACAAAAAAAAATATAATAAGAAAGAGAGAGAGAGGATGGGGGGGAGAGGGTGGGAGGGGAGGGGGACAGGGGTGTGGGTGGG--GTAAAGAGAAAAAAAAGACAGACAAGCGAGTAGAGCACACAAGGCTATGTAC"
		seq2 := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGGAAGGAGGAACAGGAGAAGGAGAGCGGGGCAGGCGCGACTGAGGGTTTGTTATTTGATGCTGATCTGTCTCCTTGTCGCTCTATGTCTGAATAAAACCTAATAAAAAAAAAAAAAATAAAAATATAA-AATGAAGAATGAGGGAGG--GGGTGGGAGGGGATGGGAGGGAGTGGAGAGAGGGGCAGAGGTAAGATGTATAAATACACATATAGCAAATGACACGATGTATCGTAACAACTAATTT--"

		qua1 := "::CCFCFGGFFCFCCE@:BCED:@F:7>+3=+3+++,,,338,***++,,,,,++33 8*8,,3,,3,,,2********/**/********/,*******2,,*28/*,0+;*0***3++,+**,*** 2*******2*0++2+<+02/:*8?//C***2/***<+2+3+0+3+*+*2*/*****18*+0/1/**0*855*/*9))/*/)-*2):)/))*20/97))*0*)*:*  ((*))**(,*.2)7()**34(2*().(/*(-(-=*+).((*/.(((//*.6)"
		qua2 := ":CCFGGDGGGDDGGGG@FE>C@FE@:7:++3+********++,8,+2+*+8,,,**++2,*,*******1*****/*2**2*,,**/,**0,,,3,,+3,+,+,3,3,+30++0++2+*/1**+++++3+++++++0+**3++32+++/***:/**:*3300<:C+3@ +:+*+++*++0++*****  ))))/)7(*()0**)*,*1.8*.**-2)05-*))))((0*()****))..).*.)..))0)-**)-)*(.---))(-,-()).(42(,))(/.))))"
		families := []string{seq1, seq2}
		qualities := []string{qua1, qua2}
		sscs, _ := MakeSSCS(families, qualities)
		fmt.Println(sscs)
	}

}
