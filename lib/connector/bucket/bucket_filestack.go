package bucket

import (
	"log"

	client "github.com/filestack/filestack-go/client"
)

//ENV variables- load filestack essentials
type filestack struct {
	Client *client.Client
}

// TODO: add apiKey to env and load to new client
func newFilestack() filestack {
	cli, err := client.NewClient("")
	if err != nil {
		log.Panic("Failed to initialize the client: %v", err)
	}
	return filestack{Client: cli}
}

func (fs filestack) Upload(string, string, string, []byte) error {
	//TODO: implement filestack upload.
	return nil
}
