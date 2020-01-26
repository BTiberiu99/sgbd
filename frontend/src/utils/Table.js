
import Column from '@/utils/Column'
import { cache, resetCacheState } from '@/utils/cache'
var WatchJS = require('melanke-watchjs')
var watch = WatchJS.watch

const indexColumns = 'Columns'

export default function (table) {
    let i
    var resetCache = []

    // Reset all cache for an table
    const reset = resetCacheState(resetCache, () => this.keyVue++)

    this.keyVue = 0

    for (i in table) {
        if (i === indexColumns) {
            this[i] = table[i].map(column => {
                const col = new Column(column)
                col.$obs.subscribe(reset)
                return col
            })
        } else {
            this[i] = table[i]
        }
    }

    for (i in has) {
        const [func, recalc] = cache(has[i], this)

        this[i] = func

        resetCache.push(recalc)
    }

    watch(this, indexColumns, function () {
        reset()
    })

    this.IsSafe = function () {
        return this.HasOneNotNull() & this.HasPrimaryKey()
    }

    return this
}

const iterateColumns = function (call, vm) {
    let i

    for (i in vm.Columns) {
        if (call(vm.Columns[i])) {
            return true
        }
    }

    return false
}

const has = {
    HasOneNotNull: function () {
        return iterateColumns((column) => {
            return column.HasNotNull()
        }, this)
    },
    HasPrimaryKey: function () {
        return iterateColumns((column) => {
            return column.HasPrimaryKey()
        }, this)
    }
}
