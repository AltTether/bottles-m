package bottles

import (
	"fmt"
	"time"
	"testing"
	
	"github.com/stretchr/testify/assert"
)


func TestValidateTokenStage(t *testing.T) {
	pool := NewTokenPool(10 * time.Second)
	stage := ValidateTokenStage(pool)

	tokenStr := "test"
	token := &Token{
		Str: &tokenStr,
	}
	bottle := &Bottle{
		Token: token,
	}

	pool.Add(token)
	err := stage(bottle)
	assert.Nil(t, err)
}

func TestValidateTokenStageInvalidToken(t *testing.T) {
	pool := NewTokenPool(10 * time.Second)
	stage := ValidateTokenStage(pool)

	tokenStr := "test"
	token := &Token{
		Str: &tokenStr,
	}
	bottle := &Bottle{
		Token: token,
	}

	err := stage(bottle)
	assert.Equal(t, fmt.Errorf("Token is Invalid"), err)
}

func TestStoreMessageStage(t *testing.T) {
	pool := NewMessagePool()
	stage := StoreMessageStage(pool)

	messageText := "test"
	message := &Message{
		Text: &messageText,
	}
	bottle := &Bottle{
		Message: message,
	}

	err := stage(bottle)
	assert.Nil(t, err)
}

func TestStoreMessageStageNilMessageText(t *testing.T) {
	pool := NewMessagePool()
	stage := StoreMessageStage(pool)

	message := &Message{
		Text: nil,
	}
	bottle := &Bottle{
		Message: message,
	}

	err := stage(bottle)
	assert.Equal(t, fmt.Errorf("Message Text is Nil"), err)
}
