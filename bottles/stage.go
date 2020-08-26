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

func AddMessageStage(pool *MessagePool) StageFunc {
	return func(b *Bottle) error {
		message, err := pool.Get()
		if err != nil {
			return err
		}
		b.Message = message

		return nil
	}
}

func AddTokenStage(pool *TokenPool) StageFunc {
	return func(b *Bottle) error {
		size := 10
		tokenStr := GenerateRandomString(size)
		token := &Token{
			Str: &tokenStr,
		}
		for pool.Add(token) != nil {
			tokenStr = GenerateRandomString(size)
			token = &Token{
				Str: &tokenStr,
			}
		}
		b.Token = token
		return nil
	}
}
