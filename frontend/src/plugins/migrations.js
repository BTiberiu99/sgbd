import Vue from 'vue'

Vue.prototype.$migrations = {
    run: async function () {
        await window.run.apply(null, arguments)
    }
}
