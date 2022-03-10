package src

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var symbols = []string{" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}

func main() {
	fmt.Printf("%c[?25l", 27)
	defer fmt.Printf("%c[?25h", 27)
	if len(os.Args) < 4 {
		fmt.Println("Usage:")
		fmt.Println("\t" + os.Args[0] + " <urlPattern> <resultFileName> <start> <end> <format>")
		fmt.Println("")
		fmt.Println("In pattern you should use '{X}' or {XXXX} placeholder. Which will be replaced with part number.")
		fmt.Println()
		fmt.Println()
		os.Exit(0)
	}

	link := os.Args[1]
	fileName := os.Args[2]
	start, err := strconv.Atoi(os.Args[3])
	end, err := strconv.Atoi(os.Args[4])
	format := len(os.Args) == 6 && os.Args[5] != ""
	if err != nil {
		fmt.Println("Error: partsCount is not an int value")
		os.Exit(1)
	}

	err = downloadMovie(link, fileName, start, end, format)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("Download complete")
}

func downloadMovie(from string, to string, start int, end int, format bool) error {
	file, err := os.OpenFile(to, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_EXCL, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := start; i < end; i++ {
		downloadPart(from, file, i, format)
		_renderProgress(end-start-1, i)
	}
	return nil
}

func downloadPart(from string, to *os.File, part int, format bool) error {
	var partNumber string
	if format {
		partNumber = fmt.Sprintf("%04d", part)
	} else {
		partNumber = strconv.Itoa(part)
	}
	downloadUrl := strings.Replace(from, "{number}", partNumber, 1)
	//fmt.Println("Downloading part", part, " url: ", downloadUrl)
	partContent, err := http.Get(downloadUrl)
	if err != nil {
		return err
	}
	defer partContent.Body.Close()

	if partContent.ContentLength < 1000 {
		fmt.Println("Warning: small content ", partContent.ContentLength)
	}
	//fmt.Println("Part", part, " downloaded, size: ", partContent.ContentLength)
	content, err := ioutil.ReadAll(partContent.Body)
	_, err = to.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func _renderProgress(parts int, part int) int {
	fmt.Printf("%d/%d\n", part, parts)
	bars := _renderBars(parts, part)

	if part != parts {
		fmt.Printf("%c", 13)
		fmt.Printf("%c[%dA", 27, bars+1)
	}
	return bars
}

func _renderBars(parts int, part int) int {
	width := parts / len(symbols)
	if width > 50 {
		width = 50
	}
	_bar(width, parts, part)
	pp := parts / width
	if pp > len(symbols)*2 {
		return _renderBars(pp, part%pp) + 1
	}
	return 1
}

func _bar(width int, parts int, position int) {
	pp := parts / width
	f := position / pp
	pt := position % pp
	fmt.Print("|")
	fmt.Print(strings.Repeat(symbols[len(symbols)-1], f))
	if f != width {
		s := (pt * (len(symbols) - 1)) / pp
		fmt.Print(symbols[s])
		fmt.Print(strings.Repeat(" ", width-1-f))
	}
	fmt.Println("|")
}
