package core

import (
	"consensus-go-lib/src/util"
	"math"
)

// parameters
var (
	OffsetRange              = 4
	AnchorRegion             = 8
	MinSequenceSize          = 2*OffsetRange + 2*AnchorRegion
	MaxMms                   = 4
	MaxDroppedSequencesRatio = 0.3
)

type ReadWithOffset struct {
	Read       string
	Qual       string
	BestOffset int
	L          int
	X          int
	Y          int
	From       int
	To         int
}

func newReadWithOffset(read, qual string, offset int) *ReadWithOffset {
	mid := len(read) / 2
	return &ReadWithOffset{
		Read:       read,
		Qual:       qual,
		BestOffset: offset,
		L:          len(read),
		X:          mid - offset,
		Y:          len(read) - mid + offset,
		From:       0,
		To:         0,
	}
}

func GetCoreSequence(sequence string, offset int) string {
	mid := len(sequence) / 2
	left := mid - offset - AnchorRegion
	right := mid - offset + AnchorRegion + 1
	newSequence := sequence[left:right]
	return newSequence
}

func MakeSSCS(families []string, quals []string) (string, string) {
	// 1. 计数
	coreSeqMap := make(map[string]int, len(families)*8)
	coreSeqList := make([]string, 0, len(families)*8)
	coreSeqDataList := make([][2]int, 0, len(families)*8)
	for _, sequence := range families {
		if len(sequence) > MinSequenceSize {
			for offset := -OffsetRange; offset <= OffsetRange; offset++ {
				coreSeq := GetCoreSequence(sequence, offset)
				coreSeqIndex, ok := coreSeqMap[coreSeq]
				if !ok {
					coreSeqIndex = len(coreSeqList)
					coreSeqList = append(coreSeqList, coreSeq)
					coreSeqDataList = append(coreSeqDataList, [2]int{0, 0})
				}
				coreSeqDataList[coreSeqIndex][0]++
				coreSeqDataList[coreSeqIndex][1] += len(sequence)
				coreSeqMap[coreSeq] = coreSeqIndex
			}
		}
	}

	// 2. 找到最大的
	bestCoreSeq := ""
	coreSeqData := [2]int{0, 0}
	for _, seq := range coreSeqList {
		i := coreSeqMap[seq]
		if coreSeqDataList[i][0] > coreSeqData[0] ||
			(coreSeqDataList[i][0] == coreSeqData[0] && coreSeqDataList[i][1] < coreSeqData[1]) {
			bestCoreSeq = seq
			coreSeqData = coreSeqDataList[i]
		}
	}
	// 用完清空
	coreSeqMap = nil
	coreSeqList = nil
	coreSeqDataList = nil

	// 3. 对比匹配
	assembledSequences := make([]*ReadWithOffset, 0, len(families))
	for i, seq := range families {
		if len(seq) > MinSequenceSize {
			bestOffset := 0
			bestOffsetMMs := AnchorRegion
			for offset := -OffsetRange; offset <= OffsetRange; offset++ {
				offsetMMs := 0
				coreSeq := GetCoreSequence(seq, offset)
				if coreSeq == bestCoreSeq {
					bestOffset = offset
					bestOffsetMMs = 0
				} else {
					for j := 0; j < len(coreSeq); j++ {
						if coreSeq[j] != bestCoreSeq[j] {
							offsetMMs++
						}
					}
					if offsetMMs < bestOffsetMMs {
						bestOffset = offset
						bestOffsetMMs = offsetMMs
					}
				}
			}

			if bestOffsetMMs <= MaxMms {
				qual := quals[i]
				readWithOffset := newReadWithOffset(seq, qual, bestOffset)
				assembledSequences = append(assembledSequences, readWithOffset)
			}
		}
	}

	n := len(assembledSequences)
	droppedSeqRatio := 1.0 - float64(n)/float64(len(families))
	if droppedSeqRatio > MaxDroppedSequencesRatio {
		return "", ""
	}

	// 4. 计算pwm
	pwm := FillPwnAndRecomputeOffsets(assembledSequences)

	// 5. 计算consensus
	sequence, quality, trimmedBasesRatio := ConstructConsensus(pwm, n)

	if trimmedBasesRatio > 0.3 {
		return "", ""
	}
	return sequence, quality
}

func FillPwnAndRecomputeOffsets(assembledSequences []*ReadWithOffset) [5][]int {
	var (
		X = 0
		Y = 0
	)
	for _, readWithOffset := range assembledSequences {
		X = util.MaxInt(X, readWithOffset.X)
		Y = util.MaxInt(Y, readWithOffset.Y)
	}
	pwm := [5][]int{}
	for i := 0; i < 5; i++ {
		pwm[i] = make([]int, X+Y)
	}

	for _, readWithOffset := range assembledSequences {
		xDelta := X - readWithOffset.X
		yDelta := Y - readWithOffset.Y

		readWithOffset.From = util.MaxInt(-xDelta, 0)
		readWithOffset.To = readWithOffset.L + util.MinInt(yDelta, 0)

		for k := readWithOffset.From; k < readWithOffset.To; k++ {
			pwmPos := xDelta + k
			code := GetCode(readWithOffset.Read[k])
			pwm[code][pwmPos]++
		}
	}
	return pwm
}

func ConstructConsensus(pwm [5][]int, n int) (string, string, float64) {
	pwmLen := len(pwm[0])
	consensusSequence := make([]byte, pwmLen)
	consensusQuality := make([]byte, pwmLen)
	goodSeqStart := 0

	for k := 0; k < pwmLen; k++ {
		mostFreqLetter := 0
		maxLetterFreq := 0
		for l := 0; l < 4; l++ {
			freq := pwm[l][k]
			if maxLetterFreq < freq {
				maxLetterFreq = freq
				mostFreqLetter = l
			}
		}
		consensusSequence[k] = GetSymbol(mostFreqLetter)

		cqs := util.MaxInt(2, int(math.Min(float64(40),
			40*((float64(maxLetterFreq)/float64(n)-0.25)/0.75))))
		consensusQuality[k] = byte(cqs + 33)

		if cqs < 10 && goodSeqStart == k {
			goodSeqStart++
		}
	}
	goodSeqEnd := pwmLen
	for goodSeqEnd >= goodSeqStart {
		if consensusQuality[goodSeqEnd-1] > 10 {
			break
		}
		goodSeqEnd--
	}
	trimmedBasesRatio := float64(pwmLen-goodSeqEnd+goodSeqStart) / float64(pwmLen)
	return string(consensusSequence), string(consensusQuality), trimmedBasesRatio
}
