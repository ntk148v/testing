namespace crawler;

static class StringExtensions
{
    public static string ToCleanURL(this string str) => str.Trim().ToLower();
}