using System.IO;
using System.Linq;
using System.Threading.Tasks;

namespace crawler;

class Crawled
{
    private readonly string path;

    public Crawled(string path)
    {
        this.path = path;
        File.Create(path).Close();
    }

    public bool HasBeenCrawled(string url) => File.ReadAllLines(path).Any(c => c == url.ToCleanURL());

    public async Task Add(string url)
    {
        using StreamWriter file = new(path, append: true);

        await file.WriteLineAsync(url.ToCleanURL());
    }
}