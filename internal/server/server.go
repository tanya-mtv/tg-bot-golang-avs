package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"tg-bot-golang/internal/config"
	"tg-bot-golang/internal/db/cashredis"
	"tg-bot-golang/internal/db/dialogredis"
	"tg-bot-golang/internal/db/mssql"
	"tg-bot-golang/internal/db/wms"
	"tg-bot-golang/internal/httpclient"
	"tg-bot-golang/internal/logger"
	"tg-bot-golang/internal/service"

	"github.com/go-telegram/bot"
)

type Server struct {
	log logger.Logger
	cfg *config.Config
	bot *bot.Bot
}

func NewServer(log logger.Logger, cfg *config.Config) *Server {
	return &Server{
		log: log,
		cfg: cfg,
	}
}

func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	sqlClient := mssql.NewMSSQLDB(s.cfg.MSSQL)
	defer sqlClient.Close() // nolint: errcheck
	s.log.Infof("MSSQL connected: %+v", sqlClient.PingContext(ctx))

	sqlAxelotClient := wms.NewAxelotConn(s.cfg.Axelot)
	defer sqlAxelotClient.Close() // nolint: errcheck
	s.log.Infof("Axelot connected: %+v", sqlAxelotClient.PingContext(ctx))

	redisClient := cashredis.NewUniversalRedisClient(s.cfg.Redis)
	defer redisClient.Close() // nolint: errcheck
	s.log.Infof("Redis connected: %+v", redisClient.PoolStats())

	sqlRepo := mssql.NewmsSQLStorage(s.log, s.cfg, sqlClient)
	axelotRepo := wms.NewmsAxelotStorage(s.log, s.cfg, sqlAxelotClient)

	redisRepo := cashredis.NewRedisStorage(s.log, s.cfg, redisClient)
	raecRepo := httpclient.NewRaecRepository(s.log, s.cfg)

	ps := service.NewProductService(s.log, s.cfg, sqlRepo, redisRepo, raecRepo)
	usr := mssql.NewmsUsersStorage(s.log, s.cfg, sqlClient)
	point := dialogredis.NewDialogStorage(s.log, s.cfg, redisClient)

	handl := NewHandler(ps, usr, axelotRepo, point, s.log, s.cfg)
	opts := []bot.Option{

		bot.WithDefaultHandler(handl.defaultHandler),
		bot.WithCallbackQueryDataHandler("Btn", bot.MatchTypeContains, handl.callbackHandler),
	}

	telegramBotToken := s.cfg.BotToken
	b, err := bot.New(telegramBotToken, opts...)

	if err != nil {
		s.log.Fatal(err)
		return err
	}

	s.bot = b

	go func() {
		s.log.Infof("Bot is listening!", "AVSrobot")
		s.bot.Start(ctx)
	}()

	<-ctx.Done()
	return nil
}
