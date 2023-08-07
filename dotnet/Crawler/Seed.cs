using System.IO;

namespace crawler;

class Seed
{
    /// <summary>
    /// Returns all seed URLs.
    /// </summary>
    public string[] Items
    {
        get => File.ReadAllLines(path);
    }

    private readonly string path;

    public Seed(string path)
    {
        this.path = path;

        string[] seedURLs = new string[]
        {
                "https://ntk148v.github.io"
        };

        using StreamWriter file = File.CreateText(path);

        foreach (string url in seedURLs)
            file.WriteLine(url.ToCleanURL());
    }
}
