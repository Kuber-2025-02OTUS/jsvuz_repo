package controllers

import (
    "context"
    appsv1 "k8s.io/api/apps/v1"
    corev1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    "k8s.io/apimachinery/pkg/util/intstr"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log"

    otusv1 "mysql-operator/api/v1"
)

type MySQLReconciler struct {
    client.Client
    Scheme *runtime.Scheme
}

func (r *MySQLReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    logger := log.FromContext(ctx)

    var mysql otusv1.MySQL
    if err := r.Get(ctx, req.NamespacedName, &mysql); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    labels := map[string]string{"app": mysql.Name}

    pvc := &corev1.PersistentVolumeClaim{
        ObjectMeta: metav1.ObjectMeta{
            Name:      mysql.Name + "-pvc",
            Namespace: mysql.Namespace,
        },
        Spec: corev1.PersistentVolumeClaimSpec{
            AccessModes: []corev1.PersistentVolumeAccessMode{
                corev1.ReadWriteOnce,
            },
            Resources: corev1.ResourceRequirements{
                Requests: corev1.ResourceList{
                    corev1.ResourceStorage: resource.MustParse(mysql.Spec.StorageSize),
                },
            },
        },
    }
    _ = ctrl.SetControllerReference(&mysql, pvc, r.Scheme)
    _ = r.Create(ctx, pvc)

    deploy := &appsv1.Deployment{
        ObjectMeta: metav1.ObjectMeta{
            Name:      mysql.Name,
            Namespace: mysql.Namespace,
        },
        Spec: appsv1.DeploymentSpec{
            Selector: &metav1.LabelSelector{MatchLabels: labels},
            Template: corev1.PodTemplateSpec{
                ObjectMeta: metav1.ObjectMeta{Labels: labels},
                Spec: corev1.PodSpec{
                    Containers: []corev1.Container{
                        {
                            Name:  "mysql",
                            Image: mysql.Spec.Image,
                            Env: []corev1.EnvVar{
                                {Name: "MYSQL_DATABASE", Value: mysql.Spec.Database},
                                {Name: "MYSQL_ROOT_PASSWORD", Value: mysql.Spec.Password},
                            },
                            VolumeMounts: []corev1.VolumeMount{
                                {Name: "mysql-storage", MountPath: "/var/lib/mysql"},
                            },
                        },
                    },
                    Volumes: []corev1.Volume{
                        {
                            Name: "mysql-storage",
                            VolumeSource: corev1.VolumeSource{
                                PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
                                    ClaimName: mysql.Name + "-pvc",
                                },
                            },
                        },
                    },
                },
            },
        },
    }
    _ = ctrl.SetControllerReference(&mysql, deploy, r.Scheme)
    _ = r.Create(ctx, deploy)

    svc := &corev1.Service{
        ObjectMeta: metav1.ObjectMeta{
            Name:      mysql.Name,
            Namespace: mysql.Namespace,
        },
        Spec: corev1.ServiceSpec{
            Selector: labels,
            Ports: []corev1.ServicePort{
                {Port: 3306, TargetPort: intstr.FromInt(3306)},
            },
        },
    }
    _ = ctrl.SetControllerReference(&mysql, svc, r.Scheme)
    _ = r.Create(ctx, svc)

    logger.Info("Resources created for MySQL", "name", mysql.Name)
    return ctrl.Result{}, nil
}

func (r *MySQLReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&otusv1.MySQL{}).
        Complete(r)
}
