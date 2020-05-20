package replicasets

import (
	"fmt"
	"strconv"

	"github.com/asobti/kube-monkey/config"
	"github.com/asobti/kube-monkey/victims"

	"k8s.io/api/apps/v1"
)

type ReplicaSet struct {
	*victims.VictimBase
}

// New creates a new instance of ReplicaSet
func New(rep *v1.ReplicaSet) (*ReplicaSet, error) {
	ident, err := identifier(rep)
	if err != nil {
		return nil, err
	}
	mtbf, err := meanTimeBetweenFailures(rep)
	if err != nil {
		return nil, err
	}
	kind := fmt.Sprintf("%T", *rep)

	return &ReplicaSet{VictimBase: victims.New(kind, rep.Name, rep.Namespace, ident, mtbf)}, nil
}

// Returns the value of the label defined by config.IdentLabelKey
// from the deployment labels
// This label should be unique to a deployment, and is used to
// identify the pods that belong to this deployment, as pods
// inherit labels from the ReplicaSet
func identifier(kubekind *v1.ReplicaSet) (string, error) {
	identifier, ok := kubekind.Labels[config.IdentLabelKey]
	if !ok {
		return "", fmt.Errorf("%T %s does not have %s label", kubekind, kubekind.Name, config.IdentLabelKey)
	}
	return identifier, nil
}

// Read the mean-time-between-failures value defined by the ReplicaSet
// in the label defined by config.MtbfLabelKey
func meanTimeBetweenFailures(kubekind *v1.ReplicaSet) (int, error) {
	mtbf, ok := kubekind.Labels[config.MtbfLabelKey]
	if !ok {
		return -1, fmt.Errorf("%T %s does not have %s label", kubekind, kubekind.Name, config.MtbfLabelKey)
	}

	mtbfInt, err := strconv.Atoi(mtbf)
	if err != nil {
		return -1, err
	}

	if !(mtbfInt > 0) {
		return -1, fmt.Errorf("Invalid value for label %s: %d", config.MtbfLabelKey, mtbfInt)
	}

	return mtbfInt, nil
}
