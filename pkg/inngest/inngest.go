package inngest

import (
	"draw/internal/db/repo"
	"draw/pkg/config"

	"github.com/gin-gonic/gin"
	"github.com/inngest/inngestgo"
)

type Inngest struct {
	client       inngestgo.Client
	awsConfig    *config.AWSConfig
	geminiConfig *config.GeminiConfig
	queries      *repo.Queries
}

func NewInngest(awsConfig *config.AWSConfig,
	geminiConfig *config.GeminiConfig,
	queries *repo.Queries,
) (*Inngest, error) {
	client, err := inngestgo.NewClient(inngestgo.ClientOpts{
		AppID: "core",
		Dev:   inngestgo.BoolPtr(true),
	})
	if err != nil {
		return nil, err
	}

	i := &Inngest{
		client:       client,
		awsConfig:    awsConfig,
		geminiConfig: geminiConfig,
		queries:      queries,
	}

	err = i.RegisterFunctions()
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Inngest) Handler() gin.HandlerFunc {
	inngestHandler := i.client.Serve()
	return func(c *gin.Context) {
		inngestHandler.ServeHTTP(c.Writer, c.Request)
	}
}
