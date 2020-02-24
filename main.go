package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Libraries []Library

type Book struct {
	Id    int
	Value int
	Flag  bool
}

func (b Book) String() string {
	return fmt.Sprintf("%d %d %b", b.Id, b.Value, b.Flag)
}

type Library struct {
	Id             int
	Books          []*Book
	SignupDuration int
	Speed          int
	ScoreMax       int
}

func (l Library) GetMaxScore(nbDay int, useFlag bool) int {
	result := 0
	bookDayIndex := 0
	for bookIndex := 0; bookDayIndex < nbDay*l.Speed && bookIndex < len(l.Books); bookIndex++ {
		if !useFlag || !l.Books[bookIndex].Flag {
			result += l.Books[bookIndex].Value
			bookDayIndex++
		}
	}
	return result
}

type Score int

func (s Score) String() string {
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", s)
}

func ParseInput(fileName string) (*Libraries, int, error) {
	file, err := os.Open("data/" + fileName)
	if err != nil {
		return nil, 0, fmt.Errorf("open file error for %s : %v", fileName, err)
	}
	defer file.Close()
	scanner := bufio.NewReader(file)
	firstLine, _ := scanner.ReadString('\n')
	firstLine = firstLine[0 : len(firstLine)-1]
	firstLineData := strings.Split(firstLine, " ")

	nbBooks, _ := strconv.Atoi(firstLineData[0])
	nbLibrary, _ := strconv.Atoi(firstLineData[1])
	nbDays, _ := strconv.Atoi(firstLineData[2])

	allBook := []*Book{}

	secondLine, _ := scanner.ReadString('\n')
	secondLine = secondLine[0 : len(secondLine)-1]
	secondLineData := strings.Split(secondLine, " ")
	for i := 0; i < nbBooks; i++ {
		v := secondLineData[i]
		var book Book
		book.Id = i
		book.Flag = false
		book.Value, err = strconv.Atoi(v)
		if err != nil {
			return nil, 0, fmt.Errorf("read integer error %s : %v", v, err)
		}
		allBook = append(allBook, &book)
	}

	libraries := Libraries{}

	for i := 0; i < nbLibrary; i++ {
		libraryLine1, _ := scanner.ReadString('\n')
		libraryLine1 = libraryLine1[0 : len(libraryLine1)-1]
		libraryLine1Data := strings.Split(libraryLine1, " ")

		var library Library
		library.Id = i

		library.SignupDuration, _ = strconv.Atoi(libraryLine1Data[1])
		library.Speed, _ = strconv.Atoi(libraryLine1Data[2])

		libraryLine2, _ := scanner.ReadString('\n')
		libraryLine2 = libraryLine2[0 : len(libraryLine2)-1]
		for _, v := range strings.Split(libraryLine2, " ") {
			bookId, _ := strconv.Atoi(v)
			library.Books = append(library.Books, allBook[bookId])
		}

		sort.Slice(library.Books, func(i, j int) bool {
			return library.Books[i].Value > library.Books[j].Value
		})

		libraries = append(libraries, library)
	}

	return &libraries, nbDays, nil
}

func (libraries Libraries) Dump(fileName string) (int, error) {

	t := time.Now()
	output := fmt.Sprintf("result/%s-%s.result", fileName, t.Format("20060102T150405"))
	f, err := os.Create(output)
	if err != nil {
		return 0, fmt.Errorf("error create file %s : %v", output, err)
	}
	defer f.Close()

	score := 0
	var alreadyShip = map[int]int{}

	w := bufio.NewWriter(f)

	num := len(libraries)
	w.WriteString(fmt.Sprintf("%d", num))
	w.WriteString("\n")
	for _, l := range libraries {
		w.WriteString(fmt.Sprintf("%d %d", l.Id, len(l.Books)))
		w.WriteString("\n")
		var stringBookId []string
		for _, b := range l.Books {
			stringBookId = append(stringBookId, fmt.Sprintf("%d", b.Id))

			_, inside := alreadyShip[b.Id]
			if !inside {
				alreadyShip[b.Id] = 1
				score += b.Value
			}
		}
		w.WriteString(strings.Join(stringBookId, " "))
		w.WriteString("\n")
	}
	w.Flush()

	p := message.NewPrinter(language.English)
	p.Printf("Result : %d\n", score)
	return score, nil
}

func (libraries Libraries) Score() (int, error) {
	score := 0
	var alreadyShip = map[int]int{}

	for _, l := range libraries {
		var stringBookId []string
		for _, b := range l.Books {
			stringBookId = append(stringBookId, fmt.Sprintf("%d", b.Id))

			_, inside := alreadyShip[b.Id]
			if !inside {
				alreadyShip[b.Id] = 1
				score += b.Value
			}
		}
	}

	fmt.Printf("Result : %s\n", Score(score))
	return score, nil
}

func main() {
	total := 0
	total += compute("a_example.txt", 1)
	total += compute("b_read_on.txt", 10)
	total += compute("c_incunabula.txt", 100)
	total += compute("d_tough_choices.txt", 130)
	total += compute("e_so_many_books.txt", 1000)
	total += compute("f_libraries_of_the_world.txt", 280000)

	fmt.Printf("Total : %s\n", Score(total))
}

func compute(fileName string, minScore int) int {

	writeFile := false

	libraries, nbDay, err := ParseInput(fileName)
	if err != nil {
		fmt.Printf("%v", err)
		panic("open file error")
	}
	var result Libraries

	dayIndex := 0

	for len(*libraries) > 0 {

		for i := 0; i < len(*libraries); i++ {
			val := (*libraries)[i].GetMaxScore(nbDay-(*libraries)[i].SignupDuration-dayIndex, true)
			(*libraries)[i].ScoreMax = val
		}

		sort.Slice(*libraries, func(i, j int) bool {
			vi := (*libraries)[i].ScoreMax / (*libraries)[i].SignupDuration
			vj := (*libraries)[j].ScoreMax / (*libraries)[j].SignupDuration
			return vi > vj
		})

		lib := (*libraries)[0]

		if (dayIndex + lib.SignupDuration) < nbDay {
			var newLib Library
			newLib.Id = lib.Id
			newLib.Speed = lib.Speed
			dayIndex += lib.SignupDuration

			bookDayIndex := 0
			for bookIndex := 0; bookDayIndex < ((nbDay-dayIndex)*lib.Speed) && bookIndex < len(lib.Books); bookIndex++ {
				if !lib.Books[bookIndex].Flag {
					newLib.Books = append(newLib.Books, lib.Books[bookIndex])
					lib.Books[bookIndex].Flag = true
					bookDayIndex++
				}
			}

			maxScore := newLib.GetMaxScore(nbDay-dayIndex, false)
			if len(newLib.Books) == 0 || maxScore < minScore {
				dayIndex -= lib.SignupDuration
			} else {
				result = append(result, newLib)
			}
		}
		*libraries = (*libraries)[1:len(*libraries)]
	}

	score := 0
	if writeFile {
		score, _ = result.Dump(fileName)
	} else {
		score, _ = result.Score()
	}
	return score
}
