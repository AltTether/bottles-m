<template>
<div class="bottle-list">
    <ul>
        <bottle-item v-for="(bottle, index) in reversedBottles"
                     v-bind:key="index"
                     v-bind:message="bottle.message.message"
                     v-bind:token="bottle.token.token"
                     v-bind:tokenDeletedAt="bottle.token.deleted_at">
        </bottle-item>
    </ul>
</div>
</template>

<script>
import BottleItem from './bottle_item.vue';
import store from '../store.js';

export default {
    data () {
        return store.state;
    },

    components: {
        BottleItem,
    },

    mounted () {
        const evtSource = new EventSource("/api/stream",
                                          {withCredentials: true});
        evtSource.addEventListener("ping", (e) => {
            const bottle = JSON.parse(e.data);
            store.addBottle(bottle);
        });
    },

    computed: {
        reversedBottles: function () {
            return store.getBottles().reverse();
        }
    }
}
</script>

<style scoped>
div {
    border:1px solid #000;
    font-family:arial;
    height:500px;
    width:400px;
}

ul {
    list-style:none;
    max-height:500px;
    margin:0;
    overflow:auto;
    padding:0;
    text-indent:10px;
}

.bottle-list {
    margin: 15px;
}
</style>
