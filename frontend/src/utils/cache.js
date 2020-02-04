
export const cache = function (call, vm) {
    var m = null
    return [function () {
        if (m === null) {
            m = call.apply(vm, null)
        }

        return m
    }, function () {
        m = null
    }]
}

export const resetCacheState = function (resetCache) {
    return async function () {
        let i = 0

        for (; i < resetCache.length; i++) {
            resetCache[i]()
        }
    }
}
