use anyhow::{Context, Result};
use clap::Parser;
use std::{
    fs::File,
    io::{stdout, BufRead, BufReader, Write},
};

/// Search for a pattern in a file and display the lines that contain it.
#[derive(Parser)]
struct Cli {
    /// The pattern to look for
    pattern: String,
    /// The path to the file to read
    path: std::path::PathBuf,
}

fn main() -> Result<()> {
    let args = Cli::parse();

    let content = std::fs::read_to_string(&args.path)
        .with_context(|| format!("could not read file `{}`", args.path.display()))?;

    find_matches(&content, &args.pattern, &mut stdout());

    // Pass ownership of the open file to a BufReader struct. BufReader uses an internal
    // buffer to reduce intermediate allocations
    let file = File::open(&args.path)
        .with_context(|| format!("could not read file `{}`", args.path.display()))?;
    let reader = BufReader::new(file);

    for (_index, line_result) in reader.lines().enumerate() {
        let line = line_result?;
        if line.contains(&args.pattern) {
            println!("{}", line);
        }
    }

    Ok(())
}

fn find_matches(content: &str, pattern: &str, mut writer: impl Write) {
    for line in content.lines() {
        if line.contains(pattern) {
            writeln!(writer, "{}", line);
        }
    }
}

#[test]
fn find_a_match() {
    let mut result = Vec::new();
    find_matches("lorem ipsum\ndolor sit amet", "lorem", &mut result);
    assert_eq!(result, b"lorem ipsum\n");
}
