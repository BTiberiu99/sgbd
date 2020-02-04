import 'babel-polyfill'

import Vue from 'vue'
import 'vuetify/dist/vuetify.min.css'
import '@mdi/font/css/materialdesignicons.css'

import App from './App.vue'

import Vuetify from 'vuetify'
import '@/plugins'

import * as Wails from '@wailsapp/runtime'

import { WAILSINIT } from '@/store/events'

Vue.use(Vuetify)
const opts = { theme: { dark: false } }

var vuetify = new Vuetify(opts)
Vue.config.productionTip = false
Vue.config.devtools = true

var cApp = null

function transform (call) {
	return async function () {
		const transformedArguments = []
		let i
		for (i in arguments) {
			transformedArguments.push(typeof arguments[i] === 'object' ? JSON.stringify(arguments[i]) : arguments[i])
		}

		console.log('SEND ==>', transformedArguments)
		var response = await call(...transformedArguments)
		try {
			if (typeof response === 'string') {
				response = JSON.parse(response)
			}
			console.log('RECEIVE <==', response)

			return response
		} catch (e) { console.log(e) }

		return false
	}
}
Wails.Init(() => {
	if (cApp == null) {
		Vue.prototype.$backend = (function () {
			var functions = {}
			Object.keys(window.backend).map(key => {
				functions[key] = transform(window.backend[key])
			})
			return functions
		}())

		cApp = new Vue({
			vuetify,
			render: h => h(App)
		}).$mount('#app')
	} else {
		cApp.$root.$emit(WAILSINIT)
	}
})
