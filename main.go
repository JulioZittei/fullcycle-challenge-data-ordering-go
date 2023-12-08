package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Person struct {
	Name string
	Age int
	Score int
}

type People []*Person

func (p People) Len() int {
	return len(p)
}

func (p People) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p People) Less(i, j int) bool { 
	return p[i].Name < p[j].Name 
}

func (p People) LessByAge(i, j int) bool {
	return p[i].Age < p[j].Age
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Use: go run main.go <arquivo-origem.csv> <arquivo-destino.csv>")
		os.Exit(1)
	}

	input := os.Args[1]
	output := os.Args[2]

	people, err := readFile(input)
	if err != nil {
		fmt.Printf("Error while reading file: %v\n", err)
		os.Exit(1)
	}

	sort.Sort(People(people))
	saveFile(people, "Ordered by name", output)

	sort.Slice(people, people.LessByAge)
	saveFile(people, "Ordered by age", output)

}

func readFile(filePath string) (People, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var people People

	if len(lines) > 0 {
		begin := 0
		_, err := strconv.Atoi(lines[0][1]) 
		if err != nil {
			begin = 1
		}

		for _, line := range lines[begin:] {
			name:= line[0]
			age, err := strconv.Atoi(line[1])
			if err != nil {
				fmt.Printf("Age not found: %v\n", err)
			}
			score, err := strconv.Atoi(line[2])
			if err != nil {
				fmt.Printf("Score not found: %v\n", err)
			}
	
			if err == nil {
				person := &Person{
					Name: name,
					Age: age,
					Score: score,
				}
				people = append(people, person)
			}			
		}
	}
	return people, nil
}

func saveFile(people People, msg string, outputFilePath string) {
	file, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Printf("Error while creating file %s: %v\n", msg, err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escreve o cabeçalho.
	writer.Write([]string{"Nome", "Idade", "Pontuação"})

	// Escreve os dados ordenados.
	for _, person := range people {
		writer.Write([]string{person.Name, strconv.Itoa(person.Age), strconv.Itoa(person.Score)})
	}

	fmt.Printf("Data %s was saved at %s with success!\n", msg, outputFilePath)
}