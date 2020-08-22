<template>
<div class="box-container">
    <bottle-list/>
    <message-form/>
</div>
</template>

<script>
import BottleList from './bottle_list.vue'
import MessageForm from './message_form.vue'

export default {
    components: {
        BottleList,
        MessageForm,
    },
}
</scripts>

<style scoped>
.box-container {
    display: flex;
}
</style>
