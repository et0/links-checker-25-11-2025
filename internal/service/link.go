package service

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/et0/links-checker-25-11-2025/internal/config"
	"github.com/et0/links-checker-25-11-2025/internal/model"
	"github.com/et0/links-checker-25-11-2025/internal/pdf"
	"github.com/et0/links-checker-25-11-2025/internal/repository"
)

type LinkService interface {
	Check(ctx context.Context, urls []string) (*model.CheckResponse, error)
	CheckContinue(ctx context.Context) error
	Report(ctx context.Context, ids []int) ([]byte, error)
	Shutdown()
}

type linkService struct {
	db            repository.LinkRepository
	workerCount   int
	globWaitGroup sync.WaitGroup
	exit          bool // сигнал к выходу
}

func NewLinkService(db repository.LinkRepository, cfg *config.Config) LinkService {
	service := &linkService{
		db:          db,
		workerCount: cfg.App.WorkerCount,
	}

	return service
}

func (s *linkService) Check(ctx context.Context, urls []string) (*model.CheckResponse, error) {
	if s.isShutdown() {
		return nil, fmt.Errorf("Server is shutting down")
	}

	list := &model.LinkList{}

	// Сохраняем новую проверку в памяти
	if err := s.db.CreateList(list); err != nil {
		return nil, err
	}

	// Глобальная WG
	s.globWaitGroup.Add(1)
	defer s.globWaitGroup.Done()

	list.Status = model.Processing
	list.Links = s.checkLinkSync(urls)
	list.Status = model.Completed

	if err := s.db.UpdateList(list); err != nil {
		return nil, err
	}

	return &model.CheckResponse{
		ID:    list.ID,
		Links: *list.Links,
	}, nil
}

func (s *linkService) CheckContinue(ctx context.Context) error {
	if s.isShutdown() {
		return nil
	}

	for _, l := range s.db.GetAllList() {
		if l.Status != model.Pending {
			continue
		}

		s.globWaitGroup.Add(1)
		defer s.globWaitGroup.Done()

		urls := make([]string, 0, len(*l.Links))
		for _, link := range *l.Links {
			urls = append(urls, link.URL)
		}

		l.Status = model.Processing
		l.Links = s.checkLinkSync(urls)
		l.Status = model.Completed

		s.db.UpdateList(l)
	}

	return nil
}

func (s *linkService) Report(ctx context.Context, ids []int) ([]byte, error) {
	if s.isShutdown() {
		return nil, fmt.Errorf("Server is shutting down")
	}

	lists, err := s.db.GetList(ids)
	if err != nil {
		return nil, err
	}

	if len(lists) == 0 {
		return nil, fmt.Errorf("Empty link list")
	} else if len(lists) > 10 {
		return nil, fmt.Errorf("Limit 10 link id")
	}

	result, err := pdf.Create(lists)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *linkService) Shutdown() {
	s.exit = true

	// Ожидаем завершения глобальной Wait Group
	s.globWaitGroup.Wait()
}

func (s *linkService) isShutdown() bool {
	return s.exit
}

func (s *linkService) checkLinkSync(urls []string) *[]model.Link {
	var mu sync.Mutex

	// локальная WG для текущей проверки
	wg := &sync.WaitGroup{}

	semaphore := make(chan struct{}, s.workerCount)

	result := make([]model.Link, 0, len(urls))
	for _, u := range urls {
		wg.Add(1)

		go func() {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() {
				<-semaphore
			}()

			fmt.Println("Run ", u)

			status := s.checkLinkStatus(u)

			mu.Lock()
			result = append(result, model.Link{Status: status, URL: u})
			mu.Unlock()
		}()
	}

	wg.Wait()

	return &result
}

func (s *linkService) checkLinkStatus(url string) model.Status {
	if len(url) < 3 {
		return model.NotAvailable
	}

	if len(url) > 7 && (url[0:7] != "http://" && url[0:8] != "https://") {
		url = "https://" + url
	}

	time.Sleep(10 * time.Second)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Head(url)
	if err != nil {
		return model.NotAvailable
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return model.NotAvailable
	}

	return model.Available
}
