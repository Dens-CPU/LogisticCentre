package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// Структуры продукта
type Product struct {
	ID       int
	Category string
	Weight   float64
}

// Создание нового продукта
func NewProduct(id int, catehory string, weight float64) Product {
	return Product{ID: id, Category: catehory, Weight: weight}
}

const (
	WarehouseWorkers int = 5
	DeliveryPeoples  int = 10
)

// Функция для реальной работы рандомайзера в фунциях программы
func init() {
	rand.Seed(time.Now().Unix())
}

// Создания товара и отправление его в логистический центр
func makeProduct(transit chan Product, categoryMap map[string][]Product, wg *sync.WaitGroup) {
	//Запуск пула горутин
	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			category := 0
			switch category {
			case 0:
				for {
					id := rand.Intn(100) + 1 //Генерация случайного ID
					if !Checking(id, categoryMap["Dresses"]) {
						product := NewProduct(id, "Dresses", rand.Float64())
						transit <- product
						break
					}

				}
			}
		}()
	}
}

// Функция проверки наличия ID в списке товаров
func Checking(id int, array []Product) bool {
	var mapMutex = sync.RWMutex{}

	for _, element := range array {
		mapMutex.Lock()
		if element.ID == id {
			return true
		}
		mapMutex.Unlock()
	}
	return false
}
func main() {

	// Структура для синхронизации
	var wg sync.WaitGroup

	//Карта категорий товаров
	CategoryMap := map[string][]Product{
		"Dresses":   make([]Product, 0),
		"Car":       make([]Product, 0),
		"Furniture": make([]Product, 0),
	}
	//Канал для передачи товаров. Вместительность канала создается равно количеству работников склада.
	transitChanel := make(chan Product, 10)
	makeProduct(transitChanel, CategoryMap, &wg)
	wg.Wait()
	log.Println(CategoryMap)
}
