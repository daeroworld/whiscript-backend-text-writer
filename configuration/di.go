package configuration

import (
	"fmt"
	"net"
	"net/http"
	"text/writer/internal/api/application/controller"
	"text/writer/internal/api/application/handler"
	"text/writer/internal/api/domain/business"
	indexBiz "text/writer/internal/api/domain/business/index"
	"text/writer/internal/api/domain/service"
	"text/writer/internal/api/infrastructure/client/voice"
	"text/writer/internal/api/infrastructure/repository/postgresql"

	sharedModel "github.com/daeroworld/shared/model"
	"github.com/gin-gonic/gin"

	pb "github.com/daeroworld/shared/proto/text"

	"github.com/daeroworld/shared/database"

	"google.golang.org/grpc"
)

type Container struct {
	Router            *gin.Engine
	Variable          *Variable
	PostgresqlWrapper *database.PostgresqlWrapper
	GRPCHandler       *handler.GRPCHandler
	ctrl              controller.IController
}

func (c *Container) GetHttpHandler() http.Handler {
	return c.Router.Handler()
}

func (c *Container) InitVariable() error {
	c.Variable = NewVariable()
	return nil
}

func (c *Container) SetRouter(router any) {
	return
}

func (c *Container) DefineDatabase(databaseWrappers ...any) error {
	c.PostgresqlWrapper = database.ConnectPostgresqlDatabase(c.Variable.Database)
	//ctnrMysqlWrapper := databaseWrappers[0].(*database.MysqlWrapper)
	err := c.PostgresqlWrapper.Driver.AutoMigrate(&sharedModel.Text{}, &sharedModel.TextConversion{})
	if err != nil {
		return err
	}
	//c.MysqlWrapper = ctnrMysqlWrapper
	return nil
}

func (c *Container) DefineRoute(router interface{}) error {

	return nil
}

func (c *Container) DefineGrpc() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", c.Variable.Api.Ip, c.Variable.Api.Port))
	if err != nil {
		return fmt.Errorf("failed to listen : %d, %w", c.Variable.Api.Port, err)
	}
	server := grpc.NewServer()
	pb.RegisterTextWriterServer(server, c.GRPCHandler)
	go func() {
		if servErr := server.Serve(lis); servErr != nil {
			return
		}
	}()
	return nil
}

func (c *Container) InitDependency(db interface{}) error {
	whisperBiz := business.NewWhisperBusiness(2)
	indexBiz := indexBiz.NewIndexBusiness()
	voiceClnt := voice.NewVoiceReaderClient(c.Variable.VoiceReaderApi.Ip, c.Variable.VoiceReaderApi.Port)
	textRepo := postgresql.NewTextRepository(c.PostgresqlWrapper)
	conversionRepo := postgresql.NewConversionRepository(c.PostgresqlWrapper)
	svc := service.NewService(indexBiz, whisperBiz, voiceClnt, textRepo, conversionRepo)
	c.ctrl = controller.NewController(svc)
	//c.Handler = handler.NewHttpHandler(c.ctrl)
	c.GRPCHandler = handler.NewGRPCHandler(c.ctrl)
	return nil
}

func NewContainer() *Container {
	ctnr := &Container{}
	ctnr.InitVariable()
	ctnr.DefineDatabase(nil)
	ctnr.InitDependency(ctnr.PostgresqlWrapper)

	return ctnr
}
