package server

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/logger"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/go-telegram/ui/dialog"
)

const (
	prefBtnOrder      = "BtnOrder"
	prefBtnProperties = "BtnProperties"
)

var (
	dialogNodesADM = []dialog.Node{
		{ID: "admin", Text: "Администрирование", Keyboard: [][]dialog.Button{{{Text: "Добавить", NodeID: "BtnAddUser"}, {Text: "Удалить", NodeID: "BtnDeleteUser"}}}},
	}
	dialogNodes = []dialog.Node{
		// {ID: "start", Text: "Start Node", Keyboard: [][]dialog.Button{{{Text: "Товар", NodeID: "Product"}, {Text: "Заказ", NodeID: "Order"}, {Text: "Доставка", NodeID: "BtnDelivery"}}, {{Text: "Go Telegram UI", URL: "https://github.com/go-telegram/ui"}}}},
		{ID: "start", Text: "Start Node", Keyboard: [][]dialog.Button{{{Text: "Товар", NodeID: "Product"}, {Text: "Заказ", NodeID: "Order"}}, {{Text: "Go Telegram UI", URL: "https://github.com/go-telegram/ui"}}}},
		{ID: "BarCode", Text: "Выберете единицу", Keyboard: [][]dialog.Button{{{Text: "шт", NodeID: "Btn796"}, {Text: "м", NodeID: "Btn006"}, {Text: "упак", NodeID: "Btn778"}, {Text: "Бухта", NodeID: "Btn123"}, {Text: "Go to start", NodeID: "start"}}}},
		{ID: "Product", Text: "Выберете действие", Keyboard: [][]dialog.Button{{{Text: "Остатки", NodeID: "BtnRemains"}, {Text: "ГХ", NodeID: "Properties"}, {Text: "Info", NodeID: "BtnInfo"}, {Text: "Go to start", NodeID: "start"}}}},
		{ID: "Properties", Text: "Выберете действие", Keyboard: [][]dialog.Button{{{Text: "ШК", NodeID: "BarCode"}, {Text: "Вес", NodeID: "BtnWheight"}, {Text: "Объем", NodeID: "BtnVolume"}, {Text: "Go to start", NodeID: "start"}}}},
		{ID: "Order", Text: "Выберете действие", Keyboard: [][]dialog.Button{{
			{Text: "АЭ", NodeID: "Btn11"},
			{Text: "Чек", NodeID: "Btn12"},
			{Text: "РНк", NodeID: "Btn14"}},
			{{Text: "РНк/у", NodeID: "Btn21"},
				{Text: "РНр/у", NodeID: "Btn22"},
				{Text: "РнБ", NodeID: "Btn23"},
				{Text: "Прм/у", NodeID: "Btn24"}},
			{{Text: "ИМ1", NodeID: "Btn31"},
				{Text: "ИМ2", NodeID: "Btn32"},
				{Text: "ИМ3", NodeID: "Btn33"},
				{Text: "ИМ4", NodeID: "Btn34"}},
			{{Text: "ИП1", NodeID: "Btn41"},
				{Text: "ИП2", NodeID: "Btn42"},
				{Text: "ИП3", NodeID: "Btn43"},
				{Text: "ИП4", NodeID: "Btn44"}},
			{
				{Text: "Go to start", NodeID: "start"}}}},
		{ID: "BtnInfo", Text: "Введите код товара"},
		{ID: "BtnRemains", Text: "Введите код товара или ячейки"},
	}
)

type Handler struct {
	productservice productservice
	userservice    userservice
	axelotRepo     axelot
	log            logger.Logger
	cfg            *config.Config
	menuBtn        map[string]string
	point          point
}

func NewHandler(service productservice, userservice userservice, axelotRepo axelot, point point, log logger.Logger, cfg *config.Config) *Handler {
	return &Handler{
		productservice: service,
		userservice:    userservice,
		axelotRepo:     axelotRepo,
		log:            log,
		cfg:            cfg,
		menuBtn:        setBtnOrders(),
		point:          point,
	}
}

func (h *Handler) callbackHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var nameBtn, menuBtn string

	if update.CallbackQuery == nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Главное меню /start ",
		})
	}
	datastr := update.CallbackQuery.Data[16:]
	btnBloc := getBtnBloc(datastr)

	if btnBloc == "" {
		menuBtn = datastr
		nameBtn = datastr
	} else {
		menuBtn = btnBloc
		if strings.Contains(datastr, "Btn796") || strings.Contains(datastr, "Btn006") || strings.Contains(datastr, "Btn778") || strings.Contains(datastr, "Btn123") {
			nameBtn = fmt.Sprintf("%s/BarCode/%s", btnBloc, datastr)
		} else {
			nameBtn = fmt.Sprintf("%s/%s", btnBloc, datastr)
		}

	}

	h.point.Add(ctx, update.CallbackQuery.Sender.FirstName, nameBtn)

	switch menuBtn {
	case "BtnAddUser":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Введите имя пользователя: ",
		})
	case "BtnDeleteUser":
		h.getAllUsers(ctx, b, update)
	case "BtnInfo":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Введите Код товара: ",
		})
	case "BtnRemains":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Введите Код товара или ячейки: ",
		})
	case "BtnOrder":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Введите номер документа: ",
		})
	case "BtnProperties":
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Введите код товара: ",
		})
	default:
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.CallbackQuery.Message.Chat.ID,
			Text:   "Возврат в главное меню /start ",
		})
	}

}

func (h *Handler) defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	var point, pref string

	d := dialog.New(dialogNodes)
	if update.Message == nil {
		fmt.Println("Доделать диалог")
		return
	}

	dadm := dialog.New(dialogNodesADM)
	if update.Message == nil {
		fmt.Println("Доделать диалог")
		return
	}

	switch update.Message.Text {
	case "/start":
		h.point.Del(ctx, update.Message.Chat.FirstName)
		name := h.recognizeUser(ctx, b, update)
		if name != "" {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   fmt.Sprintf("Hello, %s", name),
			})
			d.Show(ctx, b, strconv.FormatInt(update.Message.Chat.ID, 10), "start")
		} else {
			h.point.Add(ctx, update.Message.Chat.FirstName, "BtnCheck")
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Введите имя пользователя",
			})
		}

	case "/admin":
		if h.isAdmin(ctx, b, update) {
			dadm.Show(ctx, b, strconv.FormatInt(update.Message.Chat.ID, 10), "admin")
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: update.Message.Chat.ID,
				Text:   "Пользователь не является администратором. Возврат в главное меню /start",
			})
		}

	default:
		menuStr := h.point.Get(ctx, update.Message.From.FirstName)

		switch {
		case strings.Contains(menuStr, prefBtnOrder):
			point = prefBtnOrder
			pref = strings.Replace(menuStr, prefBtnOrder+"/", "", -1)
		case strings.Contains(menuStr, prefBtnProperties):
			point = prefBtnProperties
			pref = strings.Replace(menuStr, prefBtnProperties+"/", "", -1)
		default:
			point = menuStr
		}

		switch point {
		case "BtnAddUser":
			h.addUser(ctx, b, update)
		case "BtnCheck":
			getname := h.compareUser(ctx, b, update)
			if getname == update.Message.Text {
				d.Show(ctx, b, strconv.FormatInt(update.Message.Chat.ID, 10), "start")
			} else {
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   "Данного пользователя не существует. Возврат в главное меню /start",
				})
			}
		case "BtnInfo":
			h.getProduct(ctx, b, update)
		case "BtnRemains":
			h.getRemains(ctx, b, update)
		case "BtnOrder":
			h.getOrder(ctx, b, update, pref)
		case "BtnProperties":
			cntSl := strings.Count(menuStr, "/")
			switch cntSl {
			case 1:
				h.point.Add(ctx, update.Message.From.FirstName, fmt.Sprintf("%s/%s", menuStr, update.Message.Text))

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: update.Message.Chat.ID,
					Text:   fmt.Sprintf("Введите %s:", textForUser(menuStr)),
				})
			case 2:
				menuStr = h.point.Get(ctx, update.Message.From.FirstName)
				if strings.Contains(menuStr, "BarCode") {
					h.point.Add(ctx, update.Message.From.FirstName, fmt.Sprintf("%s/%s", menuStr, update.Message.Text))

					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: update.Message.Chat.ID,
						Text:   fmt.Sprintf("Введите %s:", textForUser(menuStr)),
					})
				} else {
					h.setProperties(ctx, b, update, menuStr, update.Message.Text)
				}
			case 3:
				menuStr = h.point.Get(ctx, update.Message.From.FirstName)
				h.setProperties(ctx, b, update, menuStr, update.Message.Text)
			}
		}
	}
}

func (h *Handler) onInlineKeyboardSelect(ctx context.Context, b *bot.Bot, mes *models.Message, data []byte) {
	err := h.userservice.DeleteUser(ctx, string(data))
	if err != nil {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "Не удалось удалить пользователя: " + string(data),
		})
	} else {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: mes.Chat.ID,
			Text:   "Пользователь удален: " + string(data),
		})
	}

}
