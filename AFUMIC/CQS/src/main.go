package main

import (
	"consensus-go-lib/src/core"
	"consensus-go-lib/src/reader"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Key struct {
	Barcode string
	Order   string
	Meta    string
}

type Mig struct {
	families  []string
	qualities []string
}

type Result struct {
	Sequence string
	Quality  string
	Key      Key
	NReads   string
}

var (
	filePath      = flag.String("f", "", "input file path")
	readLines     = flag.Int("l", 0, "read lines (optional)")
	sscs1FilePath = flag.String("sscs1", "", "sscs1 file path")
	sscs2FilePath = flag.String("sscs2", "", "sscs2 file path")
	dcs1FilePath  = flag.String("dcs1", "", "dcs1 file path")
	dcs2FilePath  = flag.String("dcs2", "", "dcs2 file path")
	anchorRegion  = flag.Int("anchor", 15, "anchor region")
	offsetRange   = flag.Int("offset", 4, "offset range")
	isReadPart    = false
	familyCount   = 0
)

func ParseArgs() {
	flag.Parse()
	if *filePath == "" {
		fmt.Println("请输入文件路径")
		os.Exit(1)
	}
	if *sscs1FilePath == "" || *sscs2FilePath == "" {
		fmt.Println("请输入sscs输出文件路径")
		os.Exit(1)
	}
	if *dcs1FilePath == "" || *dcs2FilePath == "" {
		fmt.Println("请输入dcs输出文件路径")
		os.Exit(1)
	}
	if *readLines > 0 {
		isReadPart = true
	}
	if anchorRegion != nil {
		core.AnchorRegion = *anchorRegion
	}
	if offsetRange != nil {
		core.OffsetRange = *offsetRange
	}
}

func checkLines(i int) bool {
	if !isReadPart {
		return true
	}
	return i < *readLines
}

func HandleResult(handle func([]Result)) chan<- []Result {
	result := make(chan []Result, 100)
	go func() {
		for {
			select {
			case r := <-result:
				handle(r)
			}
		}
	}()
	return result
}

func WaitChan[T any](ch chan<- T) {
	for {
		if len(ch) == 0 {
			return
		}
	}
}

func Process(migs map[Key]Mig, sscsResult, dcsResult chan<- []Result, wg *sync.WaitGroup) {
	defer wg.Done()
	familyCount++
	results := make([]Result, 0)
	for k, v := range migs {
		seq, qual := core.MakeSSCS(v.families, v.qualities)
		if len(seq) != 0 {
			result := Result{
				Sequence: seq,
				Quality:  qual,
				Key:      k,
			}
			result.NReads = strconv.Itoa(len(v.families))
			results = append(results, result)
		}
	}
	sscsResult <- results
	if len(results) == 4 {
		aligns := make(map[core.OrderMate]string)
		for _, result := range results {
			orderMate := core.OrderMate{
				Order: result.Key.Order,
				Mate:  result.Key.Meta,
			}
			aligns[orderMate] = result.Sequence
		}
		dcss := make([]Result, 0)
		seqs := core.MakeDCS(aligns)
		for i, seq := range seqs {
			r := Result{
				Sequence: seq,
				Quality:  strings.Repeat("I", len(seq)),
				Key: Key{
					Barcode: results[0].Key.Barcode,
					Meta:    strconv.Itoa(i + 1),
				},
			}
			r.NReads = results[2*i].NReads + "-" + results[2*i+1].NReads
			dcss = append(dcss, r)
		}
		dcsResult <- dcss
	}
}

func main() {
	ParseArgs()
	startTime := time.Now()
	csvReader := reader.NewCSVReader()
	err := csvReader.Read(*filePath, '\t')
	if err != nil {
		panic(err)
	}

	sscs1File, err := os.Create(*sscs1FilePath)
	if err != nil {
		panic(err)
	}
	sscs2File, err := os.Create(*sscs2FilePath)
	if err != nil {
		panic(err)
	}
	sscsResult := HandleResult(func(results []Result) {
		for _, result := range results {
			if len(result.Sequence) != 0 {
				outputLine := fmt.Sprintf("@%s.%s %s\n%s\n+\n%s\n",
					result.Key.Barcode, result.Key.Order, result.Key.Meta, result.Sequence, result.Quality)
				if result.Key.Meta == "1" {
					_, _ = sscs1File.WriteString(outputLine)
				} else {
					_, _ = sscs2File.WriteString(outputLine)
				}
			}
		}
	})
	dcs1File, err := os.Create(*dcs1FilePath)
	if err != nil {
		panic(err)
	}
	dcs2File, err := os.Create(*dcs2FilePath)
	if err != nil {
		panic(err)
	}
	var lock sync.Mutex
	dcsResult := HandleResult(func(results []Result) {
		if len(results) != 2 {
			return
		}
		result1 := results[0]
		result2 := results[1]
		if len(result1.Sequence) != 0 && len(result2.Sequence) != 0 {
			outputLine1 := fmt.Sprintf("@%s %s\n%s\n+\n%s\n",
				result1.Key.Barcode, result1.NReads, result1.Sequence, result1.Quality)
			outputLine2 := fmt.Sprintf("@%s %s\n%s\n+\n%s\n",
				result2.Key.Barcode, result2.NReads, result2.Sequence, result2.Quality)
			lock.Lock()
			_, _ = dcs1File.WriteString(outputLine1)
			_, _ = dcs2File.WriteString(outputLine2)
			lock.Unlock()
		}
	})
	var families []string
	var qualities []string
	migs := make(map[Key]Mig)
	var currentKey Key
	var wg sync.WaitGroup
	var i int
	for line, err := csvReader.NextLine(); err != io.EOF && checkLines(i); line, err = csvReader.NextLine() {
		barcode, order, meta, sequence, quality := line[0], line[1], line[2], line[4], line[5]
		if barcode == "" || sequence == "" || quality == "" {
			continue
		}
		key := Key{barcode, order, meta}
		if key != currentKey {
			if len(families) != 0 {
				mig := Mig{families, qualities}
				migs[currentKey] = mig
			}
			if key.Barcode != currentKey.Barcode && len(migs) != 0 {
				wg.Add(1)
				go Process(migs, sscsResult, dcsResult, &wg)
				migs = make(map[Key]Mig, 0)
			}
			families = make([]string, 0)
			qualities = make([]string, 0)
			currentKey = key
		}
		families = append(families, sequence)
		qualities = append(qualities, quality)
		i++
	}
	wg.Add(1)
	if len(families) != 0 {
		mig := Mig{families, qualities}
		migs[currentKey] = mig
	}
	go Process(migs, sscsResult, dcsResult, &wg)
	wg.Wait()
	_ = sscs1File.Close()
	_ = sscs2File.Close()
	_ = dcs1File.Close()
	_ = dcs2File.Close()

	fmt.Println("耗时:", time.Since(startTime))
	fmt.Println("总处理:", i, "行")
}
