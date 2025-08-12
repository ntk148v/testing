use clap::Parser;
use reqwest::Client;
use serde::{Deserialize, Serialize};
use serde_json;
use std::{collections::HashMap, error::Error, fs, path::Path};

#[derive(Parser, Debug)]
#[command(name = "GitHub Star Counter")]
#[command(about = "Counts total stars for a GitHub user", long_about = None)]
struct Cli {
    /// GitHub username
    username: String,
    /// Cache TTL in seconds (default 3600 = 1 hour)
    #[arg(long, default_value_t = 3600)]
    cache_ttl: u64,
}

#[derive(Debug, Clone, Deserialize, Serialize)]
struct Repo {
    stargazers_count: u32,
}

#[derive(Debug, Deserialize, Serialize)]
struct CacheEntry {
    etag: String,
    repos: Vec<Repo>,
}

type Cache = HashMap<String, CacheEntry>;

const CACHE_FILE: &str = "/tmp/github_star_counter.json";

fn load_cache() -> Cache {
    if Path::new(CACHE_FILE).exists() {
        let data = fs::read_to_string(CACHE_FILE).unwrap_or_default();
        serde_json::from_str(&data).unwrap_or_default()
    } else {
        HashMap::new()
    }
}

fn save_cache(cache: &Cache) {
    let data = serde_json::to_string_pretty(cache).unwrap();
    fs::write(CACHE_FILE, data).unwrap();
}

// ETag-based caching
// - The Github API response includes an ETag header
// - You store this ETag/user in a local cache file.
// - On the next request, you send it as If-None-Match header
// - If Github returns 304 NOT_MODIFIED, you use the cached data instead of downloading again
// <https://docs.github.com/en/rest/using-the-rest-api/best-practices-for-using-the-rest-api?apiVersion=2022-11-28#use-conditional-requests-if-appropriate>
async fn fetch_repos(
    client: &Client,
    username: &str,
    cache: &mut Cache,
) -> Result<Vec<Repo>, Box<dyn Error>> {
    let mut all_repos = Vec::new();
    let mut page = 1;

    loop {
        let url = format!(
            "https://api.github.com/users/{}/repos?per_page=100&page={}",
            username, page
        );

        let mut req = client.get(&url).header("user-agent", "github-star-counter");

        // If we have an ETag, send it
        if let Some(entry) = cache.get(username) {
            req = req.header("if-none-match", &entry.etag);
        }

        let resp = req.send().await?;

        if resp.status() == reqwest::StatusCode::NOT_MODIFIED {
            println!("Using cached data for '{}'", username);
            return Ok(cache.get(username).unwrap().repos.clone());
        }

        let etag = resp
            .headers()
            .get(reqwest::header::ETAG)
            .and_then(|v| v.to_str().ok())
            .unwrap_or("")
            .to_string();

        let repos: Vec<Repo> = resp.json().await?;

        if repos.is_empty() {
            break; // No more pages
        }

        all_repos.extend(repos);
        page += 1;

        // Store cache after first page (GitHub API uses same ETag for all paginated calls in most cases)
        cache.insert(
            username.to_string(),
            CacheEntry {
                etag: etag.clone(),
                repos: all_repos.clone(),
            },
        );
        save_cache(cache);
    }

    Ok(all_repos)
}
#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let args = Cli::parse();
    let client = Client::new();
    let mut cache = load_cache();

    let repos = fetch_repos(&client, &args.username, &mut cache).await?;

    let total_stars: u32 = repos.iter().map(|r| r.stargazers_count).sum();

    println!(
        "User '{}' has {} total stars ⭐",
        args.username, total_stars
    );

    Ok(())
}
