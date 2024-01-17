package service

import (
	"context"

	tgAPI "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog"

	"github.com/7Maliko7/telegram-bot/internal/scenario"
	"github.com/7Maliko7/telegram-bot/internal/script"
	"github.com/7Maliko7/telegram-bot/internal/trigger"
	"github.com/7Maliko7/telegram-bot/pkg/cache"
	"github.com/7Maliko7/telegram-bot/pkg/db"
	"github.com/7Maliko7/telegram-bot/pkg/messenger/telegram"
	"github.com/7Maliko7/telegram-bot/pkg/metrics"
)

var (
	finalKB = telegram.Keyboard{
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
)

type Service struct {
	logger    *zerolog.Logger
	tg        *telegram.Telegram
	messenger *telegram.Messenger
	cache     cache.Cacher
	db        db.Databaser
	scenario  *scenario.Scenario
}

func New(logger *zerolog.Logger, tg *telegram.Telegram, cache cache.Cacher, db db.Databaser, scenario *scenario.Scenario) *Service {
	s := &Service{
		logger:    logger,
		tg:        tg,
		messenger: tg.MakeMessenger(),
		cache:     cache,
		db:        db,
		scenario:  scenario,
	}

	s.logger.Info().Msg("handlers registration...")
	tg.RegisterScenario(s.handleScenario)
	s.logger.Debug().Int("steps", s.scenario.StepCount()).Msgf("%v registered", "scenario")
	tg.RegisterHandler("ping", s.handlePing)
	s.logger.Debug().Msgf("%v registered", "ping")
	tg.RegisterHandler("start", s.handleStart)
	s.logger.Debug().Msgf("%v registered", "start")
	tg.RegisterHandler("close", s.handleClose)
	s.logger.Debug().Msgf("%v registered", "close")
	tg.RegisterHandler("Restart", s.handleStart)
	s.logger.Debug().Msgf("%v registered", "Restart")
	tg.RegisterHandler("Continue", s.handleStart)
	s.logger.Debug().Msgf("%v registered", "Continue")
	tg.RegisterHandler("Learn more", s.handleAbout)
	s.logger.Debug().Msgf("%v registered", "Learn more")
	tg.RegisterHandler("See more", s.handleEP)
	s.logger.Debug().Msgf("%v registered", "See more")
	tg.RegisterHandler("About us", s.handleBand)
	s.logger.Debug().Msgf("%v registered", "About us")
	tg.RegisterHandler("Learn about us", s.handleStartAbout)
	s.logger.Debug().Msgf("%v registered", "Learn about us")
	s.logger.Info().Msg("handlers registered")

	return s
}

func (s *Service) Run() {
	s.tg.Listen()
}

func (s *Service) handleStart(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "start").Msg("handler invoked")

	step, err := s.cache.Get(u.Message.Chat.ID)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	if step == nil {
		user, err := s.db.GetUserByDialogId(context.TODO(), u.Message.Chat.ID)
		if err != nil {
			s.logger.Err(err).Send()
			return false
		}
		if user == nil {
			err := s.db.AddUser(context.TODO(), u.Message.Chat.ID)
			s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("save user")
			if err != nil {
				s.logger.Err(err).Send()
				return false
			}
			metrics.OpsChat.Inc()
		} else {
			s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int64("user_id", user.User_id).Msg("retrieved user")
		}

		s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("start scenario")
		return s.handleScenario(u)
	}

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", *step.Int).Msg("cached step")

	_, err = s.cache.Clear(u.Message.Chat.ID)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("cache cleared")

	scenarioLen := s.scenario.StepCount()
	_, err = s.cache.SaveArray(u.Message.Chat.ID, 0, "/start", &scenarioLen)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", 0).Str("command", "/start").Msg("step cached")

	_, err = s.cache.SaveArray(u.Message.Chat.ID, 1, "repeated", &scenarioLen)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", 1).Str("command", "repeated").Msg("step cached")

	_, err = s.cache.SaveInt(u.Message.Chat.ID, 1)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", 1).Msg("current step cached")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("start scenario")
	return s.handleScenario(u)
}

func (s *Service) handleScenario(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "scenario").Msg("handler invoked")

	var (
		stepID   int
		text     string
		keyboard telegram.Keyboard
		newStep  *script.Step
	)

	cacheStep, err := s.cache.Get(u.Message.Chat.ID)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	if cacheStep == nil {
		step := s.scenario.Step(0)

		text = step.Text
		keyboard = step.Actions

		scenarioLen := s.scenario.StepCount()
		_, err = s.cache.SaveArray(u.Message.Chat.ID, stepID, u.Message.Text, &scenarioLen)
		if err != nil {
			s.logger.Err(err).Send()
			return false
		}
		s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", stepID).Str("command", u.Message.Text).Msg("step cached")
	} else {
		s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", *cacheStep.Int).Msg("cached step")

		stepID = *cacheStep.Int
		step := s.scenario.Step(stepID)

		text = step.Text
		keyboard = step.Actions

		for _, acts := range step.Actions {
			for _, act := range acts {
				if u.Message.Text == act {
					s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", stepID).Str("command", u.Message.Text).Msg("action accepted")

					stepID++
					newStep = s.scenario.Step(stepID)

					text = newStep.Text
					keyboard = newStep.Actions

					_, err = s.cache.SaveArray(u.Message.Chat.ID, stepID, u.Message.Text, nil)
					if err != nil {
						s.logger.Err(err).Send()
						continue
					}
					s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", stepID).Str("command", u.Message.Text).Msg("step cached")
				}
			}
		}
	}

	if keyboard == nil {
		_, err = s.messenger.CloseKeyboard(u.Message.Chat.ID, text)
		if err != nil {
			s.logger.Err(err).Send()
			return false
		}
	} else {
		_, err = s.messenger.SendKeyboard(u.Message.Chat.ID, text, keyboard)
		if err != nil {
			s.logger.Err(err).Send()
			return false
		}
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")

	if newStep != nil {
		s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("triggers invoked")
		for _, tr := range newStep.Triggers {
			chVal, err := s.cache.Get(u.Message.Chat.ID)
			if err != nil {
				s.logger.Err(err).Send()
				continue
			}
			s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("cache retrieved")

			s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("trigger", tr).Msg("trigger processing...")
			text, keyboard := trigger.Process(tr, chVal, u)
			s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("trigger", tr).Msg("trigger processed")

			if keyboard != nil {
				_, err = s.messenger.SendKeyboard(u.Message.Chat.ID, *text, keyboard)
				if err != nil {
					s.logger.Err(err).Send()
					return false
				}
			} else {
				_, err = s.messenger.SendText(u.Message.Chat.ID, *text)
				if err != nil {
					s.logger.Err(err).Send()
					return false
				}
			}
			s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")
		}
	}

	if s.scenario.StepCount()-1 == stepID {
		metrics.OpsScenario.Inc()
	}

	_, err = s.cache.SaveInt(u.Message.Chat.ID, stepID)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Int("step_id", stepID).Msg("current step cached")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "scenario").Msg("handler finished")
	return true
}

func (s *Service) handleClose(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "close").Msg("handler invoked")

	_, err := s.messenger.CloseKeyboard(u.Message.Chat.ID, "OK")
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "close").Msg("handler finished")
	return true
}

func (s *Service) handlePing(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "ping").Msg("handler invoked")

	text := "pong"

	_, err := s.messenger.SendText(u.Message.Chat.ID, text)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "ping").Msg("handler finished")
	return true
}

func (s *Service) handleAbout(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "about").Msg("handler invoked")
	text := "Description of bot"

	_, err := s.messenger.SendKeyboard(u.Message.Chat.ID, text, finalKB)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "about").Msg("handler finished")
	return true
}

func (s *Service) handleEP(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "ep").Msg("handler invoked")
	text := "Description of our work"

	_, err := s.messenger.SendKeyboard(u.Message.Chat.ID, text, finalKB)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "ep").Msg("handler finished")
	return true
}

func (s *Service) handleBand(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "band").Msg("handler invoked")
	text := "About us"

	_, err := s.messenger.SendKeyboard(u.Message.Chat.ID, text, finalKB)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "band").Msg("handler finished")
	return true
}

func (s *Service) handleStartAbout(u tgAPI.Update) bool {
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "start about").Msg("handler invoked")
	text := "About us"

	kb := telegram.Keyboard{
		{
			"Continue",
		},
	}

	_, err := s.messenger.SendKeyboard(u.Message.Chat.ID, text, kb)
	if err != nil {
		s.logger.Err(err).Send()
		return false
	}
	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Msg("send tg message")

	s.logger.Debug().Int64("chat_id", u.Message.Chat.ID).Str("command", u.Message.Text).Str("handler", "band").Msg("handler finished")
	return true
}
