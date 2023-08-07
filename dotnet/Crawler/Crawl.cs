using HtmlAgilityPack;
using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Threading.Tasks;

namespace crawler;

class Crawl
{
    public readonly string url;
    private string webPage;
    public List<string> parsedURLs;

    public Crawl(string url)
    {
        this.url = url;
        webPage = null;
        parsedURLs = new List<string>();
    }

    public async Task Start()
    {
        await GetWebPage();

        if (!string.IsNullOrWhiteSpace(webPage))
        {
            ParseContent();
            ParseURLs();
        }
    }

    public async Task GetWebPage()
    {
        using HttpClient client = new();

        client.Timeout = TimeSpan.FromSeconds(60);

        string responseBody = await client.GetStringAsync(url);

        if (!string.IsNullOrWhiteSpace(responseBody))
            webPage = responseBody;
    }

    public void ParseURLs()
    {
        HtmlDocument htmlDoc = new HtmlDocument();
        htmlDoc.LoadHtml(webPage);

        foreach (HtmlNode link in htmlDoc.DocumentNode.SelectNodes("//a[@href]"))
        {
            string hrefValue = link.GetAttributeValue("href", string.Empty);

            if (hrefValue.StartsWith("http"))
                parsedURLs.Add(hrefValue);
        }
    }

    public void ParseContent()
    {
        // You may want to process or parse elements of the web page here.
        // Html Agility Pack may also be useful for something like this!
    }
}
