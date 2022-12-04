package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"foodie/core"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var (
	routines         int64 = 75
	cycle            atomic.Uint64
	registerTime     atomic.Uint64
	createRecipeTime atomic.Uint64
	productsTime     atomic.Uint64
	recipesTime      atomic.Uint64
)

func main() {
	rand.Seed(time.Now().Unix())

	c := Client{
		url:    "http://127.0.0.1:13307",
		client: http.DefaultClient,
	}

	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(int(routines))

	for i := 0; i < int(routines); i++ {
		go func() {
			defer wg.Done()

			for ctx.Err() == nil {
				tstamp := time.Now()
				pp, err := c.Products(context.Background())
				if err != nil {
					fmt.Println(err)
					continue
				}
				productsTime.Add(uint64(time.Since(tstamp)))

				tstamp = time.Now()
				usr, err := c.Register(context.Background())
				if err != nil {
					fmt.Println(err)
					continue
				}
				registerTime.Add(uint64(time.Since(tstamp)))

				tstamp = time.Now()
				err = c.CreateRecipe(context.Background(), usr, pp)
				if err != nil {
					fmt.Println(err)
					continue
				}
				createRecipeTime.Add(uint64(time.Since(tstamp)))

				tstamp = time.Now()
				err = c.Recipes(context.Background())
				if err != nil {
					fmt.Println(err)
					continue
				}
				recipesTime.Add(uint64(time.Since(tstamp)))

				fmt.Println("completed cycle ", cycle.Add(1))
			}
		}()
	}

	terminationCh := make(chan os.Signal, 1)
	signal.Notify(terminationCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-terminationCh:
	}

	cancel()
	wg.Wait()

	logrus.Info("stopped")
	logrus.Info("average register time ", float64(registerTime.Load())/float64(cycle.Load()), "ns")
	logrus.Info("average create recipe time ", float64(createRecipeTime.Load())/float64(cycle.Load()), "ns")
	logrus.Info("average receive recipes time ", float64(recipesTime.Load())/float64(cycle.Load()), "ns")
	logrus.Info("average receive products time ", float64(productsTime.Load())/float64(cycle.Load()), "ns")
}

type Client struct {
	url    string
	client *http.Client
}

type User struct {
	Token    string
	Password string
	User     *core.User
}

func (c *Client) Products(ctx context.Context) ([]core.Product, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/products", c.url),
		http.NoBody,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var pp []core.Product

		if err := json.NewDecoder(resp.Body).Decode(&pp); err != nil {
			return nil, err
		}

		return pp, nil
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("invalid response %q (code: %d)", string(body), resp.StatusCode)
	}
}

func (c *Client) Register(ctx context.Context) (*User, error) {
	ui := core.UserInput{
		Name:     randSeq(rand.Intn(5) + 5),
		Password: randSeq(rand.Intn(5) + 6),
	}

	body, err := json.Marshal(ui)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/register", c.url),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		user := struct {
			User        *core.User `json:"user"`
			AccessToken string     `json:"access_token"`
		}{}

		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return nil, err
		}

		return &User{
			Token:    user.AccessToken,
			Password: ui.Password,
			User:     user.User,
		}, nil
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("invalid response %q (code: %d)", string(body), resp.StatusCode)
	}
}

func (c *Client) CreateRecipe(ctx context.Context, usr *User, pp []core.Product) error {
	ui := core.RecipeCore{
		Name:        randSeq(rand.Intn(5) + 5),
		ImageURL:    randSeq(rand.Intn(8) + 6),
		Description: randSeq(rand.Intn(8) + 20),
	}

	indexes := rand.Perm(len(pp))
	for i := 0; i < rand.Intn(len(pp)-2)+2; i++ {
		ui.Products = append(ui.Products, core.RecipeProduct{
			ProductID: pp[indexes[i]].ID,
			Quantity:  decimal.NewFromFloat(float64(rand.Intn(10000)/100) + 0.01),
		})
	}

	body, err := json.Marshal(ui)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/api/recipes", c.url),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+usr.Token)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("invalid response %q (code: %d)", string(body), resp.StatusCode)
	}
}

func (c *Client) Recipes(ctx context.Context) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/api/recipes", c.url),
		http.NoBody,
	)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	default:
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("invalid response %q (code: %d)", string(body), resp.StatusCode)
	}
}
