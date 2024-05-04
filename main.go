package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

func process_line(line string) (string, error) {
	line = strings.ToLower(line)
	for _, ch := range line {
		if !unicode.IsLetter(ch) {
			return "", errors.New("Not a valid letter")
		}
	}
	return line, nil
}

func process_file() {
	read_file, err := os.Open("wiki-100k.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer read_file.Close()
	scanner := bufio.NewScanner(read_file)
	write_file, _ := os.Create("new.txt")
	defer write_file.Close()
	writer := bufio.NewWriter(write_file)
	for scanner.Scan() {
		line := scanner.Text()
		if utf8.RuneCountInString(line) == 5 && len(line) == 5 {
			line, err := process_line(scanner.Text())
			if err != nil {
				continue
			}
			writer.WriteString(line + "\n")
			writer.Flush()
		}
	}
}

func make_words_from_scanner(s *bufio.Scanner) []string {
	var lines []string
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines
}

func get_random_word(lines []string) (string, []string) {
	line_count := len(lines)
	line := rand.Intn(line_count)
	var word string
	word = lines[line]
	// TODO: this approach is very slow can optimize this
	// currently just removes the randomly selected word from the list
	lines = append(lines[:line], lines[line+1:]...)
	return word, lines
}

// IDEA:
// The 5 lettered word should consist of 5 cells
// |X| |X| |X| |X| X|
// where each cell stores a char and an assosiated state
// STATE := CORRECT | EXITS | ABSENT

type CellState int

const (
	ABSENT CellState = iota
	PRESENT
	CORRECT
)

type Cell struct {
	ch    byte
	state CellState
}

type Wordle struct {
	cells []Cell
}

// filter the total dictionary by passing it through
// 3 filter functions
// 1. filter_absent()
// 2. filter_present()
// 3. filter_correct()
// and give out a new list of all possible words from initial dict.

func generate_has_correct(cells []Cell) [5]bool {
	has_correct := [5]bool{}
	for i, cell := range cells {
		if cell.state == CORRECT {
			has_correct[i] = true
		}
	}
	return has_correct
}

func can_remove(str string, b byte, cells []Cell) bool {
	char := rune(b)
	has_correct := generate_has_correct(cells)

	for i, c := range str {
		if c == char && !has_correct[i] {
			return true
		}
	}
	return false
}

func filter_absent(cells []Cell, words []string) []string {
	new_words := []string{}
	has_correct := generate_has_correct(cells)
	for _, word := range words {
		flag := false
		for i, cell := range cells {
			if has_correct[i] {
				continue
			}
			if cell.state == ABSENT && can_remove(word, cell.ch, cells) {
				flag = true
				break
			}
		}
		if !flag {
			new_words = append(new_words, word)
		}
	}
	fmt.Println("AFTER FILTERING ABSENT WORDS")
	fmt.Println(new_words[:10])
	return new_words
}

func filter_correct(cells []Cell, words []string) []string {
	new_words := []string{}
	for _, word := range words {
		flag := false
		for i, cell := range cells {
			if cell.state == CORRECT && word[i] == cell.ch {
				flag = true
				break
			}
		}
		if flag {
			new_words = append(new_words, word)
		}
	}
	fmt.Println("AFTER FILTERING CORRECT WORDS")
	fmt.Println(new_words)
	return new_words
}

func make_educated_guess() {

}

func parse_cell(input string) Cell {
	if len(input) > 2 && len(input) < 1 {
		log.Fatal("Cell Input has illegal amount of characters")
	}

	input_char := input[0]

	if input_char >= 65 && input_char <= 91 {
		input_char += 32
	}
	cell := Cell{
		input_char,
		ABSENT,
	}
	if len(input) == 1 {
		return cell
	}
	switch input[1] {
	case 'P':
		cell.state = PRESENT
	case 'C':
		cell.state = CORRECT
	default:
	}
	return cell
}

func take_cell_input() []Cell {
	cells := []Cell{}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	input := strings.Split(scanner.Text(), " ")
	for _, in := range input {
		cells = append(cells, parse_cell(in))
	}
	fmt.Println(cells)
	return cells
}

func main() {
	read_file, err := os.Open("new.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer read_file.Close()
	scanner := bufio.NewScanner(read_file)

	words := make_words_from_scanner(scanner)
	word_count := len(words)
	fmt.Println("LINE COUNT: ", word_count)

	// todays_word, _ := get_random_word(words)
	// fmt.Println("RANDOM WORD:", todays_word)

	//	is_solved := false
	//	for !is_solved {
	//		random_word, new_words := get_random_word(words)
	//		words = new_words
	//		if todays_word == random_word {
	//			fmt.Println("THE CORRECT WORD IS:", random_word)
	//			is_solved = true
	//		}
	//	}

	//no_of_tries := len(words)
	//fmt.Println("TRIES:", word_count-no_of_tries)

	current_cells := take_cell_input()
	correct_word := filter_correct(current_cells, words)
	filter_absent(current_cells, correct_word)

}
