package objects

// NetworkNodeInfo represents common networking information in OTN blockchain
type NetworkNodeInfo struct {
	ListeningOn     string `json:"listening_on"`
	NodePublicKey   string `json:"node_public_key"`
	NodeID          string `json:"node_id"`
	Firewalled      string `json:"firewalled"`
	ConnectionCount UInt32 `json:"connection_count"`
}

// AdvancedNodeParameters shows networking node parameters
type AdvancedNodeParameters struct {
	PeerConnectionRetryTimeout             UInt32 `json:"peer_connection_retry_timeout"`
	DesiredNumberOfConnections             UInt32 `json:"desired_number_of_connections"`
	MaximumNumberOfConnections             UInt32 `json:"maximum_number_of_connections"`
	MaximumNumberOfBlocksToHandleAtOneTime UInt32 `json:"maximum_number_of_blocks_to_handle_at_one_time"`
	MaximumNumberOfSyncBlocksToPrefetch    UInt32 `json:"maximum_number_of_sync_blocks_to_prefetch"`
	MaximumBlocksPerPeerDuringSyncing      UInt32 `json:"maximum_blocks_per_peer_during_syncing"`
}

// PeerStatusInfo represents some service p2p information
type PeerStatusInfo struct {
	Addr                       string `json:"addr"`
	Addrlocal                  string `json:"addrlocal"`
	Services                   string `json:"services"`
	Lastsend                   UInt64 `json:"lastsend"`
	Lastrecv                   UInt64 `json:"lastrecv"`
	Bytessent                  UInt64 `json:"bytessent"`
	Bytesrecv                  UInt64 `json:"bytesrecv"`
	Conntime                   string `json:"conntime"`
	Pingtime                   string `json:"pingtime"`
	Pingwait                   string `json:"pingwait"`
	Version                    string `json:"version"`
	Subver                     string `json:"subver"`
	Inbound                    bool   `json:"inbound"`
	FirewallStatus             string `json:"firewall_status"`
	Startingheight             string `json:"startingheight"`
	Banscore                   string `json:"banscore"`
	Syncnode                   string `json:"syncnode"`
	FcGitRevisionSha           string `json:"fc_git_revision_sha"`
	FcGitRevisionUnixTimestamp string `json:"fc_git_revision_unix_timestamp"`
	FcGitRevisionAge           string `json:"fc_git_revision_age"`
	Platform                   string `json:"platform"`
	CurrentHeadBlock           string `json:"current_head_block"`
	CurrentHeadBlockNumber     UInt64 `json:"current_head_block_number"`
	CurrentHeadBlockTime       string `json:"current_head_block_time"`
}

// PeerStatus represents information about peer connected to OTN blockchain node
type PeerStatus struct {
	Version int            `json:"version"`
	Host    string         `json:"host"`
	Info    PeerStatusInfo `json:"info"`
}

// PotentialPeerRecord represents struct
type PotentialPeerRecord struct {
	Endpoint                             string `json:"endpoint"`
	LastSeenTime                         string `json:"last_seen_time"`
	LastConnectionDisposition            string `json:"last_connection_disposition"`
	LastConnectionAttemptTime            string `json:"last_connection_attempt_time"`
	NumberOfSuccessfulConnectionAttempts UInt32 `json:"number_of_successful_connection_attempts"`
	NumberOfFailedConnectionAttempts     UInt32 `json:"number_of_failed_connection_attempts"`
	LastError                            *struct {
		Code    UInt32 `json:"code"`
		Name    string `json:"name"`
		Message string `json:"message"`
		Stack   []struct {
			Context struct {
				Level      string `json:"level"`
				File       string `json:"file"`
				Line       UInt32 `json:"line"`
				Method     string `json:"method"`
				Hostname   string `json:"hostname"`
				ThreadName string `json:"thread_name"`
				Timestamp  string `json:"timestamp"`
			} `json:"context"`
			Format string `json:"format"`
			Data   struct {
				Message string `json:"message"`
			} `json:"data"`
		} `json:"stack"`
	} `json:"last_error"`
}
