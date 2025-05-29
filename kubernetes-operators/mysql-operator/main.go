package main

import (
    "os"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client/config"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"

    otusv1 "mysql-operator/api/v1"
    "mysql-operator/controllers"
)

func main() {
    ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
    mgr, _ := ctrl.NewManager(config.GetConfigOrDie(), ctrl.Options{})
    otusv1.AddToScheme(mgr.GetScheme())
    controllers.MySQLReconciler{
        Client: mgr.GetClient(),
        Scheme: mgr.GetScheme(),
    }.SetupWithManager(mgr)
    mgr.Start(ctrl.SetupSignalHandler())
}
