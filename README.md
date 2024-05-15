# yadro-game-club-app
[![build](https://github.com/nikitads9/yadro-game-club-app/actions/workflows/build.yml/badge.svg)](https://github.com/nikitads9/yadro-game-club-app/actions/workflows/build.yml)
[![linters](https://github.com/nikitads9/yadro-game-club-app/actions/workflows/linter.yml/badge.svg)](https://github.com/nikitads9/yadro-game-club-app/actions/workflows/linter.yml)
![code lines](https://raw.githubusercontent.com/nikitads9/yadro-game-club-app/badges/.badges/main/lines.svg)
[![wakatime](https://wakatime.com/badge/user/018e5c64-a5fb-48a7-8d3a-b00fe4c56581/project/2782a224-4f20-4c7c-bc38-cbb4dedeed32.svg)](https://wakatime.com/badge/user/018e5c64-a5fb-48a7-8d3a-b00fe4c56581/project/2782a224-4f20-4c7c-bc38-cbb4dedeed32)

## Описание задачи

Требуется написать прототип системы, которая следит за работой компьютерного клуба, обрабатывает события и подсчитывает выручку за день и время занятости каждого стола.
Решение может быть реализовано на Golang.

Решением задания будет: файл или несколько файлов с исходным кодом программы на языке Golang (версия 1.19 и старше) с использованием go modules, инструкции по запуску и тестовые примеры (количество тестов – на усмотрение разработчика). 

Входные данные представляют собой текстовый файл. Файл указывается первым аргументом при запуске программы. Пример запуска программы: 
```
$ task.exe test_file.txt
```
Программа должна запускатьcя в Linux или Windows с использованием docker container-a (требуется написание Dockerfile). Требуется использование [стандартной библиотеки](https://pkg.go.dev/std). Использование любых сторонних библиотек, кроме стандартной, запрещено. В решении, кроме файлов с исходным кодом, требуется предоставить инструкции по запуску программы для проверки.

## Ход решения

Была написана программа, которая открывает файл, указанный как флаш `-path` при запуске исполняемого файла. Далее происходит чтение из этого файла с помощью стандартной библиотеки [bufio](https://pkg.go.dev/bufio). В решении используются проверки получаемых значений на корректность (неотрицательность, ненулевость), однако предполагается, что формат самих данных соблюден согласно обозначенному в задании. Для хранения данных о ждущих, вошедших и работающих клиентах используются три отдельные мапы (отображения), в случае с очередью и вошедшими клиентами значением в мапе выступает пустая структура, которая не требует выделения памяти. Для формирования самой очереди используется слайс, который нарезается и аппендится в зависимости от действия - клиент покидает очередь или встает в нее. Данные о столах с компьютерами реализованы в виде отображения, где ключом выступает идентификатор стола, а значением указатель на структуру `computer`, в которой хранятся:
- идентификатор стола
- время начала текущего сеанса
- время окончания сеанса
- имя последнего/текущего клиента
- время активного использования стола
- накопленная за день выручка по этому столу
Программа генерирует исходящие события в `stdout`, определяя требуемое действие, по текущему положению очереди, занятости столов и входящим событиям, которые определяются с помощью конструкции `switchcase`.

Часы работы предполагаются в рамках текущих суток, поэтому при вводе времени окончания рабочего дня, предшествующем времени его начала, программа завершает свое выполнение и указывает на эту ошибку.

## Инструкция по запуску

В проекте реализована возможность запуска как с помощью Docker, так и непосредственно в системе. 
Собрать исполняемый файл можно с помощью команды ниже (при наличии [Makefile](https://www.gnu.org/software/make/manual/make.html)). Если Makefile отсутствует, можно скопировать соответствующую команду из `Makefile`.
```bash
make build
```
Бинарник появляется в папке `bin` запускается с одним флагом - путь до файла с исходными данными для программы.
```bash
Использование ./events:
  -path string
        путь к файлу с исходными данными по итогам дня (default "./testdata/test.txt")
```
Аналогично для Windows. Исполняемый файл появляется там же.
```bash
make build-win
```
Для того, чтобы передать путь до файла с исходными данными, необходимо переименовать файл `.env.example` в `.env`.
Для запуска в контейнере достаточно выполнить команду (***Docker*** демон должен быть запущен):
```bash
make run
```
Эта команда создает образ контейнера с помощью ***Dockerfile***, расположенного в директории `deploy`. Далее она создает и запускает контейнер на основании этого образа, используя переменную окружения DATA_PATH, которая содержится в `.env` файле. При желании проверить работу программы на других данных, следует изменять данный параметр в этом файле.
Директория `testdata` смонтирована как ***Docker Volume*** для создаваемого контейнера, поэтому новые исходные данные следует располагать там.
Чтобы сменить этот параметр на новый, необходимо удалить контейнер и создать новый:
```bash
make docker-remove
make docker-run
```
Чтобы удалить только образ (контейнер должен быть удален):
```bash
make docker-delete
```
Удалить контейнер и образ можно одной командой:
```bash
make wipe
```

## Пример входных и выходных данных
| Входной файл      | Вывод в консоль              |
|-------------------|------------------------------|
| 3                 | 9:00                         |
| 09:00 19:00       | 08:48 1 client1              |
| 10                | 08:48 13   NotOpenYet        |
| 08:48 1 client1   | 09:41 1 client1              |
| 09:41 1 client1   | 09:48 1 client2              |
| 09:48 1 client2   | 09:52 3 client1              |
| 09:52 3 client1   | 09:52 13   ICanWaitNoLonger! |
| 09:54 2 client1 1 | 09:54 2 client1 1            |
| 10:25 2 client2 2 | 10:25 2 client2 2            |
| 10:58 1 client3   | 10:58 1 client3              |
| 10:59 2 client3 3 | 10:59 2 client3 3            |
| 11:30 1 client4   | 11:30 1 client4              |
| 11:35 2 client4 2 | 11:35 2 client4 2            |
| 11:45 3 client4   | 11:35 13   PlaceIsBusy       |
| 12:33 4 client1   | 11:45 3 client4              |
| 12:43 4 client2   | 12:33 4 client1              |
| 15:52 4 client4   | 12:33 12 client4   1         |
|                   | 12:43 4 client2              |
|                   | 15:52 4 client4              |
|                   | 19:00 11 client3             |
|                   | 19:00                        |
|                   | 1 70 05:58                   |
|                   | 2 30 02:18                   |
|                   | 3 90 08:01                   |

## Структура проекта

```
📦 yadro-game-club-app
├─ .env.example
├─ .github
│  └─ workflows
│     ├─ build.yml
│     ├─ linter.yml
│     └─ stats.yml
├─ .gitignore
├─ .golangci.pipeline.yaml
├─ Makefile
├─ README.md
├─ cmd
│  └─ events
│     └─ events.go
├─ deploy
│  └─ Dockerfile
├─ go.mod
├─ internal
│  ├─ format
│  │  └─ format.go
│  └─ process
│     ├─ process.go
│     ├─ read.go
│     └─ revenue.go
└─ testdata
   ├─ test.txt
   └─ test2.txt
```
