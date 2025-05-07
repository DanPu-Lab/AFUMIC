package core

import (
	"bufio"
	"correct-go/src/utils"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func MapNamesToBarcodes(barcodesPath string) map[int]string {
	res := make(map[int]string)
	// 打开文件
	f, err := os.Open(barcodesPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	// 读取文件，记录到map中
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		numLine := scanner.Text()
		num, _ := strconv.Atoi(numLine[1:])
		scanner.Scan()
		barcode := scanner.Text()
		res[num] = barcode
	}
	return res
}

func FilterAlignment(namesToBarcodes map[int]string, samPath, lostBarcodesPath string,
	posThres, mapqThres, distThres int) []AlignmentData {
	f, err := os.Open(samPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var lostBarcodes *os.File
	if len(lostBarcodesPath) != 0 {
		lostBarcodes, err = os.Create(lostBarcodesPath)
		if err != nil {
			panic(err)
		}
		defer lostBarcodes.Close()
	}

	res := make([]AlignmentData, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")

		data, ok := getBarcodeAlignment(fields, posThres, mapqThres, distThres)
		if !ok {
			if data != nil && lostBarcodes != nil {
				rseq := namesToBarcodes[data.RName]
				qseq := namesToBarcodes[data.QName]

				line := fmt.Sprintf("%s->%s\n", rseq, qseq)
				_, err = lostBarcodes.WriteString(line)
				if err != nil {
					panic(err)
				}
			}
			continue
		}

		res = append(res, *data)
	}
	return res
}

func getBarcodeAlignment(fields []string, posThres, mapqThres, distThres int) (*AlignmentData, bool) {
	aln, err := NewAlignment(fields)
	if err != nil {
		return nil, false
	}

	if len(aln.RName) == 0 || aln.RName == "*" {
		return nil, false
	}
	rnameFields := strings.Split(aln.RName, ":")
	var reversed bool
	var rnameStr string
	if len(rnameFields) == 2 && rnameFields[1] == "rev" {
		reversed = true
		rnameStr = rnameFields[0]
	} else {
		reversed = false
		rnameStr = aln.RName
	}

	qname, _ := strconv.Atoi(aln.QName)
	rname, _ := strconv.Atoi(rnameStr)

	if qname == rname ||
		aln.Unmapped() ||
		utils.Abs(aln.Pos-1) > posThres ||
		aln.Mapq < mapqThres {
		return nil, false
	}

	nm, ok := aln.Tags["NM"]
	if !ok && strings.Contains(aln.Cigar, "N") {
		return &AlignmentData{
			RName:    rname,
			QName:    qname,
			Reversed: reversed,
		}, false
	}
	if nm > distThres {
		return &AlignmentData{
			RName:    rname,
			QName:    qname,
			Reversed: reversed,
		}, false
	}

	return &AlignmentData{
		RName:    rname,
		QName:    qname,
		Reversed: reversed,
	}, true
}

func ReadAlignments(alignments []AlignmentData,
	namesToBarcodes map[int]string) (Graph, map[string]struct{}, int) {
	numGoodAlignments := 0
	reversedBarcodes := make(map[string]struct{})
	graph := NewGraph()
	for _, aln := range alignments {
		numGoodAlignments++
		rseq := namesToBarcodes[aln.RName]
		qseq := namesToBarcodes[aln.QName]
		if aln.Reversed {
			reversedBarcodes[rseq] = struct{}{}
			reversedBarcodes[qseq] = struct{}{}
		}
		graph.AddNode(rseq)
		graph.AddNode(qseq)
		graph.AddEdge(rseq, qseq)
	}
	return graph, reversedBarcodes, numGoodAlignments
}

func GetFamilyCounts(familyPath string, limit int, checkIds bool) (map[string]map[string]int, int) {
	f, err := os.Open(familyPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	familyCounts := make(map[string]map[string]int)
	lastBarcode := ""
	thisFamilyCounts := make(map[string]int)
	readPairs := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		readPairs++
		if limit != 0 && readPairs > limit {
			break
		}
		fields := strings.Split(scanner.Text(), "\t")
		if checkIds {
		}
		barcode := fields[0]
		order := fields[1]
		if barcode != lastBarcode {
			if len(thisFamilyCounts) != 0 {
				thisFamilyCounts["all"] = thisFamilyCounts["ab"] + thisFamilyCounts["ba"]
			}
			familyCounts[barcode] = thisFamilyCounts
			thisFamilyCounts = make(map[string]int)
			lastBarcode = barcode
		}
		thisFamilyCounts[order]++
	}
	thisFamilyCounts["all"] = thisFamilyCounts["ab"] + thisFamilyCounts["ba"]
	familyCounts[lastBarcode] = thisFamilyCounts
	return familyCounts, readPairs
}

func MakeCorrectionTable(metaGraph Graph, familyCounts map[string]map[string]int) map[string]string {
	corrections := make(map[string]string)
	for _, nodes := range metaGraph.ConnectedComponents() {
		sort.Slice(nodes, func(i, j int) bool {
			return familyCounts[nodes[i]]["all"] < familyCounts[nodes[j]]["all"]
		})
		correct := nodes[0]
		for _, barcode := range nodes {
			if barcode != correct {
				corrections[barcode] = correct
			}
		}
	}
	return corrections
}

func GenerateCorrectedOutput(familyPath string, corrections map[string]string,
	reversedBarcodes map[string]struct{}, output chan<- string) {
	barcodeLast := ""

	f, err := os.Open(familyPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), "\t")
		rawBarcode := fields[0]
		order := fields[1]

		if rawBarcode != barcodeLast {
			barcodeLast = rawBarcode
		}

		var correctOrder string
		var correctBarcode string
		if _, ok := corrections[rawBarcode]; ok {
			correctBarcode = corrections[rawBarcode]

			_, rawIsReversed := reversedBarcodes[rawBarcode]
			_, correctIsReversed := reversedBarcodes[correctBarcode]
			if (rawIsReversed || correctIsReversed) &&
				IsAilgnmentReversed(rawBarcode, correctBarcode) {
				if order == "ab" {
					correctOrder = "ba"
				} else {
					correctOrder = "ab"
				}
			} else {
				correctOrder = order
			}
		} else {
			correctOrder = order
			correctBarcode = rawBarcode
		}
		fields[0] = correctBarcode
		fields[1] = correctOrder
		line := strings.Join(fields, "\t")
		output <- line
	}
}

type isAilgnmentReversedCache struct {
	barcode1, barcode2 string
}

var functionCache = make(map[isAilgnmentReversedCache]bool)

func IsAilgnmentReversed(barcode1, barcode2 string) bool {
	cache := isAilgnmentReversedCache{barcode1, barcode2}
	if _, ok := functionCache[cache]; ok {
		return functionCache[cache]
	}
	half := len(barcode2) / 2
	barcode2Rev := barcode2[:half] + barcode2[half:]
	_, _, fwdScore := SmithWaterman(barcode1, barcode2)
	_, _, revScore := SmithWaterman(barcode1, barcode2Rev)
	result := fwdScore > revScore
	functionCache[cache] = result
	return result
}
