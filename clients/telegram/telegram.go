package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdateMethod   = "getUpdates"
	sendMessageMethod = "sendMessage"
)

// создает Client
func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

// клиент будет заниматься двумя вещами
// получать update (новых сообщений)
// отправкой собственных сообщений пользователю

// получает update
func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	//параметры запроса
	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequests(getUpdateMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("can't decode in UpdateResponse: %w", err)
	}

	return res.Result, nil
}

// отправляет сообщения
func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.doRequests(sendMessageMethod, q)
	if err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func (c *Client) doRequests(method string, query url.Values) ([]byte, error) {
	//сформируем url на который будет отправляться запрос
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	//подготавливаем запрос
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("can't do request: %w", err)
	}

	req.URL.RawQuery = query.Encode()
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response: %w", err)
	}

	return body, nil
}
