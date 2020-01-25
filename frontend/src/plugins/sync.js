import Vue from 'vue'
var interval = null
var count = 0
var promises = {}
const mseconds = 3 * 60 * 4 // 3 minutes * 60 seconds * 4 (1000/250ms)
function checkPromises () {
    let ok = false
    let i
    for (i in promises) {

        if (!promises[i].call()) {
            // Remove promise
            promises[i].remove()
        }

        // Set to keep interval
        if (!ok) ok = true
    }


    if (!ok) {
        const copy = interval
        interval = null
        clearInterval(copy)
    }
}

function setPromise (promise) {
    // Trigger interval if interval is disabled
    if (interval === null) interval = setInterval(checkPromises, 250)
    // Set promise
    var c = count++
    promises[c] = {
        call: promise,
        remove: () => {
            delete promises[c]
        }
    }
}


// Sync async functions with other async functions
Vue.prototype.$sync = function sync (call) {
    if (typeof call !== 'function') return

    return new Promise((resolve, reject) => {
        var counter = 0
        setPromise(() => {
            let continueWatching = true
            if (call()) {
                continueWatching = false
                resolve()
            }
            counter++
            if (counter > mseconds && continueWatching) {
                continueWatching = false
                // eslint-disable-next-line prefer-promise-reject-errors
                reject()
            }
            return continueWatching
        })
    }).catch(e => { e })
}


