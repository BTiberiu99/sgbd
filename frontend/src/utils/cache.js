
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

export const resetCacheState = function (resetCache, call) {
    return function () {
        let i = 0

        for (; i < resetCache.length; i++) {
            resetCache[i]()
        }

        // console.log('RESET CONSTRAINT', constraint)
        if (typeof call === 'function') {
            call()
        }
    }
}
