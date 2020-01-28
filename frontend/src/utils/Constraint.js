import strings from '@/plugins/static'
import { cache, resetCacheState } from '@/utils/cache'
import Observable from '@/utils/Observable'
var WatchJS = require('melanke-watchjs')
var watch = WatchJS.watch

const indexType = 'Type'
export default function (constraint) {
    var resetCache = []
    let i
    const obs = new Observable()

    // Reset all cache for an constraint
    const reset = resetCacheState(resetCache, () => obs.notify())

    this.$obs = obs

    for (i in constraint) {
        this[i] = constraint[i]
    }

    for (i in is) {
        const [func, recalc] = cache(is[i], this)

        Object.defineProperty(this, i, {
            get () {
                return func()
            },
            set () {

            }
        })

        resetCache.push(recalc)
    }

    watch(this, indexType, function () {
        reset()
    })

    return this
}

const is = {
    IsPrimaryKey: function () {
        return this.Type.indexOf(strings.SQL.PRIMARYKEY) !== -1
    },

    IsForeignKey: function () {
        return this.Type.indexOf(strings.SQL.FOREIGNKEY) !== -1
    },

    IsNotNull: function () {
        return this.Type.indexOf(strings.SQL.NOTNULL) !== -1
    },

    IsCheck: function () {
        return this.Type.indexOf(strings.SQL.CHECK) !== -1
    },

    IsUnique: function () {
        return this.Type.indexOf(strings.SQL.UNIQUE) !== -1
    }
}
