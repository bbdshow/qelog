package types

import "os"

type ServerMode string

const (
	Single          ServerMode = "single" // 默认单节点部署
	ClusterAdmin    ServerMode = "cluster_admin"
	ClusterReceiver ServerMode = "cluster_receiver"
)

// GetFlagOrOSEnvServerMode flag or os env server mode
func GetFlagOrOSEnvServerMode(mode ServerMode) ServerMode {
	if mode == Single || mode == ClusterAdmin || mode == ClusterReceiver {
		return mode
	}

	emode := os.Getenv("SERVER_MODE")
	switch ServerMode(emode) {
	case ClusterAdmin:
		return ClusterAdmin
	case ClusterReceiver:
		return ClusterReceiver
	}
	return Single
}

func GetOSEnvAdminPassword() string {
	return os.Getenv("ADMIN_PASSWORD")
}
