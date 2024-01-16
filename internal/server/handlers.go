package server

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"tg-bot-golang/internal/appmodels.go"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/keyboard/inline"
)

func (h *Handler) getProduct(ctx context.Context, b *bot.Bot, update *models.Update) {
	goodID, err := strconv.Atoi(update.Message.Text)

	if err != nil {
		h.log.Infof("Can't convert string value to int %s", err)
	}
	product, err := h.productservice.Handle(ctx, goodID)

	if (err != nil) || (reflect.DeepEqual(product, &appmodels.Product{})) {
		h.log.Infof("Can't get info about product %s \n", update.Message.Text)

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Не удалось найти товар с кодом %s. Введите новый код или вернитесь в главное меню /start ", update.Message.Text),
		})

	} else {
		inputFileData := product.Image.Max

		b.SendPhoto(ctx, &bot.SendPhotoParams{
			ChatID: update.Message.Chat.ID,
			Photo:  &models.InputFileString{Data: inputFileData},
		})

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   product.Description,
		})
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Введите новый код или вернитесь в главное меню /start ",
		})
	}
}

func (h *Handler) getRemains(ctx context.Context, b *bot.Bot, update *models.Update) {

	remains, err := h.axelotRepo.GetRemains(update.Message.Text)
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Can't get ramains: ",
		})
	}

	if len(remains) != 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("|%s|%s|%s|%s|%s|", "Cell", "Code", "Name", "EH", "Count"),
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("остатки товара или ячейки %s не найдены", update.Message.Text),
		})
	}

	for _, row := range remains {
		str := fmt.Sprintf("|%s|%s|%s|%s|%f|", row.Cell, row.Code, row.Name, row.EH, row.Count)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   str,
		})
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Введите новый код или вернитесь в главное меню /start ",
	})

}

func (h *Handler) getOrder(ctx context.Context, b *bot.Bot, update *models.Update, pref string) {
	var docNum string

	nameBtn := strings.TrimSpace(h.menuBtn[pref])
	lenBtn := len(nameBtn)

	switch lenBtn {
	case 9:
		docNum = fmt.Sprintf("%s-%06s", nameBtn, update.Message.Text)
	case 6:
		docNum = fmt.Sprintf("%s-%08s", nameBtn, update.Message.Text)
	case 5:
		docNum = fmt.Sprintf("%s-%08s", nameBtn, update.Message.Text)
	case 4:
		docNum = fmt.Sprintf("%s-%09s", nameBtn, update.Message.Text)
	case 3:
		docNum = fmt.Sprintf("%s-%11s", nameBtn, update.Message.Text)
	case 2:
		docNum = fmt.Sprintf("%s-%11s", nameBtn, update.Message.Text)
	}

	orders, err := h.axelotRepo.GetOrder(docNum)

	if err != nil {
		h.log.Errorf("Can't get order information", err)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Получен номер документа  %s", docNum),
	})

	if len(orders) != 0 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("|%s    |%s |", "Executor", "zone"),
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   fmt.Sprintf("Документ  %s не найден в ВМС системе", docNum),
		})
	}

	for _, val := range orders {
		text := fmt.Sprintf("|%s    |%s |", val.Executor, val.Zone)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   text,
		})
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "Введите новый код или вернитесь в главное меню /start ",
	})
}

func (h *Handler) setProperties(ctx context.Context, b *bot.Bot, update *models.Update, pref, value string) {

	code := getCode(pref)

	err := h.productservice.SetProperties(code, update.Message.Text, pref)
	if err != nil {
		h.log.Errorf("Can't change properties for code ", code)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось установить значения. Возврат в главное меню /start",
		})
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Данные записаны. Код  %s, значение %s. Возврат в главное меню /start", code, update.Message.Text),
	})
}

func (h *Handler) addUser(ctx context.Context, b *bot.Bot, update *models.Update) {
	err := h.userservice.CreateUser(ctx, update.Message.Text)
	if err != nil {
		h.log.Errorf("Can't add user ", update.Message.Text)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Не удалось создать пользователя. Возврат в главное меню /start",
		})
		return
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   fmt.Sprintf("Пользователь  %s добавлен. Возврат в главное меню /start", update.Message.Text),
	})
}

func (h *Handler) getAllUsers(ctx context.Context, b *bot.Bot, update *models.Update) {

	kb := inline.New(b)

	users, err := h.userservice.GetAllUsers(ctx)
	if err != nil {
		h.log.Errorf("Can't get users list ", err)
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Не удалось получить список пользователей. Возврат в главное меню /start",
		})
		return
	}
	for _, row := range users {
		kb = kb.Row().Button(row, []byte(row), h.onInlineKeyboardSelect)
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Chat.ID,
		Text:        "Выберете пользователя",
		ReplyMarkup: kb,
	})
}

func (h *Handler) recognizeUser(ctx context.Context, b *bot.Bot, update *models.Update) string {
	usrID := strconv.FormatInt(update.Message.Chat.ID, 10)
	name, err := h.userservice.GetUserByID(ctx, usrID)

	if err != nil {
		h.log.Errorf("Can't check user with id", usrID)
		return name
	}
	return name
}

func (h *Handler) compareUser(ctx context.Context, b *bot.Bot, update *models.Update) string {
	getname, err := h.userservice.CheckUser(ctx, update.Message.Text, strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		h.log.Errorf("Can't check user with id", update.Message.Text)
		return getname
	}
	return getname
}
func (h *Handler) isAdmin(ctx context.Context, b *bot.Bot, update *models.Update) bool {
	isAdm, err := h.userservice.CheckAdmin(ctx, strconv.FormatInt(update.Message.Chat.ID, 10))

	if err != nil {
		h.log.Errorf("Can't check user with id", update.Message.Text)

		return false
	}
	return isAdm

}
