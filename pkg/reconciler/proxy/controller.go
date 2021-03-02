package proxy

import (
	"context"
	"knative.dev/pkg/configmap"
	"os"

	"knative.dev/pkg/injection"
	"knative.dev/pkg/signals"
	// Injection stuff
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	mwhinformer "knative.dev/pkg/client/injection/kube/informers/admissionregistration/v1/mutatingwebhookconfiguration"
	"knative.dev/pkg/controller"
	secretinformer "knative.dev/pkg/injection/clients/namespacedkube/informers/core/v1/secret"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/system"
	"knative.dev/pkg/webhook"
)

// NewAdmissionController constructs a reconciler
func NewAdmissionController(
	ctx context.Context,
	name, path string,
	wc func(context.Context) context.Context,
	disallowUnknownFields bool,
) *controller.Impl {

	client := kubeclient.Get(ctx)
	mwhInformer := mwhinformer.Get(ctx)
	secretInformer := secretinformer.Get(ctx)
	options := webhook.GetOptions(ctx)

	key := types.NamespacedName{Name: name}

	wh := &reconciler{
		LeaderAwareFuncs: pkgreconciler.LeaderAwareFuncs{
			// Have this reconciler enqueue our singleton whenever it becomes leader.
			PromoteFunc: func(bkt pkgreconciler.Bucket, enq func(pkgreconciler.Bucket, types.NamespacedName)) error {
				enq(bkt, key)
				return nil
			},
		},

		key:  key,
		path: path,

		withContext:           wc,
		disallowUnknownFields: disallowUnknownFields,
		secretName:            options.SecretName,

		client:       client,
		mwhlister:    mwhInformer.Lister(),
		secretlister: secretInformer.Lister(),
	}

	logger := logging.FromContext(ctx)
	c := controller.NewImpl(wh, logger, "DefaultingWebhook")

	// Reconcile when the named MutatingWebhookConfiguration changes.
	mwhInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithName(name),
		// It doesn't matter what we enqueue because we will always Reconcile
		// the named MWH resource.
		Handler: controller.HandleAll(c.Enqueue),
	})

	// Reconcile when the cert bundle changes.
	secretInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithNameAndNamespace(system.Namespace(), wh.secretName),
		// It doesn't matter what we enqueue because we will always Reconcile
		// the named MWH resource.
		Handler: controller.HandleAll(c.Enqueue),
	})

	return c
}

func Getctx() context.Context {
	serviceName := os.Getenv("WEBHOOK_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "tekton-operator-proxy-webhook"
	}

	secretName := os.Getenv("WEBHOOK_SECRET_NAME")
	if secretName == "" {
		secretName = "proxy-webhook-certs"
	}
	systemNamespace := os.Getenv("SYSTEM_NAMESPACE")

	// Scope informers to the webhook's namespace instead of cluster-wide
	ctx := injection.WithNamespaceScope(signals.NewContext(), systemNamespace)

	// Set up a signal context with our webhook options
	ctx = webhook.WithOptions(ctx, webhook.Options{
		ServiceName: serviceName,
		Port:        8443,
		SecretName:  secretName,
	})
	return ctx
}

func NewProxyDefaultingAdmissionController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {

	return NewAdmissionController(ctx,

		// Name of the resource webhook.
		"proxy.operator.tekton.dev",

		// The path on which to serve the webhook.
		"/defaulting",

		// A function that infuses the context passed to Validate/SetDefaults with custom metadata.
		func(ctx context.Context) context.Context {
			return ctx
		},

		// Whether to disallow unknown fields.
		true,
	)
}
