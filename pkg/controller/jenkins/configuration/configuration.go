package configuration

import (
	"context"
	"fmt"

	"github.com/jenkinsci/kubernetes-operator/pkg/apis/jenkins/v1alpha2"
	"github.com/jenkinsci/kubernetes-operator/pkg/controller/jenkins/configuration/base/resources"
	"github.com/jenkinsci/kubernetes-operator/pkg/controller/jenkins/notifications"

	stackerr "github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Configuration holds required for Jenkins configuration
type Configuration struct {
	Client        client.Client
	ClientSet     kubernetes.Clientset
	Notifications *chan notifications.Event
	Jenkins       *v1alpha2.Jenkins
}

// RestartJenkinsMasterPod terminate Jenkins master pod and notifies about it
func (c *Configuration) RestartJenkinsMasterPod() error {
	currentJenkinsMasterPod, err := c.getJenkinsMasterPod()
	if err != nil {
		return err
	}

	*c.Notifications <- notifications.Event{
		Jenkins:         *c.Jenkins,
		Phase:           notifications.PhaseBase,
		LogLevel:        v1alpha2.NotificationLogLevelInfo,
		Message:         fmt.Sprintf("Terminating Jenkins Master Pod %s/%s.", currentJenkinsMasterPod.Namespace, currentJenkinsMasterPod.Name),
		MessagesVerbose: []string{},
	}

	return stackerr.WithStack(c.Client.Delete(context.TODO(), currentJenkinsMasterPod))
}

func (c *Configuration) getJenkinsMasterPod() (*corev1.Pod, error) {
	jenkinsMasterPodName := resources.GetJenkinsMasterPodName(*c.Jenkins)
	currentJenkinsMasterPod := &corev1.Pod{}
	err := c.Client.Get(context.TODO(), types.NamespacedName{Name: jenkinsMasterPodName, Namespace: c.Jenkins.Namespace}, currentJenkinsMasterPod)
	if err != nil {
		return nil, err // don't wrap error
	}
	return currentJenkinsMasterPod, nil
}
