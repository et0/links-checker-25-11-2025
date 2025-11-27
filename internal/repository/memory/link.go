package memory

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/et0/links-checker-25-11-2025/internal/model"
)

type Link struct {
	filepath string
	mu       *sync.Mutex
	data     map[int]*model.LinkList
}

func NewLink(filepath string) (*Link, error) {
	repo := Link{
		filepath: filepath,
		mu:       &sync.Mutex{},
		data:     make(map[int]*model.LinkList),
	}

	if err := repo.loadJSON(); err != nil {
		return nil, err
	}

	return &repo, nil
}

func (l *Link) loadJSON() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	file, err := os.Open(l.filepath)
	if err != nil {
		// если файл ещё не создан
		if os.IsNotExist(err) {
			return l.saveJSON()
		}

		return err
	}
	defer file.Close()

	var data []*model.LinkList
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	for _, d := range data {
		l.data[d.ID] = d
	}

	return nil
}

func (l *Link) saveJSON() error {
	file, err := os.Create(l.filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	data := make([]*model.LinkList, 0, len(l.data))
	for _, d := range l.data {
		data = append(data, d)
	}

	return json.NewEncoder(file).Encode(data)
}

func (l *Link) CreateList(list *model.LinkList) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Println("New id ", len(l.data)+1)
	list.ID = len(l.data) + 1
	list.CreatedAt = time.Now()
	list.Status = model.Pending

	l.data[list.ID] = list

	return l.saveJSON()
}

func (l *Link) UpdateList(list *model.LinkList) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.data[list.ID]; !exists {
		return fmt.Errorf("Wrong list ID")
	}

	l.data[list.ID] = list

	return l.saveJSON()
}

func (l *Link) GetList(ids []int) ([]*model.LinkList, error) {
	lists := make([]*model.LinkList, 0, len(ids))
	for _, id := range ids {
		list, exists := l.data[id]
		if !exists {
			continue
		}
		lists = append(lists, list)
	}

	return lists, nil
}

func (l *Link) GetAllList() []*model.LinkList {
	l.mu.Lock()
	defer l.mu.Unlock()

	lists := make([]*model.LinkList, 0, len(l.data))
	for _, d := range l.data {
		lists = append(lists, d)
	}

	return lists
}
