package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/totegamma/concurrent/x/core"
	"github.com/totegamma/concurrent/x/util"

	"github.com/bradfitz/gomemcache/memcache"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	version      = "unknown"
	buildMachine = "unknown"
	buildTime    = "unknown"
	goVersion    = "unknown"
)

func main() {

	e := echo.New()

	// Configファイルの読み込み
	config := util.Config{}
	configPath := os.Getenv("CONCURRENT_CONFIG")
	if configPath == "" {
		configPath = "/etc/concurrent/config.yaml"
	}
	err := config.Load(configPath)
	if err != nil {
		e.Logger.Fatal(err)
	}

	gwConf := GatewayConfig{}
	gwConfPath := os.Getenv("GATEWAY_CONFIG")
	if gwConfPath == "" {
		gwConfPath = "/etc/concurrent/gateway.yaml"
	}
	err = gwConf.Load(gwConfPath)
	if err != nil {
		e.Logger.Fatal(err)
	}

	log.Print("Concurrent ", version, " starting...")
	log.Print("Config loaded! I am: ", config.Concurrent.CCID)

	// Echoの設定
	e.HidePort = true
	e.HideBanner = true

	e.Use(middleware.Recover())

	if config.Server.EnableTrace {
		cleanup, err := setupTraceProvider(config.Server.TraceEndpoint, config.Concurrent.FQDN+"/ccgateway", version)
		if err != nil {
			panic(err)
		}
		defer cleanup()

		skipper := otelecho.WithSkipper(
			func(c echo.Context) bool {
				return c.Path() == "/metrics" || c.Path() == "/health"
			},
		)
		e.Use(otelecho.Middleware(config.Concurrent.FQDN, skipper))

		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				span := trace.SpanFromContext(c.Request().Context())
				c.Response().Header().Set("trace-id", span.SpanContext().TraceID().String())
				return next(c)
			}
		})
	}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/metrics" || c.Path() == "/health"
		},
		Format: `{"time":"${time_rfc3339_nano}",${custom},"remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","status":${status},` +
			`"error":"${error}","latency":${latency},"latency_human":"${latency_human}",` +
			`"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTagFunc: func(c echo.Context, buf *bytes.Buffer) (int, error) {
			span := trace.SpanFromContext(c.Request().Context())
			buf.WriteString(fmt.Sprintf("\"%s\":\"%s\"", "traceID", span.SpanContext().TraceID().String()))
			buf.WriteString(fmt.Sprintf(",\"%s\":\"%s\"", "spanID", span.SpanContext().SpanID().String()))
			return 0, nil
		},
	}))

	e.Use(echoprometheus.NewMiddlewareWithConfig(echoprometheus.MiddlewareConfig{
		Namespace: "ccgateway",
		LabelFuncs: map[string]echoprometheus.LabelValueFunc{
			"service": func(c echo.Context, err error) string {
				service := c.Response().Header().Get("cc-service")
				if service == "" {
					service = "unknown"
				}
				return service
			},
			"url": func(c echo.Context, err error) string {
				return "REDACTED"
			},
		},
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/metrics" || c.Path() == "/health"
		},
	}))

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             300 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Warn,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,                   // Enable color
		},
	)

	// Postrgresqlとの接続
	db, err := gorm.Open(postgres.Open(config.Server.Dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB, err := db.DB() // for pinging
	if err != nil {
		panic("failed to connect database")
	}
	defer sqlDB.Close()

	err = db.Use(tracing.NewPlugin(
		tracing.WithDBName("postgres"),
	))
	if err != nil {
		panic("failed to setup tracing plugin")
	}

	// Redisとの接続
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Server.RedisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err = redisotel.InstrumentTracing(
		rdb,
		redisotel.WithAttributes(
			attribute.KeyValue{
				Key:   "db.name",
				Value: attribute.StringValue("redis"),
			},
		),
	)
	if err != nil {
		panic("failed to setup tracing plugin")
	}

	mc := memcache.New(config.Server.MemcachedAddr)
	defer mc.Close()

	authService := SetupAuthService(db, rdb, mc, config)

	e.Use(authService.IdentifyIdentity)

	cors := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  []string{"*"},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		ExposeHeaders: []string{"trace-id"},
	})

	// プロキシ設定
	for _, service := range gwConf.Services {
		service := service

		if service.Gone {
			e.Any(service.Path, func(c echo.Context) error {
				return c.NoContent(http.StatusGone)
			})
			e.Any(service.Path+"/*", func(c echo.Context) error {
				return c.NoContent(http.StatusGone)
			})
			continue
		}

		targetUrl, err := url.Parse("http://" + service.Host + ":" + strconv.Itoa(service.Port))
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(targetUrl)

		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = targetUrl.Scheme
			req.URL.Host = targetUrl.Host
			if service.PreservePath {
				req.URL.Path = singleJoiningSlash(targetUrl.Path, req.URL.Path)
			} else {
				req.URL.Path = singleJoiningSlash(targetUrl.Path, strings.TrimPrefix(req.URL.Path, service.Path))
			}
			otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))
		}

		proxy.Transport = otelhttp.NewTransport(http.DefaultTransport)

		middlewares := []echo.MiddlewareFunc{}
		if service.InjectCors {
			middlewares = append(middlewares, cors)
		}

		handler := func(c echo.Context) error {
			c.Response().Header().Set("cc-service", service.Name)

			requesterType, ok := c.Get(core.RequesterTypeCtxKey).(int)
			if ok {
				c.Request().Header.Set(core.RequesterTypeHeader, strconv.Itoa(requesterType))
			}

			requesterId, ok := c.Get(core.RequesterIdCtxKey).(string)
			if ok {
				c.Request().Header.Set(core.RequesterIdHeader, requesterId)
			}

			requesterTag, ok := c.Get(core.RequesterTagCtxKey).(core.Tags)
			if ok {
				c.Request().Header.Set(core.RequesterTagHeader, requesterTag.ToString())
			}

			requesterDomain, ok := c.Get(core.RequesterDomainCtxKey).(string)
			if ok {
				c.Request().Header.Set(core.RequesterDomainHeader, requesterDomain)
			}

			requesterKeyDepath, ok := c.Get(core.RequesterKeyDepathKey).(string)
			if ok {
				c.Request().Header.Set(core.RequesterKeyDepathHeader, requesterKeyDepath)
			}

			requesterDomainTags, ok := c.Get(core.RequesterDomainTagsKey).(core.Tags)
			if ok {
				c.Request().Header.Set(core.RequesterDomainTagsHeader, requesterDomainTags.ToString())
			}

			requesterRemoteTags, ok := c.Get(core.RequesterRemoteTagsKey).(core.Tags)
			if ok {
				c.Request().Header.Set(core.RequesterRemoteTagsHeader, requesterRemoteTags.ToString())
			}

			proxy.ServeHTTP(c.Response(), c.Request())
			return nil
		}

		e.Any(service.Path, handler, middlewares...)
		e.Any(service.Path+"/*", handler, middlewares...)
	}

	e.GET("/", func(c echo.Context) (err error) {
		return c.HTML(http.StatusOK, `<!DOCTYPE html>
<html>
	<head>
		<title>ccgateway</title>
		<meta charset="utf-8">
		<link rel="icon" href="`+config.Profile.Logo+`">
	</head>
	<body>
		<h1>Concurrent Domain - `+config.Concurrent.FQDN+`</h1>
		Yay! You're on ccgateway!<br>
		You might looking for <a href="https://concurrent.world">concurrent.world</a>.<br>
		This domain is currently registration: `+config.Concurrent.Registration+`<br>
		<h2>Information</h2>
		CCID: <br>`+config.Concurrent.CCID+`<br>
		PUBKEY: <br>`+config.Concurrent.PublicKey[:len(config.Concurrent.PublicKey)/2]+`<br>`+
			config.Concurrent.PublicKey[len(config.Concurrent.PublicKey)/2:]+`<br>
		<h2>Services</h2>
		<ul>
		`+func() string {
			var services string
			for _, service := range gwConf.Services {
				services += `<li><a href="` + service.Path + `">` + service.Name + `</a></li>`
			}
			return services
		}()+`
		</ul>
	</body>
</html>
`)
	})

	e.GET("/services", func(c echo.Context) (err error) {
		services := make(map[string]ServiceInfo)
		for _, service := range gwConf.Services {
			services[service.Name] = ServiceInfo{
				Path: service.Path,
			}
		}
		return c.JSON(http.StatusOK, services)
	}, cors)

	e.GET("/health", func(c echo.Context) (err error) {
		ctx := c.Request().Context()

		err = sqlDB.Ping()
		if err != nil {
			return c.String(http.StatusInternalServerError, "db error")
		}

		err = rdb.Ping(ctx).Err()
		if err != nil {
			return c.String(http.StatusInternalServerError, "redis error")
		}

		return c.String(http.StatusOK, "ok")
	})

	e.GET("/metrics", echoprometheus.NewHandler())

	e.Start(":8080")
}

func setupTraceProvider(endpoint string, serviceName string, serviceVersion string) (func(), error) {

	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)

	if err != nil {
		return nil, err
	}
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)

	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tracerProvider)

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	cleanup := func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		if err := tracerProvider.Shutdown(ctx); err != nil {
			log.Printf("Failed to shutdown tracer provider: %v", err)
		}
	}
	return cleanup, nil
}

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}
