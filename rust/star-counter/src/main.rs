use clap::Parser;
use reqwest::Client;
use serde::{Deserialize, Serialize};
use serde_json;
use std::{
    error::Error,
    fs,
    path::PathBuf,
    time::{SystemTime, UNIX_EPOCH},
};

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
struct CacheData {
    fetched_at: u64,
    repos: Vec<Repo>,
}

async fn fetch_repos(client: &Client, username: &str) -> Result<Vec<Repo>, Box<dyn Error>> {
    let mut all_repos = Vec::new();
    let mut page = 1;

    loop {
        let url = format!(
            "https://api.github.com/users/{}/repos?per_page=100&page={}",
            username, page
        );

        let repos: Vec<Repo> = client
            .get(&url)
            .header("User-Agent", "github-star-counter")
            .send()
            .await?
            .json()
            .await?;

        if repos.is_empty() {
            break; // no more pages
        }

        all_repos.extend(repos);
        page += 1;
    }

    Ok(all_repos)
}

fn cache_file_path(username: &str) -> PathBuf {
    let mut path = PathBuf::from(".cache");
    fs::create_dir_all(&path).ok();
    path.push(format!("github_stars_{}.json", username));
    path
}

fn load_cache(username: &str, ttl: u64) -> Option<Vec<Repo>> {
    let path = cache_file_path(username);
    if let Ok(data) = fs::read_to_string(&path) {
        if let Ok(cache) = serde_json::from_str::<CacheData>(&data) {
            let now = SystemTime::now()
                .duration_since(UNIX_EPOCH)
                .unwrap()
                .as_secs();
            if now - cache.fetched_at < ttl {
                return Some(cache.repos);
            }
        }
    }
    None
}

fn save_cache(username: &str, repos: &[Repo]) {
    let path = cache_file_path(username);
    let data = CacheData {
        fetched_at: SystemTime::now()
            .duration_since(UNIX_EPOCH)
            .unwrap()
            .as_secs(),
        repos: repos.to_vec(),
    };
    if let Ok(json) = serde_json::to_string(&data) {
        fs::write(path, json).ok();
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn Error>> {
    let args = Cli::parse();
    let client = Client::new();

    // Try cache first
    let repos = if let Some(cached) = load_cache(&args.username, args.cache_ttl) {
        println!("Using cached data for '{}'", args.username);
        cached
    } else {
        println!("Fetching data from Github for '{}'", args.username);
        let fresh = fetch_repos(&client, &args.username).await?;
        save_cache(&args.username, &fresh);
        fresh
    };

    let total_stars: u32 = repos.iter().map(|r| r.stargazers_count).sum();

    println!(
        "User '{}' has {} total stars ⭐",
        args.username, total_stars
    );

    Ok(())
}
