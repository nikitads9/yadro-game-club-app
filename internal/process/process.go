package process

import (
	"log"
	"time"

	"github.com/nikitads9/yadro-game-club-app/internal/format"
)

func processInboundEvents(line string, moment time.Time, eventID string, client string, computerID ...int) {
	switch eventID {
	// Клиент пришел
	case "1":
		// если клиент пришел вне часов работы
		if moment.Before(openingTime) || moment.After(closingTime) {
			format.Event(moment, 13, "NotOpenYet")
			return
		}

		// нельзя войти в одну реку дважды
		_, entered := Customers[client]
		if entered {
			format.Event(moment, 13, "YouShallNotPass")
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
			format.Event(moment, 13, "PlaceIsBusy")
			return
		}

		_, waiting := Queuers[client]
		station, playing := Players[client]

		// если клиент не входил
		_, entered := Customers[client]
		if !entered {
			format.Event(moment, 13, "ClientUnknown")
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
			format.Event(moment, 12, client, computerID[0])
		}
		Computers[computerID[0]].Client = client
		Computers[computerID[0]].SessionStart, Computers[computerID[0]].Occupied = moment, true
		Players[client] = &computerID[0]

	// Клиент ожидает
	case "3":
		// если очередь превысила количество столов
		if len(Queuers) > len(Computers) {
			format.Event(moment, 11, client)
			return
		}

		// если клиент ожидает, несмотря на то, что есть свободные столы
		var vacant bool
		for _, computer := range Computers {
			if !computer.Occupied {
				vacant = true
			}
		}
		if vacant {
			format.Event(moment, 13, "ICanWaitNoLonger!")
			return
		}

		Queuers[client] = struct{}{}
		Queue = append(Queue, client)

	// Клиент ушел
	case "4":
		// если клиент ушел не заходя
		_, entered := Customers[client]
		if !entered {
			format.Event(moment, 13, "ClientUnknown")
			return
		}

		player, playing := Players[client]

		delete(Customers, client)

		if playing && len(Queue) > 0 {
			// первый в очереди отмечается как работающий за компьютером ушедшего
			Players[Queue[0]] = player
			// сеанс ушедшего игрока завершается и подсчитывается его время и чек
			Computers[*player].SessionEnd, Computers[*player].Occupied = moment, false
			Computers[*player].UseTime += moment.Sub(Computers[*player].SessionStart)
			Computers[*player].Revenue = calculateRevenue(Computers[*player])
			format.Event(moment, 12, Queue[0], *player)
			//начинается сеанс первого игрока из очереди
			Computers[*player].SessionStart, Computers[*player].Occupied = moment, true
			// ушедший удаляется из списка активных клиентов, покинувший очередь удаляется из очереди
			delete(Players, client)
			delete(Queuers, Queue[0])
			delete(Customers, client)
			//очередь подрезается
			Queue = Queue[1:]
			return
		}

		if playing {
			// сеанс ушедшего игрока завершается и подсчитывается его время и чек, он удаляется из списка активных игроков
			Computers[*player].SessionEnd, Computers[*player].Occupied = moment, false
			Computers[*player].UseTime += moment.Sub(Computers[*player].SessionStart)
			Computers[*player].Revenue += calculateRevenue(Computers[*player])
			delete(Players, client)
		}
	}
}
