angular.module('vampLoadbalancer', [
  'ngTable',
  'ngResource',
  'vampLoadbalancer.controllers'
]);
angular.module('vampLoadbalancer.controllers', [])
  .controller('mainController',['$scope', '$http', function($scope, $http, ngTableParams) {


    $http.get('/v1/config')
      .success(function (data){
        $scope.config = data

        data.frontends.forEach(function(fe){
                data.backends.forEach(function(be){
                    if(be.name == fe.defaultBackend) {
                        fe.defaultBackend = be
                    }
                })
            })

      })

    $http.get('/v1/info')
      .success(function (data){
        $scope.haproxy = data
      })

      $scope.showDetails = false

  }]);


