package main
import (
  "fmt"
  "bufio"
  "os"
  "log"
  "math/rand"
  "time"
  "strings"
  "errors"
  "unicode"
  "unicode/utf8"
)

func process_line(line string) (string,error) {
  line = strings.ToLower(line)
  for _,ch := range line{
    if !unicode.IsLetter(ch){
        return "", errors.New("Not a valid letter")
    }
  }
  return line, nil
}

func main(){
  rand.Seed(time.Now().UnixNano())
  read_file, err := os.Open("wiki-100k.txt") 
  if err != nil {
    log.Fatal(err)
  }
  defer read_file.Close()
  scanner := bufio.NewScanner(read_file)
  
  fmt.Println("Enter a 5 letter word")
  var input string
  fmt.Scan(&input)
  
  line_count := 0

  write_file,_ := os.Create("new.txt")
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
          line_count += 1
        }
    }
  
  fmt.Println("LINE COUNT: ", line_count)
  random_index := rand.Intn(line_count)
  fmt.Println("RANDOM INDEX:", random_index)

  if len(input) != 5 {
    log.Fatal("ERROR: The word is not 5 character long!!")  
  }
  fmt.Println(input)

}
