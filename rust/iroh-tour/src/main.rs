use iroh::{Endpoint, protocol::Router};
use iroh_blobs::{BlobsProtocol, store::mem::MemStore};
use iroh_gossip::net::Gossip;

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let endpoint = Endpoint::builder().discovery_n0().bind().await?;

    // We initialize an in-memory backing store for iroh-blobs
    let store = MemStore::new();
    // Then we initialize a struct that can accept blobs requests over iroh connections
    let blobs = BlobsProtocol::new(&store, endpoint.clone(), None);

    let gossip = Gossip::builder().spawn(endpoint.clone()).await?;

    // build the router
    let router = Router::builder(endpoint)
        .accept(iroh_blobs::ALPN, blobs.clone())
        .accept(iroh_gossip::ALPN, gossip.clone())
        .spawn();

    router.shutdown().await?;

    Ok(())
}
