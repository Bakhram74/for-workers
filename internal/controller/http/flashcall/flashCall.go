package flashcall

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"time"

	"github.com/ShamilKhal/shgo/config"
	"github.com/ShamilKhal/shgo/internal/entity"
	"github.com/ShamilKhal/shgo/internal/service"
	"github.com/ShamilKhal/shgo/pkg/client/redis"
	"github.com/ShamilKhal/shgo/pkg/utils"
)

func Flashcall(ctx context.Context, config *config.Config, service *service.Service, phone string, userID string) (int, error) {
	
	return utils.Retry(3, 5*time.Second, func() error {

		pincode := utils.RandomNumbers(4)

		value := entity.PinValue{
			Pincode: pincode,
			Phone:   phone,
		}
		err := redis.SetPin(ctx, userID, value)
		if err != nil {
			return err
		}

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("public_key", config.Flashcall.PublicKey)
		_ = writer.WriteField("phone", phone)
		_ = writer.WriteField("campaign_id", config.Flashcall.CampaignID)
		_ = writer.WriteField("phone_suffix", pincode)
		writer.Close()

		url := "https://zvonok.com/manager/cabapi_external/api/v1/phones/flashcall/"

		request, err := http.NewRequest(http.MethodPost, url, payload)
		if err != nil {
			return err
		}
		request.Header.Set("Content-Type", writer.FormDataContentType())

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return err
		}

		defer response.Body.Close()

		s := response.StatusCode

		switch {
		case s >= 500:
			return fmt.Errorf("server error: %v", s)
		case s >= 400:
			return utils.Stop{Err: fmt.Errorf("client error: %v", s)}
		default:
			return nil
		}
	})
}
