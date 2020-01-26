import Vue from 'vue'

// Load store modules dynamically.
const requireContext = require.context('@/plugins/static/modules/', false, /.*\.js$/)
function recursiveFreez (obj) {
    let i
    if (typeof obj !== 'object' || obj === null) {
        return obj
    }
    for (i in obj) {
        if (typeof obj[i] === 'object') {
            recursiveFreez(obj[i])
        }
    }

    return Object.freeze(obj)
}

// Load all modules of static objects
const staticObject = requireContext.keys()
    .map(file =>
        [file.replace(/(^.\/)|(\.js$)/g, ''), requireContext(file)]
    )
    .reduce((modules, [name, module]) => {
        if (module.namespaced === undefined) {
            module.namespaced = true
        }

        return { ...modules, [name]: recursiveFreez(module.default) }
    }, {})

Object.freeze(staticObject)

Vue.prototype.$static = staticObject

export default staticObject
