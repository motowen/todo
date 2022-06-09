package test

import (
	"fmt"
	"github.com/jarcoal/httpmock"
	"log"
	"os"
	"testing"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/app/router"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/config"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/http/client"
	"viveportengineering.visualstudio.com/Viveport-Core/_git/go-base.git/internal/pkg/logger"
)

func TestMain(m *testing.M) {
	//ctx := context.Background()

	/*
		redisContainer, setupRedisErr := container.SetupRedis(ctx)
		if setupRedisErr != nil {
			log.Fatalf("SetupRedis Fail, %s\n", setupRedisErr)
		}
		defer redisContainer.Terminate(ctx)
	*/

	/*
		mongoContainer, setupMongoErr := container.SetupMongo(ctx)

		if setupMongoErr != nil {
			log.Fatalf("SetupMongo Fail, %s\n", setupMongoErr)
		}
		defer mongoContainer.Terminate(ctx)
	*/

	Setup()
	defer Close()

	httpmock.ActivateNonDefault(client.Get().GetClient())

	r := m.Run()

	if r == 0 && testing.CoverMode() != "" {
		c := testing.Coverage() * 100
		l := 0.00
		fmt.Println("=================================================")
		fmt.Println("||               Coverage Report               ||")
		fmt.Println("=================================================")
		fmt.Printf("Cover mode: %s\n", testing.CoverMode())
		fmt.Printf("Coverage  : %.2f %% (Threshold: %.2f %%)\n\n", c, l)
		if c < l {
			fmt.Println("[Tests passed but coverage failed]")
			r = -1
		}
	}

	os.Exit(r)
}

func Setup() {
	var err error

	if err = config.Setup(); err != nil {
		log.Fatal(err)
	}

	/*
		if err = cache.GetInstance().Setup(cache.Config{
			Type:         config.Env.RedisType,
			EndpointList: config.Env.RedisEndpointList,
			Password:     config.Env.RedisPassword,
		}); err != nil {
			log.Fatalf("cache Setup, error:%v", err)
		}
	*/

	/*
		if err = database.GetInstance().Setup(database.Config{
			URI: config.Env.MongoURI,
		}); err != nil {
			log.Fatalf("database Setup, error:%v", err)
		}
	*/

	/*
		if err = postgres.GetInstance().Setup(postgres.Config{
			Username:                config.Env.PostgresUsername,
			Password:                config.Env.PostgresPassword,
			Host:                    config.Env.PostgresHost,
			Port:                    config.Env.PostgresPort,
			TableName:               config.Env.PostgresName,
			MinConnSize:             config.Env.PostgresMinConnSize,
			MaxConnSize:             config.Env.PostgresMaxConnSize,
			MaxConnIdleTimeBySecond: time.Duration(config.Env.PostgresMaxConnIdleTimeBySecond),
			MaxConnLifetimeBySecond: time.Duration(config.Env.PostgresMaxConnLifeTimeBySecond),
		}); err != nil {
			log.Fatalf("postgres Setup, error:%v", err)
		}
	*/

	/*
		if err = queue.GetInstance().Setup(queue.Config{
			Url: config.Env.NatsUrl,
		}); err != nil {
			log.Fatalf("queue Setup, error:%v", err)
		}
	*/

	/*
		if err = search.GetInstance().Setup(search.Config{
			Url:         config.Env.ElasticsearchUrl,
			IndexPrefix: config.Env.ElasticsearchIndexPrefix,
		}); err != nil {
			log.Fatalf("search Setup, error:%v", err)
		}
	*/

	if err = logger.Setup(config.Env.LogLevel); err != nil {
		log.Fatal(err)
	}

	if err = client.Setup(); err != nil {
		log.Fatal(err)
	}

	if err = router.Setup(); err != nil {
		log.Fatal(err)
	}
}

func Close() {
}
