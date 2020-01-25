import 'babel-polyfill';
import Vue from 'vue';

import 'vuetify/dist/vuetify.min.css';
import '@mdi/font/css/materialdesignicons.css';

import App from './App.vue';

import Vuetify from 'vuetify'
import '@/plugins'

Vue.use(Vuetify)
const opts = { theme: { dark: false } }

var vuetify = new Vuetify(opts)
Vue.config.productionTip = false;
Vue.config.devtools = true;

import * as Wails from '@wailsapp/runtime';

var cApp = null



function parse (key) {
	return async function () {
		const transformedArguments = []
		let i
		for (i in arguments) {
			transformedArguments.push(JSON.stringify(arguments[i]))
		}

		console.log('SEND==>', transformedArguments)
		const response = await window.backend[key](...transformedArguments)
		console.log('RECIVE<==', JSON.parse(response))
		return JSON.parse(response)
	}
}
Wails.Init(() => {
	if (cApp == null) {

		var functions = {}
		Object.keys(window.backend).map(key => {
			functions[key] = parse(key)
		})
		Vue.prototype.$backend = functions

		cApp = new Vue({
			vuetify,
			render: h => h(App)
		}).$mount('#app');


	} else {

		cApp.$root.$emit('wailsinit')
	}

});
