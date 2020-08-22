<template>
<li>
    <div class="message">{{ message }}</div>
    <div class="box-container">
        <div class="token">{{ token }}</div>
        <div class="deleted-at">{{ convertDate(tokenDeletedAt) }}</div>
    </div>
</li>
</template>

<script>
export default {
    props: ['message', 'token', 'tokenDeletedAt'],
    methods: {
        convertDate: function(utcDateString) {
            const utcDate = new Date(utcDateString);
            const localDate = utcDate.toString();
            return localDate;
        }
    },
}
</script>

<style scoped>
li {
    line-height: 25px;
    border-bottom: thin solid #000000;
}

.message {
    font-size: 100%;
    white-space: pre-wrap;
    word-wrap: break-word;
    text-indent: 0;
}

.box-container {
    display: flex;
}

.token {
    font-size: 30%;
}

.deleted-at {
    font-size: 30%;
}
</style>
