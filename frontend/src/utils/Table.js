
import Column, { isColumn } from '@/utils/Column'
import Observable from '@/utils/Observable'
import { cache, resetCacheState } from '@/utils/cache'
var WatchJS = require('melanke-watchjs')
var watch = WatchJS.watch

const indexColumns = 'Columns'

export default function (table) {
    let i

    this.keyVue = 0

    var resetCache = [() => {
        this.keyVue++
        this.notify()
    }]

    Observable(this)

    // Reset all cache for an table
    const reset = resetCacheState(resetCache)

    const defineProxy = {
        get: function (target, name) {
            return target[name]
        },
        set: function (target, name, val) {
            if (val instanceof Column) {
                val.unsubscribe(reset).subscribe(reset)
            } else if (isColumn(val)) {
                val = new Column(val).subscribe(reset)
            }

            target[name] = val

            reset()

            return true
        }
    }

    for (i in table) {
        if (i === indexColumns) {
            var proxy = new Proxy(table[i].map(column => {
                return new Column(column).subscribe(reset)
            }), defineProxy)

            Object.defineProperty(this, i, {
                get () {
                    return proxy
                },
                set (val) {
                    if (val instanceof Proxy) {
                        proxy = val
                    } else {
                        proxy = new Proxy(val, defineProxy)
                    }
                }
            })
        } else {
            this[i] = table[i]
        }
    }

    for (i in has) {
        const [func, recalc] = cache(has[i], this)

        Object.defineProperty(this, i, {
            get () {
                return func()
            },
            set () {

            }
        })

        resetCache.splice(0, 0, recalc)
    }

    watch(this, indexColumns, function () {
        reset()
    })

    Object.defineProperty(this, 'IsSafe', {
        get () {
            return this.HasOneNotNull && this.HasPrimaryKey && this.HasCorrectPrimaryKey
        },
        set () {

        }
    })

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
            return column.HasNotNull && !column.HasPrimaryKey
        }, this)
    },
    HasPrimaryKey: function () {
        return iterateColumns((column) => {
            return column.HasPrimaryKey
        }, this)
    },
    HasCorrectPrimaryKey: function () {
        if (!this.HasPrimaryKey) {
            return false
        }

        var countPrimaryKeyNumber = 0
        var isNumericPrimaryKey = true

        iterateColumns((column) => {
            if (column.HasPrimaryKey) {
                if (!column.IsNumeric) {
                    isNumericPrimaryKey = false
                } else {
                    countPrimaryKeyNumber++
                }
            }

            return false
        }, this)

        return countPrimaryKeyNumber < 2 && isNumericPrimaryKey
    },

    Hint: function () {
        if (this.IsSafe) {
            return ''
        }
        var start = `Tabelul ${this.Name} nu are `
        var notnull = 'cel putin o coloana not null inafara de cheia primara '
        var correctPrimaryKey = 'o cheie primara formata corect '
        var primaryKey = ' o cheie primara '
        var hint = ''
        var si = 'si '
        var and = false
        if (!this.HasOneNotNull) {
            hint += `${notnull}`
            and = true
        }

        if (!this.HasPrimaryKey) {
            hint += (and ? si : '') + primaryKey
        }

        if (this.HasPrimaryKey && !this.HasCorrectPrimaryKey) {
            hint += (and ? si : '') + correctPrimaryKey
        }

        return start + hint
    }
}
