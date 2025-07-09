package service

import (
	"context"
	"fmt"
	"github.com/gogaruda/apperror"
	"github.com/irawankilmer/auth-service/internal/configs"
	"github.com/irawankilmer/auth-service/internal/model"
	"github.com/irawankilmer/auth-service/internal/repository"
	"github.com/irawankilmer/auth-service/pkg/mailer"
	"github.com/irawankilmer/auth-service/pkg/utils"
	"net/http"
	"time"
)

type EmailVerificationService interface {
	SendVerification(ctx context.Context, user model.UserModel) error
	VerifyToken(ctx context.Context, token string) error
}

type emailVerificationService struct {
	evRepo      repository.EmailVerificationRepository
	mail        *mailer.Mailer
	utilities   utils.Utility
	cfgMail     configs.EmailConfig
	userService UserService
}

func NewEmailVerificationService(
	ev repository.EmailVerificationRepository, m *mailer.Mailer, u utils.Utility,
	cm configs.EmailConfig, us UserService,
) EmailVerificationService {
	return &emailVerificationService{evRepo: ev, mail: m, utilities: u, cfgMail: cm, userService: us}
}

func (s *emailVerificationService) SendVerification(ctx context.Context, user model.UserModel) error {
	// generate token
	token, err := s.mail.GenerateRandom(32)
	if err != nil {
		return err
	}

	// create verification
	ev := &model.EmailVerificationModel{
		ID:        s.utilities.ULIDGenerate(),
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: time.Now().UTC().Add(30 * time.Minute),
	}
	if err := s.evRepo.Create(ctx, ev); err != nil {
		return err
	}

	// send mail
	url := fmt.Sprintf("%s?token=%s", s.cfgMail.FrontVerifyUrl, token)
	body := fmt.Sprintf("<p>Klik disini untuk verifikasi email: <a href='%s'>%s</a></p>", url, url)
	if err := s.mail.Send(user.Email, "Verifikasi Email", body); err != nil {
		return apperror.New("[SEND_EMAIL_VERIFICATION_VAILED]", "verifikasi email gagal dikirim", err, 505)
	}

	return nil
}

func (s *emailVerificationService) VerifyToken(ctx context.Context, token string) error {
	// periksa token
	ev, err := s.evRepo.FindByToken(ctx, token)
	if err != nil {
		return err
	}

	if ev.IsUsed {
		return apperror.New("[TOKEN_IS_USED]", "token sudah digunakan", err, http.StatusUnauthorized)
	}

	if time.Now().After(ev.ExpiresAt) {
		return apperror.New("[TOKEN_EXPIRED]", "token sudah kadaluwarsa", err, http.StatusUnauthorized)
	}

	// update verifikasi email users
	if err := s.userService.MarkEmailVerified(ctx, ev.UserID); err != nil {
		return err
	}

	// update is_used
	return s.evRepo.MarkAsUsed(ctx, ev.ID)
}
