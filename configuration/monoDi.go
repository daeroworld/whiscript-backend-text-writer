package configuration

import (
	"fmt"
	"net"
	"net/http"
	"text/writer/internal/api/application/controller"
	"text/writer/internal/api/application/handler"
	"text/writer/internal/api/domain/business"
	"text/writer/internal/api/domain/service"
	"text/writer/internal/api/infrastructure/client/voice"
	"text/writer/internal/api/infrastructure/repository/postgresql"

	sharedModel "github.com/daeroworld/shared/model"
	"github.com/gin-gonic/gin"

	pb "github.com/daeroworld/shared/proto/text"

	"github.com/daeroworld/shared/database"

	"google.golang.org/grpc"
)

type MonoContainer struct {
	Router            *gin.Engine
	Variable          *Variable
	PostgresqlWrapper *database.PostgresqlWrapper
	GRPCHandler       *handler.GRPCHandler
	ctrl              controller.IController
}

func (c *MonoContainer) GetHttpHandler() http.Handler {
	return c.Router.Handler()
}

func (c *MonoContainer) InitVariable() error {
	c.Variable = NewVariable()
	return nil
}

func (c *MonoContainer) SetRouter(router any) {
	return
}

func (c *MonoContainer) DefineDatabase(databaseWrappers ...any) error {
	c.PostgresqlWrapper = database.ConnectPostgresqlDatabase(c.Variable.Database)
	//ctnrMysqlWrapper := databaseWrappers[0].(*database.MysqlWrapper)
	err := c.PostgresqlWrapper.Driver.AutoMigrate(&sharedModel.Text{}, &sharedModel.TextConversion{})
	if err != nil {
		return err
	}
	//c.MysqlWrapper = ctnrMysqlWrapper
	return nil
}

func (c *MonoContainer) DefineRoute(router interface{}) error {

	return nil
}

func (c *MonoContainer) DefineGrpc() error {
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

func (c *MonoContainer) InitDependency(db interface{}) error {
	whisperBiz := business.NewWhisperBusiness(2)
	voiceClnt := voice.NewVoiceReaderClient("localhost", 25012)
	textRepo := postgresql.NewTextRepository(c.PostgresqlWrapper)
	conversionRepo := postgresql.NewConversionRepository(c.PostgresqlWrapper)
	svc := service.NewService(whisperBiz, voiceClnt, textRepo, conversionRepo)
	c.ctrl = controller.NewController(svc)
	//c.Handler = handler.NewHttpHandler(c.ctrl)
	c.GRPCHandler = handler.NewGRPCHandler(c.ctrl)
	return nil
}

func NewMonoContainer() *MonoContainer {
	ctnr := &MonoContainer{}
	ctnr.InitVariable()
	ctnr.DefineDatabase(nil)
	ctnr.InitDependency(ctnr.PostgresqlWrapper)
	ctnr.DefineGrpc()
	return ctnr
}
