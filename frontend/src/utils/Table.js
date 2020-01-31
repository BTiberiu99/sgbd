
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

            var proxy = new Proxy(columns, {
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

            Object.defineProperty(this, i, {
                get () {
                    return proxy
                },
                set (val) {
                    if (val instanceof Proxy) {
                        proxy = val
                    } else {
                        proxy = new Proxy(val, {
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

        resetCache.push(recalc)
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
            return column.HasNotNull
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
