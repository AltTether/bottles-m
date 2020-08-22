export default {
    state: {
        messages: [],
        tokens: [],
        inputMessage: '',
        inputToken: '',
        bottlesLength: 0,
    },

    addBottle(bottle) {
        let message = bottle.message;
        
        let token = bottle.token;
        token.disabled = false;
        
        this.addMessage(message);
        this.addToken(token);
        this.state.bottlesLength++;
    },

    addMessage(message) {
        this.state.messages.push(message);
    },

    addToken(token) {
        this.state.tokens.push(token);
    },

    getBottles() {
        let bottles = [];
        for (let i = 0; i < this.getBottlesLength(); i++) {
            let bottle = this.getBottleByIdx(i)
            bottles.push(bottle);
        }
        return bottles;
    },

    getBottleByIdx(idx) {
        let message = this.state.messages[idx];
        let token = this.state.tokens[idx];
        return {
            "message": message,
            "token": token,
        };
    },
    
    getBottlesLength() {
        return this.state.bottlesLength;
    },

    getTokenByIdx(idx) {
        return this.state.tokens[idx];
    },

    getTokens() {
        return this.state.tokens;
    },

    clearInputs() {
        this.clearInputMessage();
        this.clearInputToken();
    },

    clearInputMessage() {
        this.state.inputMessage = '';
    },
    
    clearInputToken() {
        this.state.inputToken = '';
    },
};
