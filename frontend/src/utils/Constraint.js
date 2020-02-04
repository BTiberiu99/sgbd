import strings from '@/plugins/static'
import { cache, resetCacheState } from '@/utils/cache'
import Observable from '@/utils/Observable'
var WatchJS = require('melanke-watchjs')
var watch = WatchJS.watch

const indexType = 'Type'

export const isConstraint = (obj) => {
    console.log(obj)
    if (typeof obj !== 'object') {
        return false
    }

    if (obj === null) {
        return false
    }

    if (typeof obj.Type !== 'string' || typeof obj.Name !== 'string') {
        return false
    }

    return true
}
export default function (constraint) {
    let i

    this.keyVue = 0

    Observable(this)

    var resetCache = [() => {
        this.keyVue++
        this.notify()
    }]

    // Reset all cache for an constraint
    const reset = resetCacheState(resetCache)

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

        resetCache.splice(0, 0, recalc)
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
