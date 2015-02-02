var apiURL = 'http://localhost:10001/v1/config'

var demo = new Vue({

    el: '#main',

    data: {
        config: null,
        sortedConfig: null
    },

    created: function () {
        this.fetchData();
    },

//    filters: {
//        truncate: function (v) {
//            var newline = v.indexOf('\n')
//            return newline > 0 ? v.slice(0, newline) : v
//        },
//        formatDate: function (v) {
//            return v.replace(/T|Z/g, ' ')
//        }
//    },

    methods: {
        fetchData: function () {
            var xhr = new XMLHttpRequest()
            var self = this
            xhr.open('GET', apiURL)
            xhr.onload = function () {
                self.config = JSON.parse(xhr.responseText)
                self.sortData(self.config)
            }
            xhr.send()
        },
        sortData: function(config) {
            var sortedConfig = [];
            config.frontends.forEach(function(fe){
                console.log(fe)
                config.backends.forEach(function(be){
                    if(be.name == fe.defaultBackend) {
                        fe.defaultBackend = be
                    }
                })
                sortedConfig.push(fe)
            })
            console.log(JSON.stringify(sortedConfig))
            self.sortedConfig = sortedConfig

        }
    }
})

