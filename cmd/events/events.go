package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

var path string

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
	Queue       = []string{}
	Queuers     = map[string]struct{}{}
	Players     = map[string]*int{}
	Customers   = map[string]struct{}{}
	Computers   = map[int]*computer{}
)

func init() {
	flag.StringVar(&path, "path", "./testdata/test.txt", "путь к файлу с исходными данными по итогам дня")
}

func main() {
	flag.Parse()

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("could not open file with path %s. error: %v", path, err)
	}
	defer file.Close()

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

			parseInboundEvents(line, time, message[1], message[2], computerID)
			continue
		}

		parseInboundEvents(line, time, message[1], message[2])
	}

	// Все игроки, которые остались на момент закрытия должны освободить места и быть учтены в выручке
	for key, val := range Players {
		Computers[*val].SessionEnd, Computers[*val].Occupied = closingTime, false
		Computers[*val].Revenue += calculateRevenue(Computers[*val])
		Computers[*val].UseTime += closingTime.Sub(Computers[*val].SessionStart)
		Event11(closingTime, key)
	}

	// Выводим время закрытия
	fmt.Println(closingTime.Format("15:04"))

	// Вывод выручки и активного времени по каждому столу
	for i := 1; i <= capacity; i++ {
		fmt.Println(i, Computers[i].Revenue, fmtDuration(Computers[i].UseTime))
	}
}

func parseInboundEvents(line string, moment time.Time, eventID string, client string, computerID ...int) {
	switch eventID {
	// Клиент пришел
	case "1":
		// если клиент пришел вне часов работы
		if moment.Before(openingTime) || moment.After(closingTime) {
			Event13(moment, "NotOpenYet")
			return
		}

		// нельзя войти в одну реку дважды
		_, entered := Customers[client]
		if entered {
			Event13(moment, "YouShallNotPass")
			return
		}

		// отмечаем как вошедшего
		Customers[client] = struct{}{}

	// Клиент сел за стол
	case "2":
		if computerID == nil {
			log.Fatalf("no computer id provided, line: %v", line)
		}

		// если компьютер уже занят
		if Computers[computerID[0]].Client != "" {
			Event13(moment, "PlaceIsBusy")
			return
		}

		_, waiting := Queuers[client]
		station, playing := Players[client]

		// если клиент не входил
		_, entered := Customers[client]
		if !entered {
			Event13(moment, "ClientUnknown")
			return
		}

		// если клиент уже за каким-то столом
		if playing {
			Computers[*station].SessionEnd, Computers[computerID[0]].Occupied = moment, false
			Computers[*station].UseTime += moment.Sub(Computers[*station].SessionStart)
			Computers[*station].Revenue = calculateRevenue(Computers[*station])
			Computers[computerID[0]].Client = client
			Computers[computerID[0]].SessionStart, Computers[computerID[0]].Occupied = moment, true
			Players[client] = &computerID[0]
			return
		}

		// если клиент в очереди
		if waiting {
			delete(Queuers, client)
			Queue = Queue[1:]
			Event12(moment, client, computerID[0])
		}
		Computers[computerID[0]].Client = client
		Computers[computerID[0]].SessionStart, Computers[computerID[0]].Occupied = moment, true
		Players[client] = &computerID[0]

	// Клиент ожидает
	case "3":
		// если очередь превысила количество столов
		if len(Queuers) > len(Computers) {
			Event11(moment, client)
			return
		}

		// если клиент ожидает, несмотря на то, что есть свободные столы
		var vacant bool
		for _, val := range Computers {
			if !val.Occupied {
				vacant = true
			}
		}
		if vacant {
			Event13(moment, "ICanWaitNoLonger!")
			return
		}

		Queuers[client] = struct{}{}
		Queue = append(Queue, client)

	// Клиент ушел
	case "4":
		// если клиент ушел не заходя
		_, entered := Customers[client]
		if !entered {
			Event13(moment, "ClientUnknown")
			return
		}

		delete(Customers, client)

		if len(Queue) > 0 {
			// первый в очереди отмечается как работающий за компьютером ушедшего
			Players[Queue[0]] = Players[client]
			// сеанс ушедшего игрока завершается и подсчитывается его время и чек
			Computers[*Players[client]].SessionEnd, Computers[*Players[client]].Occupied = moment, false
			Computers[*Players[client]].UseTime += moment.Sub(Computers[*Players[client]].SessionStart)
			Computers[*Players[client]].Revenue = calculateRevenue(Computers[*Players[client]])
			Event12(moment, Queue[0], *Players[client])
			//начинается сеанс первого игрока из очереди
			Computers[*Players[client]].SessionStart, Computers[*Players[client]].Occupied = moment, true
			// ушедший удаляется из списка активных клиентов, покинувший очередь удаляется из очереди
			delete(Players, client)
			delete(Queuers, Queue[0])
			//очередь подрезается
			Queue = Queue[1:]
			return
		}
		// сеанс ушедшего игрока завершается и подсчитывается его время и чек, он удаляется из списка активных игроков
		Computers[*Players[client]].SessionEnd, Computers[*Players[client]].Occupied = moment, false
		Computers[*Players[client]].UseTime += moment.Sub(Computers[*Players[client]].SessionStart)
		Computers[*Players[client]].Revenue += calculateRevenue(Computers[*Players[client]])
		delete(Players, client)
	}
}

// calculateRevenue Функция вычисления выручки с округлением по часам
func calculateRevenue(pc *computer) int64 {
	dur := pc.SessionEnd.Sub(pc.SessionStart)
	return int64(math.Ceil(dur.Hours())) * price
}

// TODO: попробовать единую функцию для исходящих событий с вариадическими аргументами

// Event11 Клиент ушел
func Event11(moment time.Time, client string) {
	fmt.Printf("%v 11 %s\n", moment.Format("15:04"), client)
}

// Event12 Клиент сел за стол
func Event12(moment time.Time, client string, computerID int) {
	fmt.Printf("%v 12 %s %d\n", moment.Format("15:04"), client, computerID)
}

// Event13 Ошибка
func Event13(moment time.Time, err string) {
	fmt.Printf("%v 13 %s\n", moment.Format("15:04"), err)
}

// fmtDuration ЫФункция форматирования длительности сеанса
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}
