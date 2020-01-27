
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
            var columns = table[i].map(column => {
                const col = new Column(column)
                col.$obs.subscribe(reset)
                return col
            })

            this[i] = new Proxy(columns, {
                get: function (target, name) {
                    return target[name]
                },
                set: function (target, name, val) {
                    if (val instanceof Column) {
                        val.$obs.unsubscribe(reset)
                        val.$obs.subscribe(reset)
                    }

                    target[name] = val

                    return true
                }
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

    this.modifyColumns = function (index, column) {
        column.$obs.subscribe(reset)
        this.Columns[index] = column
        reset()
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
