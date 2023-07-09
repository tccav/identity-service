package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/extra/redisotel/v9"
	goredis "github.com/redis/go-redis/v9"
	httpswagger "github.com/swaggo/http-swagger"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
	"github.com/twmb/franz-go/plugin/kotel"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"moul.io/chizap"

	_ "github.com/tccav/identity-service/api"
	"github.com/tccav/identity-service/pkg/config"
	"github.com/tccav/identity-service/pkg/domain/identities/idusecases"
	"github.com/tccav/identity-service/pkg/gateways/httpserver"
	"github.com/tccav/identity-service/pkg/gateways/kafka"
	"github.com/tccav/identity-service/pkg/gateways/opentelemetry"
	"github.com/tccav/identity-service/pkg/gateways/postgres"
	"github.com/tccav/identity-service/pkg/gateways/redis"
)

var (
	AppVersion = "unknown"
	GoVersion  = "unknown"
	Time       = "unknown"
)

// @title Identity Service API
// @version 1.0
// @description Service responsible for identity management of the Aluno Online's system.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url https://github.com/tccav
// @contact.email pedroyremolo@gmail.com
// @license.name No License
// @license.url https://choosealicense.com/no-permission/
func main() {
	logConfig := zap.NewProductionConfig()
	logConfig.DisableStacktrace = true

	logger, err := logConfig.Build(
		zap.Fields(
			zap.String("version", AppVersion),
			zap.String("go_version", GoVersion),
			zap.String("build_time", Time),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to initialize logger: %s", err))
	}
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	logger.Info("application init started, configs will be loaded")

	configs, err := config.LoadConfigs()
	if err != nil {
		logger.Error("failed to load configs", zap.Error(err))
		return
	}

	logger = logger.With(zap.String("environment", configs.Telemetry.Environment))

	tpClose, err := opentelemetry.InitProvider(logger, AppVersion, configs.Telemetry)
	if err != nil {
		logger.Error("failed to initialize tracer", zap.Error(err))
		return
	}
	defer tpClose()

	ctx := context.Background()
	dbCfg, err := pgxpool.ParseConfig(configs.DB.URL())
	if err != nil {
		logger.Error("failed to parse db url", zap.Error(err))
	}

	dbCfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(ctx, dbCfg)
	if err != nil {
		logger.Error("failed to start db", zap.Error(err))
		return
	}
	defer pool.Close()
	logger.Info("db conn pool fetched")

	redisOptions := goredis.Options{
		Addr: configs.MemoryDB.URL(),
	}
	if configs.MemoryDB.User != "" {
		redisOptions.Username = configs.MemoryDB.User
		redisOptions.Password = configs.MemoryDB.Password
	}

	redisClient := goredis.NewClient(&redisOptions)
	err = redisClient.Ping(ctx).Err()
	if err != nil {
		logger.Error("failed to fetch memory db conn", zap.Error(err))
		return
	}
	defer redisClient.Close()
	logger.Info("memory db conn pool fetched")

	if err = redisotel.InstrumentTracing(redisClient); err != nil {
		logger.Error("failed to init memory db tracing", zap.Error(err))
		return
	}
	if err = redisotel.InstrumentMetrics(redisClient); err != nil {
		logger.Error("failed to init memory db metrics", zap.Error(err))
		return
	}

	// Create a new kotel tracer.
	tracerOpts := []kotel.TracerOpt{
		kotel.TracerProvider(otel.GetTracerProvider()),
		kotel.TracerPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{})),
	}
	tracer := kotel.NewTracer(tracerOpts...)

	// Create a new kotel service.
	kotelOps := []kotel.Opt{
		kotel.WithTracer(tracer),
	}
	kotelService := kotel.NewKotel(kotelOps...)

	kOpts := []kgo.Opt{
		kgo.SeedBrokers(configs.Kafka.URL()),
		kgo.WithHooks(kotelService.Hooks()...),
	}
	if configs.Kafka.User != "" {
		kOpts = append(kOpts, kgo.SASL(plain.Auth{
			User: configs.Kafka.User,
			Pass: configs.Kafka.Password,
		}.AsMechanism()))
	}

	kafkaClient, err := kgo.NewClient(kOpts...)
	if err != nil {
		logger.Error("unable to connect to kafka broker", zap.Error(err))
		return
	}
	defer kafkaClient.Close()

	err = kafkaClient.Ping(ctx)
	if err != nil {
		logger.Error("kafka broker unreachable", zap.Error(err))
		return
	}
	logger.Info("kafka client created")

	producer := kafka.NewProducer(kafkaClient)

	studentsProducer := kafka.NewStudentsProducer(producer)

	repository := postgres.NewStudentsRepository(pool)
	tokenRepository := redis.NewTokensRepository(redisClient)

	useCase := idusecases.NewRegisterUseCase(repository, studentsProducer)
	authUseCase := idusecases.NewStudentJWTAuthenticator(repository, tokenRepository, configs.Auth)

	studentsHandler := httpserver.NewStudentsHandler(useCase, logger)
	authHandler := httpserver.NewAuthenticationHandler(logger, authUseCase)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(chizap.New(logger, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	}))
	if configs.Swagger.Enabled {
		router.Get("/swagger/*", httpswagger.Handler())
	}
	router.MethodFunc(http.MethodPost, "/v1/identities/students", studentsHandler.RegisterStudent)
	router.MethodFunc(http.MethodPost, "/v1/identities/students/login", authHandler.AuthenticateStudent)
	router.MethodFunc(http.MethodPost, "/v1/identities/students/verify-auth", authHandler.VerifyAuthentication)
	router.Get("/healthcheck", httpserver.Healthcheck)
	logger.Info("handlers and routes configured")

	handler := otelhttp.NewHandler(router, "", otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
		return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
	}))

	server := http.Server{
		Addr:              fmt.Sprintf(":%d", configs.API.Port),
		Handler:           handler,
		ReadTimeout:       configs.API.ReadTimeout,
		ReadHeaderTimeout: configs.API.ReadTimeout,
		WriteTimeout:      configs.API.WriteTimeout,
		IdleTimeout:       configs.API.IdleTimeout,
	}

	notifyContext, stop := signal.NotifyContext(ctx, os.Kill, os.Interrupt)
	defer stop()

	go func(sigCtx context.Context) {
		<-sigCtx.Done()
		logger.Info("shutdown signal received")
		shutdownCtx, c := context.WithTimeout(ctx, 30*time.Second)
		defer c()
		shutdownErr := server.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			logger.Error("server shutdown failed", zap.Error(shutdownErr))
			return
		}
		logger.Info("bye bye!")
	}(notifyContext)

	logger.Info("server will be started", zap.String("addr", server.Addr))
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error("server listening has failed", zap.Error(err))
		return
	}
}
