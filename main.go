package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache — структура, представляющая потокобезопасный кэш.
// Включает RWMutex для конкурентного доступа и саму map.
type Cache struct {
	mu   sync.RWMutex
	data map[int]int
}

// NewCache — конструктор, создаёт кэш с заданной ёмкостью.
func NewCache(capacity int) *Cache {
	return &Cache{data: make(map[int]int, capacity)}
}

// Get — безопасное конкурентное чтение из кэша.
// Использует RLock, чтобы допустить параллельное чтение.
func (c *Cache) Get(key int) int {
	c.mu.RLock()         // Блокируем на чтение
	defer c.mu.RUnlock() // Освобождаем после чтения
	value, ok := c.data[key]
	if !ok {
		fmt.Println("Cache miss") // Ключ не найден
		return 0
	}
	return value
}

// Set — безопасная конкурентная запись в кэш.
// Использует Lock, так как запись должна быть эксклюзивной.
func (c *Cache) Set(key int, value int) {
	c.mu.Lock()         // Блокируем на запись
	c.data[key] = value // Записываем значение
	c.mu.Unlock()       // Освобождаем мьютекс
}

func main() {
	cache := NewCache(30) // Создаём кэш на 30 элементов

	var wg sync.WaitGroup // WaitGroup для ожидания завершения горутин

	// Горутина записи в кэш
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 30; i++ {
			cache.Set(i, i) // Записываем пары ключ-значение: 0→0, 1→1, ...
		}
	}()

	// Горутина чтения из кэша
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 30; i++ {
			time.Sleep(100 * time.Millisecond) // Ждём, чтобы дать шанс писателю сначала заполнить
			value := cache.Get(i)              // Пытаемся прочитать значение по ключу i
			fmt.Printf("value of map[%d] is %d\n", i, value)
		}
	}()

	wg.Wait() // Ожидаем завершения обеих горутин
}
