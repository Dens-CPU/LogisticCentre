package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// Структуры продукта
type Product struct {
	ID       int
	Category string
	Weight   int
}

// Создание нового продукта
func NewProduct(id int, catehory string, weight int) Product {
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
func makeProduct(transit chan Product, categoryMap map[string][]Product, wgProduser *sync.WaitGroup) {

	//Запуск пула горутин
	for i := 0; i < DeliveryPeoples; i++ {
		wgProduser.Add(1)
		go func() {
			defer wgProduser.Done()
			for i := 0; i < rand.Intn(20)+1; i++ {
				category := rand.Intn(3) //Случайный выборо категории для созданного товара
				switch category {
				case 0:
					Transit(transit, categoryMap["Dresses"], "Dresses")
				case 1:
					Transit(transit, categoryMap["Car"], "Car")
				case 2:
					Transit(transit, categoryMap["Furniture"], "Furniture")
				}
				time.Sleep(time.Duration(rand.Intn(3)) * time.Second)
			}
		}()
	}
	//Горутина для закрытия канала transit
	go func() {
		wgProduser.Wait()
		close(transit)
	}()
}

// Функция отправления товара
func Transit(transit chan Product, queue []Product, category string) {
	for {
		id := rand.Intn(100) + 1 //Генерация случайного ID
		if !Checking(id, queue) {
			product := NewProduct(id, category, rand.Intn(100)+1)
			transit <- product
			log.Printf("(ОТПРРАВЛЕН)Товар №:%d; Категория:%s; Вес:%d кг.\n", product.ID, product.Category, product.Weight)
			return
		}
	}
}

var checkMutex = sync.RWMutex{}

// Функция проверки наличия ID в списке товаров
func Checking(id int, array []Product) bool {

	//Использование мьютекса для реальизаации синхронизации доступа к очереди

	for _, element := range array {
		checkMutex.RLock()
		if element.ID == id {
			return true
		}
		checkMutex.RUnlock()
	}
	return false
}

// Функция обработки товара и распределения его по зонам хранения
func Distribution(transit chan Product, categoryMap map[string][]Product, wgConsumer *sync.WaitGroup) {

	//Мьютекс для минхронизации доступа к мапе
	var mapMutex = sync.Mutex{}

	// Запуск работы обработчиков
	for i := 0; i < WarehouseWorkers; i++ {
		wgConsumer.Add(1)
		go func() {
			defer wgConsumer.Done()
			for product := range transit {
				log.Printf("(ПРИБЫЛ)Товар №:%d; Категория:%s; Вес:%d кг.\n", product.ID, product.Category, product.Weight)
				time.Sleep(time.Second)
				mapMutex.Lock()
				categoryMap[product.Category] = append(categoryMap[product.Category], product)
				log.Printf("Товар №%d весом %d кг был помещен в на склад категории %s.\n", product.ID, product.ID, product.Category)
				mapMutex.Unlock()
				time.Sleep(time.Second)

			}
			// fmt.Println("Канал закрыт")
		}()
	}
}
func main() {

	// Структура для синхронизации
	var wgProduser sync.WaitGroup
	var wgConsumer sync.WaitGroup

	//Карта категорий товаров
	CategoryMap := map[string][]Product{
		"Dresses":   make([]Product, 0),
		"Car":       make([]Product, 0),
		"Furniture": make([]Product, 0),
	}
	//Канал для передачи товаров. Вместительность канала создается равно количеству работников склада.
	transitChanel := make(chan Product, 10)
	Distribution(transitChanel, CategoryMap, &wgConsumer)
	makeProduct(transitChanel, CategoryMap, &wgProduser)
	wgConsumer.Wait()
	fmt.Println()
	k := 0
	for category := range CategoryMap {
		k = 0
		for i := 0; i < len(CategoryMap[category]); i++ {
			k++
		}
		fmt.Printf("В логистическом %d товаров категории %s\n", k, category)
	}
}
