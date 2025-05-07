package core

import (
	"correct-go/src/utils"
	"math"
	"strings"
)

type Entry struct {
	score float64
	prev  [2]int
}

const (
	MATCH    float64 = 2.0
	MISMATCH float64 = -0.5
	GAP      float64 = -1.0
)

func SmithWaterman(seq1, seq2 string) (string, string, float64) {
	m := len(seq1) + 1
	n := len(seq2) + 1
	S := make([][]Entry, m)
	for i := 0; i < m; i++ {
		S[i] = make([]Entry, n)
	}

	for i := 1; i < m; i++ {
		S[i][0].prev[0] = i - 1
	}
	for j := 1; j < n; j++ {
		S[0][j].prev[1] = j - 1
	}

	for i := 1; i < len(seq1); i++ {
		for j := 1; j < len(seq2); j++ {
			nw_score := MISMATCH
			if seq1[i-1] == seq2[j-1] {
				nw_score = MATCH
			}

			S[i][j].score = math.SmallestNonzeroFloat64
			S[i][j].prev[0] = 0
			S[i][j].prev[1] = 0

			for k := 0; k < 2; k++ {
				for l := 0; l < 2; l++ {
					val := 0.0
					if k == 0 && l == 0 {
						continue
					} else if k > 0 && l > 0 {
						val = nw_score
					} else if k > 0 || l > 0 {
						if (i == len(seq1) && k == 0) || (j == len(seq2) && l == 0) {
							val = 0.0
						} else {
							val = GAP
						}
					}

					val += S[i-k][j-k].score
					if val > S[i][j].score {
						S[i][j].score = val
						S[i][j].prev[0] = i - k
						S[i][j].prev[1] = j - l
					}
				}
			}
		}
	}
	return TraceBack(seq1, seq2, S)
}

func TraceBack(seq1, seq2 string, S [][]Entry) (string, string, float64) {
	i := len(S) - 1
	j := len(S[0]) - 1
	k := 0
	a := make([]byte, i+j+3)
	b := make([]byte, i+j+3)
	score := 0.0

	max := math.SmallestNonzeroFloat64
	for l := 0; l < len(S); l++ {
		for m := 0; m < len(S[0]); m++ {
			if S[l][m].score > max {
				i = l
				j = m
				max = S[i][j].score
			}
		}
	}

	for S[i][j].prev[0] != 0 && S[i][j].prev[1] != 0 {
		for i > 0 || j > 0 {
			new_i := S[i][j].prev[0]
			new_j := S[i][j].prev[1]

			if new_i < i {
				a[k] = seq1[i-1]
			} else {
				a[k] = '-'
			}

			if S[i][j].score > score {
				score = S[i][j].score
			}

			if new_j < j {
				b[k] = seq2[j-1]
			} else {
				b[k] = '-'
			}

			k++
			i = new_i
			j = new_j
		}
	}
	utils.Reverse(a)
	utils.Reverse(b)
	result1 := strings.TrimFunc(string(a), func(r rune) bool {
		return r <= 3
	})
	result2 := strings.TrimFunc(string(b), func(r rune) bool {
		return r <= 3
	})
	return result1, result2, score
}
