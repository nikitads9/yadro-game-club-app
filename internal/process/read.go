package process

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nikitads9/yadro-game-club-app/internal/format"
)

// computer структура игрового стола
type computer struct {
	// Имя клиента
	Client string
	// Занят ли стол
	Occupied bool
	// Начало сеанса
	SessionStart time.Time
	// Конец Сеанса
	SessionEnd time.Time
	// Активное время работы стола
	UseTime time.Duration
	// Суммарная выручка
	Revenue int64
}

var (
	price       int64
	openingTime time.Time
	closingTime time.Time
	// Queue Очередь клиентов
	Queue = []string{}
	// Queuers Отображение для быстрого (константного) доступа к клиентам в очереди
	Queuers = map[string]struct{}{}
	// Players Отображение для быстрого (константного) доступа к клиентам за компьютерами
	Players = map[string]*int{}
	// Customers Отображение для быстрого (константного) доступа к клиентам, вошедшим в клуб
	Customers = map[string]struct{}{}
	// Computers Отображение для быстрого (константного) доступа к столам с компьютерами
	Computers = map[int]*computer{}
)

// ReadLogs читает события из файла с помощью bufio и вызывает функцию по их обработке. В конце своей работы она выводит статистику по столам.
func ReadLogs(file *os.File) {
	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	capacity, err := strconv.Atoi(line)
	if err != nil {
		log.Fatalf("could not convert to integer: %v, line: %v", err, line)
	}

	//проверки на отрицательные и нулевые значения
	if capacity <= 0 {
		log.Fatalf("pc club capacity lower than or equal to zero, line: %v", line)
	}

	scanner.Scan()
	line = scanner.Text()
	workingHours := strings.Split(line, " ")

	openingTime, err = time.Parse("15:04", workingHours[0])
	if err != nil {
		log.Fatalf("could not parse time: %v, line: %v", err, line)
	}

	closingTime, err = time.Parse("15:04", workingHours[1])
	if err != nil {
		log.Fatalf("could not parse time: %v, line: %v", err, line)
	}
	if closingTime.Before(openingTime) {
		log.Fatalf("closing time is beforehand opening time. line: %v", line)
	}

	scanner.Scan()
	line = scanner.Text()
	price, err = strconv.ParseInt(line, 10, 64)
	if err != nil {
		log.Fatalf("could not convert to integer: %v, line: %v", err, line)
	}

	// проверка цены на отрицательные значения
	if price <= 0 {
		log.Fatalf("negative service price, line: %v", line)
	}

	// инициализация и нумерация столов
	for i := 1; i <= capacity; i++ {
		Computers[i] = &computer{}
	}

	// начало работы клуба
	fmt.Println(openingTime.Format("15:04"))

	for scanner.Scan() {
		line = scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		fmt.Println(line)

		// предполагается верный формат исходных данных: время, идентификатор события, тело и опционально идентификатор стола
		message := strings.Split(line, " ")
		time, err := time.Parse("15:04", message[0])
		if err != nil {
			log.Fatalf("could not parse time: %v, line: %v", err, line)
		}

		if message[2] == "" {
			log.Fatalf("no client name provided, line: %v", line)
		}

		// для события 2 необходимо указать, за какой стол сядет клиент
		if message[1] == "2" {
			if len(message) != 4 {
				log.Fatalf("no computer id provided, line: %v", line)
			}

			computerID, err := strconv.Atoi(message[3])
			if err != nil {
				log.Fatalf("could not convert computerID to integer: %v, line: %v", err, line)
			}

			processInboundEvents(line, time, message[1], message[2], computerID)
			continue
		}

		processInboundEvents(line, time, message[1], message[2])
	}

	// Все игроки, которые остались на момент закрытия должны освободить места и быть учтены в выручке
	for key, val := range Players {
		Computers[*val].SessionEnd, Computers[*val].Occupied = closingTime, false
		Computers[*val].Revenue += calculateRevenue(Computers[*val])
		Computers[*val].UseTime += closingTime.Sub(Computers[*val].SessionStart)
		format.Event(closingTime, 11, key)
	}

	// Выводим время закрытия
	fmt.Println(closingTime.Format("15:04"))

	// Вывод выручки и активного времени по каждому столу
	for i := 1; i <= capacity; i++ {
		fmt.Println(i, Computers[i].Revenue, format.FmtDuration(Computers[i].UseTime))
	}
}
