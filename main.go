package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorCyan   = "\033[36m"
	ColorPink   = "\033[35m"
)

type GameResult struct {
	Date     string `json:"date"`
	Outcome  string `json:"outcome"`
	Attempts int    `json:"attempts"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(ColorPink + "Игра 'Угадай число'!" + ColorReset)

	for {
		playGame(scanner)

		fmt.Print(ColorPink + "\nХотите сыграть еще раз? (y/n): " + ColorReset)
		if !readYesNo(scanner) {
			fmt.Println("Спасибо за игру!")
			break
		}
	}
}

func playGame(scanner *bufio.Scanner) {
	maxNum, maxAttempts := chooseDifficulty(scanner)
	secretNumber := generateNumber(maxNum)

	fmt.Printf("\nИгра 'Угадай число' - от 1 до %d началась!\n", maxNum)
	fmt.Printf("Угадайте число за %d попыток!\n", maxAttempts)

	var history []int
	won := false
	attemptsUsed := 0

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		attemptsUsed = attempt

		if len(history) > 0 {
			fmt.Printf("Ваши прошлые попытки: %v\n", history)
		}

		fmt.Printf(ColorYellow+"Попытка #%d - Введите число: "+ColorReset, attempt)

		guess := readInt(scanner, 1, maxNum)
		history = append(history, guess)

		isWin, directionHint, distanceHint := checkGuess(guess, secretNumber)

		if isWin {
			fmt.Println(ColorGreen + "Вы угадали!")
			fmt.Println(ColorGreen + "Игра закончена победой!")
			won = true
			break
		} else {
			fmt.Printf("%s. %s!\n", directionHint, distanceHint)
		}
	}

	if !won {
		fmt.Println(ColorRed + "Вы проиграли! Попытки закончились.")
		fmt.Printf(ColorRed+"Секретное число было: %d\n"+ColorReset, secretNumber)
	}

	saveResult(won, attemptsUsed)
}

func chooseDifficulty(scanner *bufio.Scanner) (maxNum int, maxAttempts int) {
	fmt.Println("\nВыберите сложность:")
	fmt.Println("1 - Easy: 1–50, 15 попыток")
	fmt.Println("2 - Medium: 1–100, 10 попыток")
	fmt.Println("3 - Hard: 1–200, 5 попыток")
	fmt.Print("Ваш выбор (1/2/3): ")

	choice := readInt(scanner, 1, 3)

	switch choice {
	case 1:
		return 50, 15
	case 2:
		return 100, 10
	case 3:
		return 200, 5
	}
	return 100, 10
}

func generateNumber(max int) int {
	return rand.Intn(max) + 1
}

func readInt(scanner *bufio.Scanner, min, max int) int {
	for {
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		number, err := strconv.Atoi(input)
		if err != nil {
			fmt.Print(ColorRed + "Ошибка: введите целое число! Попробуйте снова: " + ColorReset)
			continue
		}

		if number < min || number > max {
			fmt.Printf(ColorRed+"Ошибка: число должно быть от %d до %d! Попробуйте снова: "+ColorReset, min, max)
			continue
		}

		return number
	}
}

func readYesNo(scanner *bufio.Scanner) bool {
	for {
		scanner.Scan()
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if answer == "y" || answer == "да" || answer == "yes" {
			return true
		}
		if answer == "n" || answer == "нет" || answer == "no" {
			return false
		}

		fmt.Print(ColorRed + "Пожалуйста, введите 'y' (да) или 'n' (нет): " + ColorReset)
	}
}

func checkGuess(guess, secret int) (isWin bool, direction string, distance string) {
	if guess == secret {
		return true, "", ""
	}

	if guess > secret {
		direction = "Секретное число меньше"
	} else {
		direction = "Секретное число больше"
	}

	diff := guess - secret
	if diff < 0 {
		diff = -diff
	}

	if diff <= 5 {
		distance = ColorRed + "Горячо" + ColorReset
	} else if diff <= 15 {
		distance = ColorYellow + "Тепло" + ColorReset
	} else {
		distance = ColorCyan + "Холодно" + ColorReset
	}

	return false, direction, distance
}

func saveResult(won bool, attempts int) {
	filename := "results.json"

	outcome := "Проигрыш"
	if won {
		outcome = "Победа"
	}

	newResult := GameResult{
		Date:     time.Now().Format("2006-01-02 15:04:05"),
		Outcome:  outcome,
		Attempts: attempts,
	}

	var results []GameResult
	fileData, err := os.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(fileData, &results)
	}

	results = append(results, newResult)
	updatedData, _ := json.MarshalIndent(results, "", "  ")
	_ = os.WriteFile(filename, updatedData, 0644)
}
