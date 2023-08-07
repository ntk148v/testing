using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Threading.Tasks;

namespace crawler;


class Queue
{
    /// <summary>
    /// Returns the first item in the queue.
    /// </summary>
    public string Top
    {
        get => File.ReadAllLines(path).First();
    }

    /// <summary>
    /// Returns all items in the queue;
    /// </summary>
    public string[] All
    {
        get => File.ReadAllLines(path);
    }

    /// <summary>
    /// Returns a value based on whether there are URLs in the queue.
    /// </summary>
    public bool HasURLs
    {
        get => File.ReadAllLines(path).Length > 0;
    }

    private readonly string path;

    public Queue(string path, string[] seedURLs)
    {
        this.path = path;

        using StreamWriter file = File.CreateText(path);

        foreach (string url in seedURLs)
            file.WriteLine(url.ToCleanURL());
    }

    public async Task Add(string url)
    {
        using StreamWriter file = new(path, append: true);

        await file.WriteLineAsync(url.ToCleanURL());
    }

    public async Task Remove(string url)
    {
        IEnumerable<string> filteredURLs = All.Where(u => u != url);

        await File.WriteAllLinesAsync(path, filteredURLs);
    }

    public bool IsInQueue(string url) => All.Where(u => u == url).Any();
}