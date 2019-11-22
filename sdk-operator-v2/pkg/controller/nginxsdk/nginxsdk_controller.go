package nginxsdk

import (
	"context"
	"reflect"

	sdkoperatorv1alpha1 "github.com/arbhoj/ansible-operator/pkg/apis/sdkoperator/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	// Arvind added custom
	"k8s.io/apimachinery/pkg/util/intstr"
)

var log = logf.Log.WithName("controller_nginxsdk")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new NginxSDK Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileNginxSDK{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("nginxsdk-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource NginxSDK
	err = c.Watch(&source.Kind{Type: &sdkoperatorv1alpha1.NginxSDK{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner NginxSDK
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sdkoperatorv1alpha1.NginxSDK{},
	})
	if err != nil {
		return err
	}

	// Arvind added custom
	//Watch for changes in the service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sdkoperatorv1alpha1.NginxSDK{},
	})
	if err != nil {
		return err
	}
	// Arvind added custom
	//Watch for changes in the configMap
	err = c.Watch(&source.Kind{Type: &corev1.ConfigMap{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sdkoperatorv1alpha1.NginxSDK{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileNginxSDK implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileNginxSDK{}

// ReconcileNginxSDK reconciles a NginxSDK object
type ReconcileNginxSDK struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a NginxSDK object and makes changes based on the state read
// and what is in the NginxSDK.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileNginxSDK) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling NginxSDK")

	// Fetch the NginxSDK instance
	instance := &sdkoperatorv1alpha1.NginxSDK{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	deploy := newDeployForCR(instance)
	// Set NginxSDK instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deploy, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Deploy already exists
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		err = r.client.Create(context.TODO(), deploy)
		reqLogger.Info("Created a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err == nil && (instance.Spec.Size != *found.Spec.Replicas || instance.Spec.DefaultBody != found.Spec.Template.Spec.Containers[0].Env[0].Value || deploy.Spec.Template.Spec.Containers[0].Image != found.Spec.Template.Spec.Containers[0].Image) {
		found.Spec.Replicas = &instance.Spec.Size
		found.Spec.Template.Spec.Containers[0].Env[0].Value = instance.Spec.DefaultBody
                found.Spec.Template.Spec.Containers[0].Image = deploy.Spec.Template.Spec.Containers[0].Image
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			reqLogger.Error(err, "Failed to update Depployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Updating Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name, "Deploy.Size", deploy.Spec.Replicas)
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Info("Error", err)
	}

	// Deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", deploy.Name)

	// Arvind added custom
	// Define a new Service object
	service := newServiceForCR(instance)
	// Set NginxSDK instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	foundService := &corev1.Service{}
	//Arvind added custom
	//Get the node port specified in the owner CRD to check if the
	//service port matches the spec in the owner CRD
	crdNodeport := instance.Spec.NodePort

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
		//Arvind Added Custom
		//NodePort does not mactch spec of the owner CRD
	} else if err == nil && crdNodeport != foundService.Spec.Ports[0].NodePort {
		foundService.Spec.Ports[0].NodePort = crdNodeport
		err = r.client.Update(context.TODO(), foundService)
		if err != nil {
			reqLogger.Error(err, "Failed to update service", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Updating Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name, "Service.Port", service.Spec.Ports[0].NodePort)
		return reconcile.Result{Requeue: true}, nil

		// Service created successfully - don't requeue
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Service already exists - don't requeue
	reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)

	// Arvind added custom
	// Define a new ConfigMap object
	configMap := newConfigMapForCR(instance)
	// Set NginxSDK instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, configMap, r.scheme); err != nil {
		return reconcile.Result{}, err
	}
	foundConfigMap := &corev1.ConfigMap{}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, foundConfigMap)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name)
		err = r.client.Create(context.TODO(), configMap)
		if err != nil {
			return reconcile.Result{}, err
		}

		// ConfigMap created successfully - don't requeue
	} else if err == nil && instance.Spec.DefaultBody != foundConfigMap.Data["index.html"] {
		foundConfigMap.Data["index.html"] = instance.Spec.DefaultBody
		err = r.client.Update(context.TODO(), foundConfigMap)
		if err != nil {
			reqLogger.Error(err, "Failed to update configmap", "ConfigMap.Namespace", foundConfigMap.Namespace, "ConfigMap.Name", foundConfigMap.Name)
			return reconcile.Result{}, err
		}
		reqLogger.Info("Updating ConfigMap", "ConfigMap.Namespace", configMap.Namespace, "ConfigMap.Name", configMap.Name, "ConfigMap.Data", configMap.Data["index.html"])
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// ConfigMap already exists - don't requeue
	reqLogger.Info("Skip reconcile: ConfigMap already exists", "ConfigMap.Namespace", foundConfigMap.Namespace, "ConfigMap.Name", foundConfigMap.Name)

	//Arvind added custom
	statusNodes := map[string]string{
		"Deployment": deploy.ObjectMeta.Name,
		"Service":    service.ObjectMeta.Name,
		"ConfigMap":  configMap.ObjectMeta.Name,
	}
	if !reflect.DeepEqual(statusNodes, instance.Status.Nodes) {
		reqLogger.Info("Updating Nodes", "Deployment", deploy.ObjectMeta.Name, "Service", service.ObjectMeta.Name, "ConfigMap", configMap.ObjectMeta.Name)
		instance.Status.Nodes = statusNodes
		err = r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update the status of instance")
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, nil
}

// newDeployment returns a nginx deployment object matching the specs specified in the main CRD
func newDeployForCR(cr *sdkoperatorv1alpha1.NginxSDK) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}
	//Arvind added custom
	annotations := map[string]string{
		"arvind": "This is the updated annotation for v1.0.5",
	}
	configMapObject := corev1.ConfigMapVolumeSource{
		LocalObjectReference: corev1.LocalObjectReference{
			Name: cr.Name + "-conf",
		},
	}
	size := cr.Spec.Size
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
			//Arvind added custom
			Annotations: annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
					Name:   cr.Name + "-pod",
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "my-nginx",
							Image: "nginx:latest",
							// Arvind added custom
							Env: []corev1.EnvVar{
								{
									Name:  "Test",
									Value: cr.Spec.DefaultBody,
								},
 								{
									Name: "New",
									Value: "NewValue",	
 								},	
							},	
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 80,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "myconfigmap",
									MountPath: "/usr/share/nginx/html/index.html",
									SubPath:   "index.html",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "myconfigmap",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &configMapObject,
							},
						},
					},
				},
			},
		},
	}
}

//Arvind added custom
func newServiceForCR(cr *sdkoperatorv1alpha1.NginxSDK) *corev1.Service {
	selectors := map[string]string{
		"app": cr.Name,
	}
	targetPort := intstr.IntOrString{
		IntVal: cr.Spec.ContainerPort,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-service",
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       cr.Name,
					Protocol:   "TCP",
					Port:       80,
					TargetPort: targetPort,
					NodePort:   cr.Spec.NodePort,
				},
			},
			Selector: selectors,
			Type:     corev1.ServiceTypeNodePort,
		},
	}
}

//Arvind added custom
func newConfigMapForCR(cr *sdkoperatorv1alpha1.NginxSDK) *corev1.ConfigMap {
	configData := map[string]string{
		"index.html": cr.Spec.DefaultBody,
	}
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-conf",
			Namespace: cr.Namespace,
		},
		Data: configData,
	}
}
