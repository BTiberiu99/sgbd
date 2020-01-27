import { cache, resetCacheState } from '@/utils/cache'
import Constraint from '@/utils/Constraint'
import staticObj from '@/plugins/static'
import Observable from '@/utils/Observable'

var WatchJS = require('melanke-watchjs')
var watch = WatchJS.watch

const indexConstraints = 'Constraints'
const indexType = 'Type'
export default function (column) {
    var resetCache = []
    let i
    const obs = new Observable()

    // Reset all cache for an column
    const reset = resetCacheState(resetCache, obs.notify())

    this.$obs = obs

    for (i in column) {
        if (i === indexConstraints) {
            var constraints = column[i].map(constrain => {
                const cst = new Constraint(constrain)
                cst.$obs.subscribe(reset)
                return cst
            })

            this[i] = new Proxy(constraints, {
                get: function (target, name) {
                    return target[name]
                },
                set: function (target, name, val) {
                    if (val instanceof Constraint) {
                        val.$obs.unsubscribe(reset)
                        val.$obs.subscribe(reset)
                    }

                    target[name] = val

                    return true
                }
            })
        } else {
            this[i] = column[i]
        }
    }

    for (i in has) {
        const [func, recalc] = cache(has[i], this)

        this[i] = func

        resetCache.push(recalc)
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
            return constraint.IsNotNull() & !constraint.IsPrimaryKey()
        }, this)
    },

    HasPrimaryKey: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsPrimaryKey()
        }, this)
    },

    HasForeignKey: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsForeignKey()
        }, this)
    },

    HasUnique: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsUnique()
        }, this)
    },

    HasCheck: function () {
        return iterateConstraints((constraint) => {
            return constraint.IsCheck()
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
