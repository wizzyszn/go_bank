package service

import (
	"fmt"
	"time"

	"github.com/wizzyszn/go_bank/db"
	"github.com/wizzyszn/go_bank/models"
	"github.com/wizzyszn/go_bank/repository"
	"github.com/wizzyszn/go_bank/utils"
)

type AuthService struct {
	db              *db.DB
	accountRepo     *repository.AccountRepository
	sessionRepo     *repository.SessionRepository
	sessionDuration time.Duration
}

func NewAuthService(database *db.DB, accountRepo *repository.AccountRepository, sessionRepo *repository.SessionRepository, sessionDuration time.Duration) *AuthService {

	return &AuthService{
		db:              database,
		accountRepo:     accountRepo,
		sessionRepo:     sessionRepo,
		sessionDuration: sessionDuration,
	}
}

func (s *AuthService) RegisterAccount(req *models.CreateAccountRequest) (*models.AccountResponse, error) {

	//Validate
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}
	if err := utils.ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}
	if err := utils.ValidateName(req.FirstName, "first_name"); err != nil {
		return nil, err
	}
	if err := utils.ValidateName(req.LastName, "last_name"); err != nil {
		return nil, err
	}

	// Sanitize
	req.Email = utils.SanitizeString(req.Email)
	req.FirstName = utils.SanitizeString(req.FirstName)
	req.LastName = utils.SanitizeString(req.LastName)

	exists, err := s.accountRepo.EmailExists(req.Email)

	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("Email already in use")
	}

	passwordHash, err := utils.HashPassword(req.Password)

	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	account, err := s.accountRepo.Create(req.Email, passwordHash, req.FirstName, req.LastName)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return account.ToResponse(), nil

}

func (s *AuthService) Login(req models.LoginAccountRequest) (*models.LoginResponse, error) {
	if err := utils.ValidateEmail(req.Email); err != nil {
		return nil, err
	}

	if err := utils.ValidateRequired(req.Password, "password"); err != nil {
		return nil, err
	}

	account, err := s.accountRepo.GeyByEmail(req.Email)

	if err != nil {
		return nil, fmt.Errorf("Invalid email or password")
	}

	if account.Status != models.AccountStatusActice {
		return nil, fmt.Errorf("account is %s", account.Status)
	}

	if err := utils.CheckPassword(req.Password, account.PasswordHash); err != nil {
		return nil, fmt.Errorf("Invalid Email or Password")
	}

	sessionID, err := utils.GenerateSessionID()

	if err != nil {
		return nil, fmt.Errorf("failed to generate session: %w", err)
	}

	session, err := s.sessionRepo.Create(sessionID, account.ID, time.Now().Add(s.sessionDuration))

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &models.LoginResponse{
		Account:   account.ToResponse(),
		SessionID: session.ID,
		ExpiresAt: session.ExpiresAt,
	}, nil

}

func (s *AuthService) Logout(sessionID string) error {
	if err := utils.ValidateRequired(sessionID, "session_id"); err != nil {
		return err
	}

	if err := s.sessionRepo.Delete(sessionID); err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}
	return nil
}

func (s *AuthService) LogoutAll(accountID int) error {
	if err := s.sessionRepo.DeleteAccountByID(accountID); err != nil {
		return fmt.Errorf("failed to logout all sessions: %w", err)
	}
	return nil
}

func (s *AuthService) ValidateSession(sessionID string) (*models.Account, error) {
	session, err := s.sessionRepo.GetByID(sessionID)

	if err != nil {
		return nil, fmt.Errorf("invalid session")
	}

	if session.IsExpired() {
		s.sessionRepo.Delete(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	account, err := s.accountRepo.GetByID(session.AccountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	if account.Status != models.AccountStatusActice {
		return nil, fmt.Errorf("account is %s", account.Status)
	}

	return account, nil
}

func (s *AuthService) GetAccount(accountID int) (*models.AccountResponse, error) {
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}
	return account.ToResponse(), nil
}

func (s *AuthService) UpdateAccount(accountID int, req *models.UpdateAccountRequest) (*models.AccountResponse, error) {
	account, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("account not found")
	}

	firstName := account.FirstName
	lastName := account.LastName

	if req.FirstName != "" {
		if err := utils.ValidateName(req.FirstName, "first_name"); err != nil {
			return nil, err
		}
		firstName = utils.SanitizeString(req.FirstName)
	}

	if req.LastName != "" {
		if err := utils.ValidateName(req.LastName, "last_name"); err != nil {
			return nil, err
		}
		lastName = utils.SanitizeString(req.LastName)
	}

	if err := s.accountRepo.Update(accountID, firstName, lastName); err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	updated, err := s.accountRepo.GetByID(accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch updated account: %w", err)
	}

	return updated.ToResponse(), nil
}

// CleanupExpiredSessions removes all expired sessions
// Intended to be called on a schedule (e.g. every hour via a goroutine in main.go)
func (s *AuthService) CleanupExpiredSessions() (int, error) {
	count, err := s.sessionRepo.DeleteExpired()
	if err != nil {
		return 0, fmt.Errorf("failed to cleanup sessions: %w", err)
	}
	return count, nil
}
