package model

type Repository struct {
	RemoteURL string
	Version   string
	FileInfo  []*FileInfo
	Type      string
}
