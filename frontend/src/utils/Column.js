import { cache, resetCacheState } from '@/utils/cache'
import Constraint, { isConstraint } from '@/utils/Constraint'
import staticObj from '@/plugins/static'
import Observable from '@/utils/Observable'

var WatchJS = require('melanke-watchjs')
var watch = WatchJS.watch

const indexConstraints = 'Constraints'
const indexType = 'Type'

export const isColumn = (obj) => {
    if (typeof obj !== 'object') {
        return false
    }

    if (obj === null) {
        return false
    }

    if (typeof obj.Name !== 'string' || typeof obj.Type !== 'string' || !Array.isArray(obj.Constraints)) {
        return false
    }

    return true
}
export default function (column) {
    Observable(this)

    this.keyVue = 0

    var resetCache = [() => {
        this.keyVue++
        this.notify()
    }]

    // Reset all cache for a column
    const reset = resetCacheState(resetCache)

    const defineProxy = {
        get: function (target, name) {
            return target[name]
        },
        set: function (target, name, val) {
            if (val instanceof Constraint) {
                val.unsubscribe(reset).subscribe(reset)
            } else if (isConstraint(val)) {
                val = new Constraint(val).subscribe(reset)
            }

            target[name] = val

            reset()

            return true
        }
    }

    let i
    for (i in column) {
        if (i === indexConstraints) {
            var proxy = new Proxy(column[i].map(constrain => {
                return new Constraint(constrain).subscribe(reset)
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
            this[i] = column[i]
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

    watch(this, indexType, function () {
        reset()
    })

    watch(this, indexConstraints, function () {
        reset()
    })

    return this
}

const has = {
    HasNotNull: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsNotNull
        }, this)
    },

    HasPrimaryKey: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsPrimaryKey
        }, this)
    },

    HasForeignKey: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsForeignKey
        }, this)
    },

    HasUnique: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsUnique
        }, this)
    },

    HasCheck: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsCheck
        }, this)
    },

    // Autocreate truth functions to check type of column
    ...(function () {
        const obj = {}
        const types = ['Numeric', 'String', 'Binary', 'Date', 'Boolean', 'Geometric']
        let i
        for (i in types) {
            const type = types[i].toUpperCase()

            obj['Is' + types[i]] = function () {
                return staticObj.SQLTypes.postgres[type].values.indexOf(this.Type) !== -1
            }
        }
        return obj
    }())
}

const iterateConstraints = function (call, vm) {
    let i
    for (i in vm.Constraints) {
        if (call(vm.Constraints[i])) {
            return true
        }
    }

    return false
}
