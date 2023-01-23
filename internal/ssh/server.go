package ssh

import (
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/robherley/snips.sh/internal/config"
	"github.com/robherley/snips.sh/internal/db"
)

func New(cfg *config.Config, db *db.DB) (*Server, error) {
	sessionHandler := &SessionHandler{db}

	return wish.NewServer(
		wish.WithMaxTimeout(MaxTimeout),
		wish.WithIdleTimeout(IdleTimeout),
		wish.WithAddress(cfg.SSHAddress()),
		wish.WithHostKeyPath(cfg.SSH.HostKeyPath),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		wish.WithPasswordAuth(func(ctx ssh.Context, password string) bool {
			// accept pw auth so we can display a helpful message
			return true
		}),
		// note: middleware is evaulated in reverse order
		wish.WithMiddleware(
			sessionHandler.HandleFunc,
			AssignUser(db),
			BlockIfNoPublicKey,
			WithLogger,
		),
	)
}