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
	SendVerification(ctx context.Context, user *model.UserModel, urlTo, actionType string, duration time.Duration) (string, error)
	VerifyToken(ctx context.Context, token string) error
	CheckToken(ctx context.Context, token string) (*model.UserModel, error)
}

type emailVerificationService struct {
	evRepo    repository.EmailVerificationRepository
	mail      *mailer.Mailer
	utilities utils.Utility
	cfgMail   configs.EmailConfig
	userRepo  repository.UserRepository
}

func NewEmailVerificationService(
	ev repository.EmailVerificationRepository, m *mailer.Mailer, u utils.Utility,
	cm configs.EmailConfig, ur repository.UserRepository,
) EmailVerificationService {
	return &emailVerificationService{evRepo: ev, mail: m, utilities: u, cfgMail: cm, userRepo: ur}
}

func (s *emailVerificationService) SendVerification(ctx context.Context, user *model.UserModel, urlTo, actionType string, duration time.Duration) (string, error) {
	// generate token
	token, err := s.mail.GenerateRandom(64)
	if err != nil {
		return "", err
	}

	// create verification
	ev := &model.EmailVerificationModel{
		ID:         s.utilities.ULIDGenerate(),
		UserID:     user.ID,
		Token:      token,
		ExpiresAt:  time.Now().UTC().Add(duration),
		ActionType: actionType,
	}
	if err := s.evRepo.Create(ctx, ev); err != nil {
		return "", err
	}

	// send mail
	url := fmt.Sprintf("%s/%s?token=%s", s.cfgMail.FrontVerifyUrl, urlTo, token)
	body := fmt.Sprintf(`
	<h2>Verifikasi Email Anda</h2>
	<p>Halo,</p>
	<p>Email Anda telah terdaftar. Untuk menyelesaikan proses pendaftaran, silakan verifikasi alamat email Anda dengan mengklik tombol di bawah ini:</p>
	<p><a href='%s' style='
		display: inline-block;
		padding: 10px 20px;
		background-color: #4CAF50;
		color: white;
		text-decoration: none;
		border-radius: 5px;
		font-weight: bold;
	'>Verifikasi Email</a></p>
	<p>Jika tombol di atas tidak bekerja, salin dan tempel URL berikut ke browser Anda:</p>
	<p><code>%s</code></p>
	<p>Link ini akan kadaluarsa dalam 24 jam.</p>
	<p>Salam hangat,<br><strong>Tim Support %s</strong></p>
`, url, url, "Sekolah Kita")

	if err := s.mail.Send(user.Email, "Verifikasi Email", body); err != nil {
		return "", apperror.New("[SEND_EMAIL_VERIFICATION_FAILED]", "verifikasi email gagal dikirim", err, 505)
	}

	return token, nil
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
		return apperror.New("[TOKEN_EXPIRED]", "token sudah kadaluwarsa", err, http.StatusUnauthorized).WithResponseStatus("expired")
	}

	// cek user by ID
	user, err := s.userRepo.FindByID(ctx, ev.UserID)
	if err != nil {
		return err
	}

	// update verifikasi email
	if err := s.userRepo.UpdateEmailVerified(ctx, user); err != nil {
		return err
	}

	// update is_used
	return s.evRepo.MarkAsUsed(ctx, ev.ID)
}

func (s *emailVerificationService) CheckToken(ctx context.Context, token string) (*model.UserModel, error) {
	// cek token
	ev, err := s.evRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}
	if ev.IsUsed {
		return nil, apperror.New("[TOKEN_IS_USED]", "user sudah aktif", err, http.StatusUnauthorized)
	}
	if time.Now().After(ev.ExpiresAt) {
		return nil, apperror.New("[TOKEN_EXPIRED]", "token sudah kadaluwarsa", err, http.StatusUnauthorized).WithResponseStatus("expired")
	}

	// cek user by ID
	user, err := s.userRepo.FindByID(ctx, ev.UserID)
	if err != nil {
		return nil, err
	}

	var userModel model.UserModel
	userModel = model.UserModel{
		ID:    user.ID,
		Email: user.Email,
	}

	// update is_used
	if err := s.evRepo.MarkAsUsed(ctx, ev.ID); err != nil {
		return nil, err
	}
	return &userModel, nil
}
