<template>
<div class="message-form">
    <div><textarea v-model="inputMessage"></textarea></div>
    <div class="box-container">
        <select v-model="inputTokenIndex">
            <option disabled value="">Select Token</option>
            <option v-for="(token, index) in tokens"
                    v-if="!token.disabled"
                    v-bind:value="index">
                {{ token.token }}
            </option>
        </select>
        <button
            v-on:click="send">
            Send
        </button>
        {{arrayTokens}}
    </div>
</div>
</template>

<script>
import axios from 'axios';
import store from '../store.js'

export default {
    data () {
         return store.state;
    },

    methods: {
        send: function(event) {
            const message = this.inputMessage;
            const token = store.getTokenByIdx(this.inputTokenIndex).token;

            axios
            .post('/api', {message, token})
            .then(response => {
                if (response.status == 200) {
                    store.clearInputs();
                    // TODO store内で変更する
                    this.$set(this.tokens[this.inputTokenIndex], "disabled", true);
                }
            });
        },
    },
}
</script>

<style scoped>
textarea {
    resize: none;
    height: 450px;
    width: 400px;
}

.box-container {
    display: flex;
}

.message-form {
    margin: 15px;
}
</style>
