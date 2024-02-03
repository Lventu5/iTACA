package retrieve

/*
// MongoRetriever : implements Retriever for MongoDB, fetching reconstructed TCP streams from a MongoDB database
// address must be the address of the machine which is hosting the database, port must be the port which exposes the database
type MongoRetriever struct {
	address string
	port    uint16
	storage storage.Storage
}

// Retrieve : retrieves tcp streams from a MongoDB database
func (r *MongoRetriever) Retrieve(ctx context.Context, results chan<- Result) {
	var connections []storage.Connection

	for {
		select {
		case <-ctx.Done():
			return
		default:
			query := r.storage.Find("connections").Context(ctx)
			query = query.Limit(50)
			if err := query.All(&connections); err != nil {
				fmt.Printf("failed to retrieve connections: %s\n", err)
				os.Exit(1)
			}
			for _, connection := range connections {

			}
		}
	}
}
*/
