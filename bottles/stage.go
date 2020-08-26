package bottles


type StageFunc func(*Bottle) error

func StoreMessageStage(pool *MessagePool) StageFunc {
	return func(b *Bottle) error {
		if err := pool.Add(b.Message); err != nil {
			return err
		}
		return nil
	}
}

func ValidateTokenStage(pool *TokenPool) StageFunc {
	return func(b *Bottle) error {
		if err := pool.Use(b.Token); err != nil {
			return err
		}
		return nil
	}
}
