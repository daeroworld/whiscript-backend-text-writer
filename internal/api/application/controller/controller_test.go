package controller_test

import (
	"os"
	"testing"
	"time"

	"text/writer/internal/api/application/controller"
	"text/writer/internal/api/domain/business"
	"text/writer/internal/api/domain/service"
	"text/writer/internal/api/infrastructure/client/voice"
	"text/writer/internal/api/infrastructure/repository/postgresql"

	"github.com/daeroworld/shared/configuration"
	"github.com/daeroworld/shared/database"
	sharedModel "github.com/daeroworld/shared/model"
)

func measure(t *testing.T, name string, fn func()) {
	start := time.Now()
	fn()
	t.Logf("%s took %v", name, time.Since(start))
}

var (
	testCtrl   *controller.Controller
	testFileId = ""
)

func TestMain(m *testing.M) {
	whisperBiz := business.NewWhisperBusiness(2)
	voiceClnt := voice.NewVoiceReaderClient("localhost", 25012)

	//ctnrMysqlWrapper := databaseWrappers[0].(*database.MysqlWrapper)

	
	err := postgresqlWrapper.Driver.AutoMigrate(&sharedModel.Text{}, &sharedModel.TextConversion{})
	if err != nil {
		return
	}
	textRepo := postgresql.NewTextRepository(postgresqlWrapper)
	conversionRepo := postgresql.NewConversionRepository(postgresqlWrapper)

	svc := service.NewService(whisperBiz, voiceClnt, textRepo, conversionRepo)
	testCtrl = controller.NewController(svc)

	// Run all tests & benchmarks
	code := m.Run()
	os.Exit(code)
}

func BenchmarkCreateSync(b *testing.B) {
	b.ReportAllocs()

	start := time.Now()

	_, err := testCtrl.CreateSync(testFileId)
	if err != nil {
		b.Fatalf("Create error: %v", err)
	}

	latency := time.Since(start)
	b.ReportMetric(float64(latency.Milliseconds()), "ms/op")
}

func BenchmarkCreate(b *testing.B) {
	b.ReportAllocs()

	start := time.Now()

	_, err := testCtrl.Create(testFileId)
	if err != nil {
		b.Fatalf("Create error: %v", err)
	}

	latency := time.Since(start)
	b.ReportMetric(float64(latency.Milliseconds()), "ms/op")
}
