import Vue from 'vue'

// Load store modules dynamically.
const requireContext = require.context('@/plugins/static/modules/', false, /.*\.js$/)

const staticObject = requireContext.keys()
    .map(file =>
        [file.replace(/(^.\/)|(\.js$)/g, ''), requireContext(file)]
    )
    .reduce((modules, [name, module]) => {
        if (module.namespaced === undefined) {
            module.namespaced = true
        }

        return { ...modules, [name]: Object.freeze(module.default) }
    }, {})

Object.freeze(staticObject)

Vue.prototype.$static = staticObject
