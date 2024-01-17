package trigger

import (
	"fmt"
	"strings"

	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/7Maliko7/telegram-bot/pkg/cache"
	"github.com/7Maliko7/telegram-bot/pkg/messenger/telegram"
	"github.com/7Maliko7/telegram-bot/pkg/test_result"
)

const (
	TriggerProcessHistory = "PROCESS_HISTORY"
	TriggerFinalize       = "FINALIZE"
	TriggerDoEnd          = "DO_END"
)

func Process(trigger string, c *cache.Value, u tgAPI.Update) (*string, [][]string) {
	var (
		text     string
		keyboard telegram.Keyboard
	)
	switch trigger {
	case TriggerProcessHistory:
		text = processHistory(c)
	case TriggerFinalize:
		text = finalize()
	case TriggerDoEnd:
		text = doEnd()
		keyboard = telegram.Keyboard{
			{
				"Restart",
			},
			{
				"Learn more",
			},
			{
				"See more",
			},
			{
				"About us",
			},
		}
	}

	return &text, keyboard
}

func processHistory(value *cache.Value) string {
	var (
		result string
		arr    [25]bool
		count  int
	)
	for i, v := range value.Array {
		result = fmt.Sprintf("%v\n%v", result, v)

		if i-1 >= 0 {
			switch v {
			case "Yes":
				arr[i-1] = true
				count++
			}
		}
	}

	res := test_result.MakeResult(arr)

	m := res.GetChooseSubTotal()
	mechanism := strings.Join(m, "\n")

	return fmt.Sprintf("%v", mechanism)

}

func finalize() string {
	return "Result"
}

func doEnd() string {
	return "Congratulations"
}
