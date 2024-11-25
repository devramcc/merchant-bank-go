package repository

import (
	"encoding/json"
	"log"
	"os"
)

type WhitelistAccessTokenRepository struct {
	filePath string
}

func NewWhitelistAccessTokenRepository(filePath string) *WhitelistAccessTokenRepository {
	return &WhitelistAccessTokenRepository{
		filePath: filePath,
	}
}

func (r *WhitelistAccessTokenRepository) LoadWhitelistTokens() ([]string, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var tokens []string
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}

func (r *WhitelistAccessTokenRepository) SaveWhitelistToken(token string) error {
	tokens, err := r.LoadWhitelistTokens()
	if err != nil {
		return err
	}

	tokens = append(tokens, token)
	data, err := json.MarshalIndent(tokens, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *WhitelistAccessTokenRepository) RemoveWhitelistToken(token string) error {
	tokens, err := r.LoadWhitelistTokens()
	if err != nil {
		return err
	}

	newTokens := []string{}
	for _, t := range tokens {
		if t != token {
			newTokens = append(newTokens, t)
		}
	}

	data, err := json.MarshalIndent(newTokens, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func (r *WhitelistAccessTokenRepository) IsTokenWhitelisted(token string) bool {
	tokens, err := r.LoadWhitelistTokens()
	if err != nil {
		log.Printf("Error loading whitelist tokens: %v", err)
		return false
	}

	for _, t := range tokens {
		if t == token {
			return true
		}
	}

	return false
}
