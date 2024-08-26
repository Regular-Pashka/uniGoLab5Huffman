package main

import (
	"bufio"
	"fmt"
	// "io"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Node struct {
	left    *Node
	right   *Node
	value   int
	isLeaf  bool
	symbol  byte
}

type ByValue []struct {
	Key   byte
	Value int
}

func (bv ByValue) Len() int           { return len(bv) }
func (bv ByValue) Less(i, j int) bool { return bv[i].Value < bv[j].Value }
func (bv ByValue) Swap(i, j int)      { bv[i], bv[j] = bv[j], bv[i] }

func sortByValue(alphabet map[byte]int) []struct {
	Key   byte
	Value int
} {
	sorted := make([]struct {
		Key   byte
		Value int
	}, len(alphabet))

	i := 0
	for k, v := range alphabet {
		sorted[i] = struct {
			Key   byte
			Value int
		}{k, v}
		i++
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Value < sorted[j].Value
	})

	return sorted
}

func createQueue(alphabet []struct {
	Key   byte
	Value int
}) []*Node {
	queue := []*Node{}
	for _, v := range alphabet {
		queue = append(queue, &Node{symbol: v.Key, value: v.Value, isLeaf: true})
	}
	return queue
}

func generateFixedLengthCodes(alphabet []struct {
	Key   byte
	Value int
}) map[byte]string {
	fmt.Println(alphabet, "alphabet")
	codes := make(map[byte]string)
	codeLength := int(math.Ceil(math.Log(float64(len(alphabet))) / math.Log(2)))
	codeNumber := 0
	for _, letter := range alphabet {
		code := strconv.FormatInt(int64(codeNumber), 2)
		/* 
			возможно эти коды генерятся неверно
		*/
		for len(code) < codeLength {
			code = "0" + code
		}
		codes[byte(letter.Key)] = code
		codeNumber++
	}
	fmt.Println("Fixed-length codes generated")
	return codes
}

func generateHuffmanCodes(queue []*Node) *Node {
	for len(queue) > 1 {
		a := queue[0]
		queue = queue[1:]
		b := queue[0]
		queue = queue[1:]
		newNode := &Node{left: a, right: b, value: a.value + b.value}
		queue = append(queue, newNode)
		sort.Slice(queue, func(i, j int) bool { return queue[i].value < queue[j].value })
	}
	return queue[0]
}

func printHuffmanCodes(node *Node, code string) {
	if node == nil {
		return
	}
	if node.isLeaf {
		fmt.Printf("%c: %s\n", node.symbol, code)
	}
	printHuffmanCodes(node.left, code+"0")
	printHuffmanCodes(node.right, code+"1")
}

func generateMap(node *Node, codes map[byte]string, code string) {
	if node == nil {
		return
	}
	if node.isLeaf {
		codes[node.symbol] = code
	}
	generateMap(node.left, codes, code+"0")
	generateMap(node.right, codes, code+"1")
}

func generateAlphabet(fileLines []string) /* map[byte]int */ []struct {
	Key   byte
	Value int
} {
	/* эту функцию нужно переделать  */
	alphabet := make(map[byte]int)
	for _, s := range fileLines {
		for _, char := range []byte(s) {
			alphabet[char]++
		}
	}
	sortedAlphabet := sortByValue(alphabet)
	
	return sortedAlphabet
}

func compress(fileLines []string, codes map[byte]string, filename string) {
	fmt.Println(fileLines, codes)
	/* 
		мне не нравится что кодирование идет вот так: [POKO] map[75:10 79:0 80:11]
		непонятно где что.
		Скорее всего процесс ломается при кодировании, а именно в функции генераторе кода. ломается он возможно из-за изначально кривой генерации алфавита.
	*/
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, str := range fileLines {
		for i := 0; i < len(str); i++ {
			_, err := file.WriteString(codes[str[i]])
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	fileLines := []string{}
	alphabet := make([]struct {
		Key   byte
		Value int
	}, 5) // надо менять эту переменную, ибо сортировать мапу невозможно
	fixedLengthCodes := make(map[byte]string)
	huffmanCodes := make(map[byte]string)

	for {
		fmt.Println("Меню")
		fmt.Println("1. Открыть текстовый файл")
		fmt.Println("2. Вывести содержимое текстового файла")
		fmt.Println("3. Сгенерировать алфавит с указанием частоты появления символов в файле")
		fmt.Println("4. Сгенерировать коды фиксированной длины")
		fmt.Println("5. Сгенерировать коды Хаффмана")
		fmt.Println("6. Сжать файл с кодами фиксированной длины")
		fmt.Println("7. Сжать файл с кодами Хаффмана")
		fmt.Println("8. Выйти из программы")
		fmt.Print("> ")

		var ans int
		fmt.Scan(&ans)

		switch ans {
		case 1:
			fmt.Print("Enter the path to the text file: ")
			filename, _ := in.ReadString('\n')
			filename = strings.TrimSpace(filename)

			file, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			reader := bufio.NewScanner(file)
			for reader.Scan() {
				line := reader.Text()
				fileLines = append(fileLines, line)
			}

			fmt.Println("File opened and ready to work")
		case 2:
			for _, s := range fileLines {
				fmt.Println(s)
			}
		case 3:
			alphabet = generateAlphabet(fileLines)
			fmt.Print(alphabet)
			for _, entry := range alphabet {
				fmt.Printf("%c: %d\n", entry.Key, entry.Value)
			}
		case 4:
			for _, v := range alphabet {
				fixedLengthCodes[v.Key] = ""
			}
			fmt.Println(alphabet, "pered generaciei")
			fixedLengthCodes = generateFixedLengthCodes(alphabet)
			// fmt.Print(fixedLengthCodes)
			for key, value := range fixedLengthCodes {
				fmt.Println(key, value)
			}
		case 5:
			queue := createQueue(alphabet)
			root := generateHuffmanCodes(queue)
			generateMap(root, huffmanCodes, "")
			fmt.Println("Huffman codes generated")
		case 6:
			fmt.Print("Enter the filename to compress with fixed-length codes: ")
			filename, _ := in.ReadString('\n')
			filename = strings.TrimSpace(filename)
			compress(fileLines, fixedLengthCodes, filename)
		case 7:
			fmt.Print("Enter the filename to compress with Huffman codes: ")
			filename, _ := in.ReadString('\n')
			filename = strings.TrimSpace(filename)
			compress(fileLines, huffmanCodes, filename)
		case 8:
			fmt.Println("Exiting the program")
			return
		default:
			fmt.Println("Invalid command!")
		}
	}
}

/* TODO:

Лабораторная работа №5. Коды Хаффмана.
Написать программу, реализующую кодирование символов алфавита
входного текстового файла в виде двоичных кодов:
- фиксированной длины;
- переменной длины.
Для генерации кодов переменной длины использовать жадный
алгоритм Хаффмана (код Хаффмана).
Реализовать меню с пунктами:
1 Открыть текстовый файл;
2 Вывести содержимое текстового файла;
3 Вывести символы алфавита с указанием их частоты появлен ия с
сортировкой по частоте;
4 Сгенерировать коды для символов алфавита входного файла
	
		TODO:
			4.1
			Вывести алфавит входного файла с кодами фиксированной
			длины для каждого символа алфавита;
			4.2
			Вывести алфавит входного файла с кодами Хаффмана для
			каждого символа алфавита;
		
5 Сжать содержимое текстового файла с помощью кодов фиксированной
длины с сохранением данных в файл;
6 Сжать содержимое текстового файла с помощью кодов Хаффмана с
сохранением данных в файл.
7 Сравнить размеры файлов исходного текстового файла и двух
зашифрованных.
Результаты лабораторной работы оформить в виде отчета с результатами
работы программы. 

*/


/* 
	ВЫВОДЫ:
	Алгосы сжатия не отрабатывают так, как должны на мой взглядж
	Надои изучить тему и порешить что как.

	25.08
	Не генерится нормально алфавит

	Требования к алфавиту:
*/